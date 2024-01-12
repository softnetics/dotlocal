//
//  ClientManager.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 7/1/2567 BE.
//

import Foundation
import SecureXPC

class ClientManager: ObservableObject {
    static let shared = ClientManager()
    
    @Published var installed = false
    private let target = "/usr/local/bin/dotlocal"
    
    private init() {}
    
    func installCli() async {
        do {
            try await HelperManager.shared.xpcClient.sendMessage(Bundle.main.bundleURL, to: SharedConstants.installClientRoute)
            checkInstalled()
        } catch {
            print("error installing cli: \(error)")
        }
    }
    
    func uninstallCli() async {
        do {
            try await HelperManager.shared.xpcClient.send(to: SharedConstants.uninstallClientRoute)
            checkInstalled()
        } catch {
            print("error uninstalling cli: \(error)")
        }
    }
    
    func checkInstalled() {
        DispatchQueue.main.async {
            self.installed = FileManager.default.isExecutableFile(atPath: self.target)
        }
    }
}
