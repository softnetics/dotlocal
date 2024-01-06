//
//  DotLocalApp.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 5/1/2567 BE.
//

import SwiftUI
import AppKit

@main
struct DotLocalApp: App {
    @NSApplicationDelegateAdaptor(AppDelegate.self) var appDelegate
    
    var body: some Scene {
        MenuBarExtra("DotLocal", systemImage: "server.rack", isInserted: .constant(true)) {
            AppMenu()
        }
        
        WindowGroup() {
            ContentView()
        }
    }
}
