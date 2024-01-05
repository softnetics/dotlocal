//
//  DaemonManager.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 5/1/2567 BE.
//

import Foundation

class DaemonManager: ObservableObject {
    static let shared = DaemonManager()
    
    private let binUrl = Bundle.main.bundleURL.appendingPathComponent("Contents/Resources/bin")
    private let daemonUrl = Bundle.main.bundleURL.appendingPathComponent("Contents/Resources/dotlocal-daemon")
    @Published var state: DaemonState = .stopped
    private var task: Process? = nil
    
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
                    self.state = .started
                }
            }
            handle.waitForDataInBackgroundAndNotify()
        }
        handle.waitForDataInBackgroundAndNotify()
        
        DispatchQueue.global().async {
            task.waitUntilExit()
            DispatchQueue.main.async {
                self.onStop()
            }
        }
        self.task = task
    }
    
    private func onStop() {
        task = nil
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
