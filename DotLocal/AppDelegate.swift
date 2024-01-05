//
//  AppDelegate.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 5/1/2567 BE.
//

import Foundation
import AppKit

class AppDelegate: NSObject, NSApplicationDelegate {
    func applicationDidFinishLaunching(_ notification: Notification) {
        DaemonManager.shared.start()
    }
    
    func applicationShouldTerminateAfterLastWindowClosed(_ sender: NSApplication) -> Bool {
        sender.setActivationPolicy(.accessory)
        return false
    }
    
    func applicationWillTerminate(_ notification: Notification) {
        DaemonManager.shared.stop()
        DaemonManager.shared.wait()
    }
}
