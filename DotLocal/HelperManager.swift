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
    
    private let service = SMAppService.daemon(plistName: "helper.plist")
    @Published var status: SMAppService.Status
    
    let xpcClient = XPCClient.forMachService(named: "dev.suphon.DotLocal.helper")
    
    private init() {
        status = service.status
        Task {
            print("sending exit to current helper")
            do {
                try await xpcClient.send(to: SharedConstants.exitRoute)
            } catch {
                print("error sending exit: \(error)")
            }
            do {
                print("registering service")
                try service.register()
                print("registered service")
            } catch {
                print("error registering service: \(error)")
            }
            await checkStatus()
        }
    }
    
    @MainActor
    func checkStatus() {
        status = service.status
    }
}
