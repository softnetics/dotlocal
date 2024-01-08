//
//  DotLocalApp.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 5/1/2567 BE.
//

import SwiftUI
import AppKit
import Defaults

@main
struct DotLocalApp: App {
    @NSApplicationDelegateAdaptor(AppDelegate.self) var appDelegate
    @Default(.showInMenuBar) var showInMenuBar
    
    var body: some Scene {
        MenuBarExtra("DotLocal", systemImage: "server.rack", isInserted: .constant(showInMenuBar)) {
            AppMenu()
        }
        
        WindowGroup() {
            ContentView()
        }.commands {
            CommandGroup(replacing: .newItem) {}
        }
        
        Settings {
            SettingsView()
        }
    }
}
