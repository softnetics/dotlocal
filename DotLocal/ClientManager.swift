//
//  ClientManager.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 7/1/2567 BE.
//

import Foundation

class ClientManager: ObservableObject {
    static let shared = ClientManager()
    
    @Published var installed = false
    private let clientUrl = Bundle.main.bundleURL.appendingPathComponent("Contents/Resources/bin/dotlocal")
    private let target = "/usr/local/bin/dotlocal"
    
    private init() {}
    
    func installCli() async {
        _ = await Sudo.run(path: clientUrl.path(percentEncoded: false), arguments: ["install"])
        checkInstalled()
    }
    
    func uninstallCli() async {
        _ = await Sudo.run(path: clientUrl.path(percentEncoded: false), arguments: ["uninstall"])
        checkInstalled()
    }
    
    func checkInstalled() {
        DispatchQueue.main.async {
            self.installed = FileManager.default.isExecutableFile(atPath: self.target)
        }
    }
}
