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
    static let shared = DaemonManager()
    
    @Published var state: DaemonState = .unknown
    @Published var mappings: [Mapping] = []
    
    private init() {
    }
    
    func start() async {
        do {
            print("starting daemon")
            try await HelperManager.shared.xpcClient.send(to: SharedConstants.startDaemonRoute)
            print("successfully requested start")
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
    
    func subscribeDaemonState() async {
        do {
            for try await state in HelperManager.shared.xpcClient.send(to: SharedConstants.daemonStateRoute) {
                DispatchQueue.main.async {
                    self.state = state
                    if case .started(let mappings) = state {
                        self.mappings = mappings
                    }
                }
            }
        } catch {
            print("error during state subscription: \(error)")
        }
        if HelperManager.shared.status == .enabled {
            Task {
                await subscribeDaemonState()
            }
        }
    }
}
