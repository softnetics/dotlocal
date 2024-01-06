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
        case .starting:
            Button("DotLocal is starting") {}.disabled(true)
        case .started:
            Section("Routes") {
                MappingListMenu()
            }
        }
        Divider()
        SettingsLink().keyboardShortcut(",")
        Button("Quit DotLocal") {
            NSApplication.shared.terminate(nil)
        }.keyboardShortcut("Q")
    }
}

struct MappingListMenu: View {
    @StateObject var vm = MappingListViewModel()
    @Environment(\.openURL) var openURL
    
    var body: some View {
        if vm.mappings.isEmpty {
            Button("No Routes", action: {}).disabled(true)
        } else {
            ForEach(vm.mappings) { mapping in
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
