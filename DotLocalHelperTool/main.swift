//
//  main.swift
//  DotLocalHelperTool
//
//  Created by Suphon Thanakornpakapong on 11/1/2567 BE.
//

import Foundation
import SecureXPC

NSLog("starting helper tool. PID \(getpid()). PPID \(getppid()).")
NSLog("version: \(try HelperToolInfoPropertyList.main.version.rawValue)")

if getppid() == 1 {
    let server = try XPCServer.forMachService()
    server.registerRoute(SharedConstants.installClientRoute, handler: ManageClient.install)
    server.registerRoute(SharedConstants.uninstallClientRoute, handler: ManageClient.uninstall)
    server.registerRoute(SharedConstants.exitRoute, handler: {
        NSLog("exiting")
        exit(0)
    })
    server.setErrorHandler { error in
        if case .connectionInvalid = error {
            // Ignore invalidated connections as this happens whenever the client disconnects which is not a problem
        } else {
            NSLog("error: \(error)")
        }
    }
    server.startAndBlock()
} else {
    print("not supported")
}
