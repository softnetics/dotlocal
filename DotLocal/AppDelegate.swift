//
//  AppDelegate.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 5/1/2567 BE.
//

import Foundation
import AppKit
import Defaults
import SecureXPC

class AppDelegate: NSObject, NSApplicationDelegate {
    func applicationDidFinishLaunching(_ notification: Notification) {
        DaemonManager.shared.start()
        ClientManager.shared.checkInstalled()
        _ = HelperManager.shared
    }
    
    func applicationShouldTerminateAfterLastWindowClosed(_ sender: NSApplication) -> Bool {
        if Defaults[.showInMenuBar] {
            sender.setActivationPolicy(.accessory)
        }
        return false
    }
    
    func applicationShouldHandleReopen(_ sender: NSApplication, hasVisibleWindows _: Bool) -> Bool {
        sender.setActivationPolicy(.regular)
        return true
    }
    
    func applicationWillTerminate(_ notification: Notification) {
        DaemonManager.shared.stop()
        DaemonManager.shared.wait()
    }
}
