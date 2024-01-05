//
//  DaemonManager.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 5/1/2567 BE.
//

import Foundation
import GRPC
import NIO

class DaemonManager: ObservableObject {
    static let shared = DaemonManager()
    
    private let binUrl = Bundle.main.bundleURL.appendingPathComponent("Contents/Resources/bin")
    private let daemonUrl = Bundle.main.bundleURL.appendingPathComponent("Contents/Resources/dotlocal-daemon")
    @Published var state: DaemonState = .stopped
    private var task: Process? = nil
    private(set) var apiClient: DotLocalAsyncClient? = nil
    private var group: EventLoopGroup? = nil
    
    private init() {
    }
    
    func start() {
        if state != .stopped {
            return
        }
        state = .starting
        
        let binPath = binUrl.path(percentEncoded: false)
        let launchPath = daemonUrl.path(percentEncoded: false)
        
        let task = Process()
        var environment = ProcessInfo.processInfo.environment
        environment["PATH"] = environment["PATH"]! + ":" + binPath
        task.environment = environment
        task.launchPath = launchPath
        
        let outputPipe = Pipe()
        task.standardError = outputPipe
        task.launch()
        
        let handle = outputPipe.fileHandleForReading
        let token = NotificationCenter.default.addObserver(forName: .NSFileHandleDataAvailable, object: outputPipe.fileHandleForReading, queue: nil) { _ in
            let chunk = String(decoding: handle.availableData, as: UTF8.self)
            print(chunk, terminator: "")
            if chunk.contains("API server listening") {
                DispatchQueue.main.async {
                    self.onStart()
                }
            }
            handle.waitForDataInBackgroundAndNotify()
        }
        handle.waitForDataInBackgroundAndNotify()
        
        DispatchQueue.global().async {
            task.waitUntilExit()
            DispatchQueue.main.async {
                NotificationCenter.default.removeObserver(token)
                self.onStop()
            }
        }
        self.task = task
    }
    
    private func onStart() {
        let socketPath = FileManager.default.homeDirectoryForCurrentUser.appendingPathComponent(".dotlocal/api.sock")
        let group = PlatformSupport.makeEventLoopGroup(loopCount: 1)
        self.group = group
        // TODO: try catch
        let channel = try! GRPCChannelPool.with(
            target: .unixDomainSocket(socketPath.path(percentEncoded: false)),
            transportSecurity: .plaintext,
            eventLoopGroup: group
        )
        let apiClient = DotLocalAsyncClient(channel: channel)
        self.apiClient = apiClient
        
        state = .started
    }
    
    private func onStop() {
        task = nil
        apiClient = nil
        state = .stopped
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
}

enum DaemonState {
    case stopped
    case starting
    case started
}