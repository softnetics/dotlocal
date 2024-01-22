//
//  DaemonManager.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 5/1/2567 BE.
//

import Foundation
import GRPC
import NIO
import Combine

class DaemonManager: ObservableObject {
    private let group = PlatformSupport.makeEventLoopGroup(loopCount: 1)
    let apiClient: DotLocalAsyncClient
    
    static let shared = DaemonManager()
    
    @Published var state: DaemonState = .unknown
    @Published var savedState: SavedState = SavedState()
    
    private var subscribing = false
    
    private init() {
        apiClient = try! createApiClient(group: group)
    }
    
    func start() async {
        do {
            print("starting daemon")
            try await HelperManager.shared.xpcClient.sendMessage(Bundle.main.bundleURL, to: SharedConstants.startDaemonRoute)
            print("successfully requested start")
            Task {
                await subscribeDaemonState()
            }
        } catch {
            print("error starting daemon: \(error)")
        }
    }
    
    func stop() async {
        do {
            print("stopping daemon")
            try await HelperManager.shared.xpcClient.send(to: SharedConstants.stopDaemonRoute)
            print("successfully requested stop")
        } catch {
            print("error stopping daemon: \(error)")
        }
    }
    
    private func subscribeDaemonState() async {
        if subscribing {
            return
        }
        subscribing = true
        do {
            for try await state in HelperManager.shared.xpcClient.send(to: SharedConstants.daemonStateRoute) {
                DispatchQueue.main.async {
                    self.state = state
                    if case .started(let savedState) = state {
                        self.savedState = savedState
                    }
                }
            }
        } catch {
            print("error during state subscription: \(error)")
        }
    }
}
