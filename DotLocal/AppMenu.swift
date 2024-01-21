//
//  AppMenu.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import SwiftUI
import AppKit
import LaunchAtLogin

struct AppMenu: View {
    @Environment(\.openWindow) var openWindow
    @StateObject var daemonManager = DaemonManager.shared

    var body: some View {
        switch daemonManager.state {
        case .stopped:
            Button("DotLocal is not running") {}.disabled(true)
        case .starting, .unknown:
            Button("DotLocal is starting") {}.disabled(true)
        case .started(let savedState):
            Section("Routes") {
                MappingListMenu(mappings: savedState.mappings)
            }
        }
        Divider()
        if #available(macOS 14.0, *) {
            SettingsLink().keyboardShortcut(",")
        } else {
            Button(action: {
                NSApp.sendAction(Selector(("showSettingsWindow:")), to: nil, from: nil)
            }, label: { Text("Settings...") })
        }
        Button("Quit DotLocal") {
            NSApplication.shared.terminate(nil)
        }.keyboardShortcut("Q")
    }
}

struct MappingListMenu: View {
    var mappings: [Mapping]
    @Environment(\.openURL) var openURL
    
    var body: some View {
        if mappings.isEmpty {
            Button("No Routes", action: {}).disabled(true)
        } else {
            ForEach(mappings) { mapping in
                let url = URL(string: "http://\(mapping.host)\(mapping.pathPrefix)")!
                Button(action: { openURL(url) }, label: {
                    Text(getLabel(mapping: mapping))
                })
            }
        }
    }
    
    func getLabel(mapping: Mapping) -> AttributedString {
        var title = AttributedString("\(mapping.host)\(mapping.pathPrefix)")
        title.font = .system(size: NSFont.systemFontSize)
        
        var subtitle = AttributedString("\n\(mapping.target)")
        subtitle.font = .system(size: NSFont.smallSystemFontSize)
        subtitle.foregroundColor = .secondary
        
        return title + subtitle
    }
}
