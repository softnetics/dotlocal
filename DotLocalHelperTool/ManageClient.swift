//
//  ManageCLI.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 12/1/2567 BE.
//

import Foundation

private func clientLocation() -> URL {
    let appURL = try! parentAppURL()
    return appURL.appendingPathComponent("Contents/Resources/bin/dotlocal")
}

private let installLocation = URL.init(filePath: "/usr/local/bin/dotlocal")

enum ManageClient {
    static func install() throws {
        let source = clientLocation()
        NSLog("installing client")
        NSLog("symlink \(source) to \(installLocation)")
        try FileManager.default.createDirectory(at: installLocation.deletingLastPathComponent(), withIntermediateDirectories: true)
        try FileManager.default.createSymbolicLink(at: installLocation, withDestinationURL: source)
        NSLog("installed client")
    }
    
    static func uninstall() throws {
        NSLog("uninstalling client")
        NSLog("remove \(installLocation)")
        try FileManager.default.removeItem(at: installLocation)
        NSLog("uninstalled client")
    }
}
