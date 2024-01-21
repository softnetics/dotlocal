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
import SecureXPC

struct GeneralSettingsView: View {
    @Default(.showInMenuBar) var showInMenuBar
    @StateObject var daemonManager = DaemonManager.shared
    
    let prefs: Binding<Preferences>
    
    init() {
        prefs = Binding(
            get: {
                if case .started(let savedState) = DaemonManager.shared.state {
                    return savedState.preferences
                } else {
                    return Preferences()
                }
            },
            set: { value in
                Task {
                    try await DaemonManager.shared.apiClient.setPreferences(value)
                }
            }
        )
    }
    
    var body: some View {
        Form {
            LaunchAtLogin.Toggle("Open at Login")
            Toggle("Show in Menu Bar", isOn: $showInMenuBar)
            if case .started = daemonManager.state {
                HttpsView(prefs: prefs).padding(.top, 8)
            }
            CliView().padding(.top, 8)
        }
        .padding(20)
        .frame(minWidth: 350, maxWidth: 350)
    }
}

struct HttpsView: View {
    @Binding var prefs: Preferences
    
    var body: some View {
        LabeledContent(content: {
            VStack(alignment: .leading) {
                Toggle("Enable", isOn: !$prefs.disableHTTPS)
                Toggle("Automatically redirect from HTTP", isOn: $prefs.redirectHTTPS).disabled(prefs.disableHTTPS)
            }
        }, label: {
            Text("HTTPS:")
        })
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

prefix func ! (value: Binding<Bool>) -> Binding<Bool> {
    Binding<Bool>(
        get: { !value.wrappedValue },
        set: { value.wrappedValue = !$0 }
    )
}
