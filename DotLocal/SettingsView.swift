//
//  Settings.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import SwiftUI
import LaunchAtLogin
import Defaults

struct GeneralSettingsView: View {
    @Default(.showInMenuBar) var showInMenuBar
    
    var body: some View {
        Form {
            LaunchAtLogin.Toggle("Open at Login")
            Toggle("Show in Menu Bar", isOn: $showInMenuBar)
        }
        .padding(20)
        .frame(minWidth: 350, maxWidth: 350)
    }
}

struct SettingsView: View {
    var body: some View {
        GeneralSettingsView()
    }
}
