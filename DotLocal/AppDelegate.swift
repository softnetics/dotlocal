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
    override init() {
        _ = HelperManager.shared
        ClientManager.shared.checkInstalled()
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
    
    func applicationShouldTerminate(_ sender: NSApplication) -> NSApplication.TerminateReply {
        print("applicationShouldTerminate called, stopping daemon and helper")
        if HelperManager.shared.installationStatus.isReady {
            Task {
                await DaemonManager.shared.stop()
                try? await HelperManager.shared.xpcClient.send(to: SharedConstants.exitRoute)
                NSApplication.shared.terminate(nil)
            }
            return .terminateLater
        } else {
            return .terminateNow
        }
    }
}
