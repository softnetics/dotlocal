//
//  DaemonManager.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 12/1/2567 BE.
//

import Foundation
import GRPC
import NIO
import Combine
import SecureXPC

class DaemonManager {
    static let shared = DaemonManager()
    
    private let runDirectory = URL.init(filePath: "/var/run/dotlocal")
    
    var internalState: DaemonStateInternal = .stopped
    private var task: Process? = nil
    private(set) var apiClient: DotLocalAsyncClient? = nil
    private let group = PlatformSupport.makeEventLoopGroup(loopCount: 1)
    private let _updates = PassthroughSubject<Void, Never>()
    
    private let _internalStates = PassthroughSubject<DaemonStateInternal, Never>()
    private let states = CurrentValueSubject<DaemonState, Never>(.stopped)
    
    private var subscriptions = Set<AnyCancellable>()
    
    private init() {
        _internalStates
            .flatMap { state in
                Future<DaemonState, Never> { promise in
                    Task {
                        promise(.success(await DaemonManager.mapState(state)))
                    }
                }
            }
            .subscribe(states)
            .store(in: &subscriptions)
    }
    
    func start(bundleURL: URL) async throws {
        let binURL = bundleURL.appending(path: "Contents/Resources/bin")
        let daemonURL = bundleURL.appending(path: "Contents/Resources/dotlocal-daemon")
        
        NSLog("received start")
        NSLog("daemonURL: \(daemonURL)")
        
        guard try CodeInfo.doesPublicKeyMatch(forExecutable: daemonURL) else {
            NSLog("start daemon failed: security requirements not met")
            return
        }
        
        NSLog("security requirements passed")
        if internalState != .stopped {
            return
        }
        setState(.starting)
        
        let binPath = binURL.path(percentEncoded: false)
        let launchPath = daemonURL.path(percentEncoded: false)
        
        try! FileManager.default.createDirectory(at: runDirectory, withIntermediateDirectories: true)
        
        let task = Process()
        var environment = ProcessInfo.processInfo.environment
        environment["PATH"] = binPath
        task.environment = environment
        task.launchPath = launchPath
        task.currentDirectoryURL = runDirectory
        
        task.standardOutput = FileHandle.standardOutput
        let outputPipe = Pipe()
        task.standardError = outputPipe
        outputPipe.fileHandleForReading.readabilityHandler = { handle in
            let chunk = String(decoding: handle.availableData, as: UTF8.self)
            print(chunk, terminator: "")
            if chunk.contains("API server listening") {
                DispatchQueue.main.async {
                    self.onStart()
                }
            } else if chunk.contains("Updated state") {
                NSLog("sending update")
                self.setState(.started)
            }
        }
        
        task.terminationHandler = { _ in
            self.onStop()
        }
        
        task.launch()
        NSLog("launched")
        
        self.task = task
    }
    
    private func onStart() {
        // TODO: try catch
        apiClient = try! createApiClient(group: group)
        
        setState(.started)
    }
    
    private func onStop() {
        task = nil
        apiClient = nil
        setState(.stopped)
    }
    
    func stop() {
        guard let task = task else {
            return
        }
        task.terminate()
    }
    
    func wait() {
        guard let task = task else {
            return
        }
        task.waitUntilExit()
    }
    
    private func setState(_ newState: DaemonStateInternal) {
        internalState = newState
        NSLog("new state: \(newState)")
        _internalStates.send(newState)
    }
    
    private static func mapState(_ internalState: DaemonStateInternal) async -> DaemonState {
        switch internalState {
        case .stopped:
            return .stopped
        case .starting:
            return .starting
        case .started:
            guard let apiClient = DaemonManager.shared.apiClient else {
                return .started(savedState: SavedState())
            }
            var res = (try? await apiClient.getSavedState(.with{_ in})) ?? SavedState()
            res.mappings.sort()
            return .started(savedState: res)
        }
    }
    
    func daemonState(provider: SequentialResultProvider<DaemonState>) async {
        var subscriptions = Set<AnyCancellable>()
        states.sink(receiveCompletion: { _ in
            if !provider.isFinished {
                provider.respond(withResult: .finished)
            }
        }, receiveValue: {
            if !provider.isFinished {
                provider.respond(withResult: .success($0))
            } else {
                subscriptions.removeAll()
            }
        })
        .store(in: &subscriptions)
    }
}

enum DaemonStateInternal {
    case stopped
    case starting
    case started
}
