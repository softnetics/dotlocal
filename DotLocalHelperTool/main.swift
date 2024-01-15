//
//  main.swift
//  DotLocalHelperTool
//
//  Created by Suphon Thanakornpakapong on 11/1/2567 BE.
//

import Foundation
import SecureXPC
import Dispatch

NSLog("starting helper tool. PID \(getpid()). PPID \(getppid()).")
NSLog("version: \(try HelperToolInfoPropertyList.main.version.rawValue)")
NSLog("code location: \(String(describing: try? CodeInfo.currentCodeLocation()))")

// Command line arguments were provided, so process them
if CommandLine.arguments.count > 1 {
    // Remove the first argument, which represents the name (typically the full path) of this helper tool
    var arguments = CommandLine.arguments
    _ = arguments.removeFirst()
    NSLog("run with arguments: \(arguments)")
    
    if let firstArgument = arguments.first {
        if firstArgument == Uninstaller.commandLineArgument {
            try Uninstaller.uninstallFromCommandLine(withArguments: arguments)
        } else {
            NSLog("argument not recognized: \(firstArgument)")
        }
    }
} else if getppid() == 1 { // Otherwise if started by launchd, start up server
    let server = try XPCServer.forMachService()
    
    server.registerRoute(SharedConstants.startDaemonRoute, handler: DaemonManager.shared.start)
    server.registerRoute(SharedConstants.stopDaemonRoute, handler: DaemonManager.shared.stop)
    server.registerRoute(SharedConstants.daemonStateRoute, handler: DaemonManager.shared.daemonState)
    
    server.registerRoute(SharedConstants.installClientRoute, handler: ManageClient.install)
    server.registerRoute(SharedConstants.uninstallClientRoute, handler: ManageClient.uninstall)
    
    server.registerRoute(SharedConstants.exitRoute, handler: gracefulExit)
    
    server.registerRoute(SharedConstants.uninstallRoute, handler: Uninstaller.uninstallFromXPC)
    server.registerRoute(SharedConstants.updateRoute, handler: Updater.updateHelperTool(atPath:))
    
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
} else { // Otherwise started via command line without arguments, print out help info
    print("Usage: \(try CodeInfo.currentCodeLocation().lastPathComponent) <command>")
    print("\nCommands:")
    print("\t\(Uninstaller.commandLineArgument)\tUnloads and deletes from disk this helper tool and configuration.")
}

func gracefulExit() {
    NSLog("exiting")
    DaemonManager.shared.stop()
    exit(0)
}
