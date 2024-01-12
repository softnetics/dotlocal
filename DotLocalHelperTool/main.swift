//
//  main.swift
//  DotLocalHelperTool
//
//  Created by Suphon Thanakornpakapong on 11/1/2567 BE.
//

import Foundation
import SecureXPC
import Dispatch

func parentAppURL() throws -> URL {
    let components = Bundle.main.bundleURL.pathComponents
    guard let contentsIndex = components.lastIndex(of: "Contents"),
          components[components.index(before: contentsIndex)].hasSuffix(".app") else {
        throw MyError.runtimeError("""
        Parent bundle could not be found.
        Path:\(Bundle.main.bundleURL)
        """)
    }
    
    return URL(fileURLWithPath: "/" + components[1..<contentsIndex].joined(separator: "/"))
}

NSLog("starting helper tool. PID \(getpid()). PPID \(getppid()).")
NSLog("version: \(try HelperToolInfoPropertyList.main.version.rawValue)")

if getppid() == 1 {
    let server = try XPCServer.forMachService(withCriteria: .forDaemon(withClientRequirement: try! .sameParentBundle))
    
    server.registerRoute(SharedConstants.startDaemonRoute, handler: DaemonManager.shared.start)
    server.registerRoute(SharedConstants.stopDaemonRoute, handler: DaemonManager.shared.stop)
    server.registerRoute(SharedConstants.daemonStateRoute, handler: DaemonManager.shared.daemonState)
    
    server.registerRoute(SharedConstants.installClientRoute, handler: ManageClient.install)
    server.registerRoute(SharedConstants.uninstallClientRoute, handler: ManageClient.uninstall)
    
    server.registerRoute(SharedConstants.exitRoute, handler: gracefulExit)
    server.setErrorHandler { error in
        if case .connectionInvalid = error {
            // Ignore invalidated connections as this happens whenever the client disconnects which is not a problem
        } else {
            NSLog("error: \(error)")
        }
    }
    
    signal(SIGINT, SIG_IGN)
    signal(SIGTERM, SIG_IGN)
    
    let sigintSrc = DispatchSource.makeSignalSource(signal: SIGINT, queue: .main)
    sigintSrc.setEventHandler(handler: gracefulExit)
    sigintSrc.resume()
    let sigtermSrc = DispatchSource.makeSignalSource(signal: SIGTERM, queue: .main)
    sigtermSrc.setEventHandler(handler: gracefulExit)
    sigtermSrc.resume()
    
    server.startAndBlock()
} else {
    print("not supported")
}

func gracefulExit() {
    NSLog("exiting")
    DaemonManager.shared.stop()
    exit(0)
}
