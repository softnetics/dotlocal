//
//  HelperManager.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 12/1/2567 BE.
//

import Foundation
import ServiceManagement
import SecureXPC

class HelperManager: ObservableObject {
    static let shared = HelperManager()
    
    private let monitor = HelperToolMonitor(constants: SharedConstants.shared)
    @Published var installationStatus: HelperToolMonitor.InstallationStatus
    private var started = false
    
    let xpcClient = XPCClient.forMachService(named: SharedConstants.shared.machServiceName)
    
    private init() {
        installationStatus = monitor.determineStatus()
        monitor.start { status in
            DispatchQueue.main.async {
                self.updateStatus(status: status)
            }
        }
        Task {
            if installationStatus.isReady {
                try await updateHelper()
            }
        }
    }
    
    private func updateHelper() async throws {
        do {
            print("updating helper")
            try await xpcClient.sendMessage(SharedConstants.shared.bundledLocation, to: SharedConstants.updateRoute)
        } catch XPCError.connectionInterrupted {
            print("update success")
            return
        } catch {
            print("update error: \(error)")
            throw error
        }
    }
    
    private func updateStatus(status: HelperToolMonitor.InstallationStatus) {
        installationStatus = status
        if status.isReady, !started {
            started = true
            Task {
                await DaemonManager.shared.start()
            }
        }
    }
}
