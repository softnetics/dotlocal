//
//  ManageCLI.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 12/1/2567 BE.
//

import Foundation

private let installLocation = URL.init(filePath: "/usr/local/bin/dotlocal")

enum ManageClient {
    static func install(bundleURL: URL) throws {
        let clientURL = bundleURL.appending(path: "Contents/Resources/bin/dotlocal")
        NSLog("installing client")
        
        guard try CodeInfo.doesPublicKeyMatch(forExecutable: clientURL) else {
            NSLog("start daemon failed: security requirements not met")
            return
        }
        
        NSLog("symlink \(clientURL) to \(installLocation)")
        try FileManager.default.createDirectory(at: installLocation.deletingLastPathComponent(), withIntermediateDirectories: true)
        try FileManager.default.createSymbolicLink(at: installLocation, withDestinationURL: clientURL)
        NSLog("installed client")
    }
    
    static func uninstall() throws {
        NSLog("uninstalling client")
        NSLog("remove \(installLocation)")
        try FileManager.default.removeItem(at: installLocation)
        NSLog("uninstalled client")
    }
}
