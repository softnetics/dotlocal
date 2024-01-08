//
//  Settings.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import SwiftUI
import LaunchAtLogin
import Defaults
import Foundation

struct GeneralSettingsView: View {
    @Default(.showInMenuBar) var showInMenuBar
    
    var body: some View {
        Form {
            LaunchAtLogin.Toggle("Open at Login")
            Toggle("Show in Menu Bar", isOn: $showInMenuBar)
            CliView().padding(.top, 8)
        }
        .padding(20)
        .frame(minWidth: 350, maxWidth: 350)
    }
}

struct CliView: View {
    @StateObject private var clientManager = ClientManager.shared
    
    var body: some View {
        LabeledContent(content: {
            if clientManager.installed {
                VStack(alignment: .leading) {
                    Text("Installed to /usr/local/bin/dotlocal")
                    Button("Uninstall", action: {
                        Task {
                            await clientManager.uninstallCli()
                        }
                    })
                }
            } else {
                Button("Install", action: {
                    Task {
                        await clientManager.installCli()
                    }
                })
            }
        }, label: {
            Text("Command Line:")
        })
        .onAppear {
            clientManager.checkInstalled()
        }
    }
}

struct SettingsView: View {
    var body: some View {
        GeneralSettingsView()
    }
}
