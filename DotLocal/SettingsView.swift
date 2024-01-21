//
//  Settings.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import SwiftUI
import SecurityInterface
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
        .frame(minWidth: 400, maxWidth: 400)
    }
}

struct HttpsView: View {
    @Binding var prefs: Preferences
    @State var rootCertificate: CertHelper.Certificate?
    @State private var window: NSWindow?
    
    @Environment(\.controlActiveState) var controlActiveState
    
    var body: some View {
        LabeledContent(content: {
            VStack(alignment: .leading) {
                Toggle("Enable", isOn: !$prefs.disableHTTPS)
                Toggle("Automatically redirect from HTTP", isOn: $prefs.redirectHTTPS).disabled(prefs.disableHTTPS)
            }
        }, label: {
            Text("HTTPS:")
        })
        .background(WindowAccessor(window: $window))
        .onAppear {
            Task {
                await loadCertificate()
            }
        }
        .onChange(of: controlActiveState) { activeState in
            Task {
                if activeState == .key {
                    await loadCertificate()
                }
            }
        }
        if let rootCertificate = rootCertificate, let window = window {
            CertificateView(certificate: rootCertificate, window: window, reload: {
                Task {
                    await loadCertificate()
                }
            })
        }
    }
    
    private func loadCertificate() async {
        do {
            rootCertificate = try await CertHelper.getRootCertificate()
        } catch {
            print("failed to load certificate: \(error)")
        }
    }
}

struct CertificateView: View {
    @State private var didError = false
    @State private var errorTitle = ""
    @State private var errorMessage = ""
    
    var certificate: CertHelper.Certificate
    var window: NSWindow
    var reload: () -> Void
    
    var body: some View {
        LabeledContent(content: {
            VStack(alignment: .leading) {
                Text(certificate.commonName)
                Text("Expires: \(certificate.notAfter.formatted(date: .abbreviated, time: .shortened))")
                    .foregroundStyle(.secondary)
                    .font(.system(size: 12))
                if certificate.trusted {
                    (Text(Image(systemName: "checkmark.seal.fill"))+Text(" Trusted"))
                        .foregroundStyle(.green)
                        .font(.system(size: 12))
                } else {
                    (Text(Image(systemName: "xmark.circle.fill"))+Text(" Not trusted"))
                        .foregroundStyle(.red)
                        .font(.system(size: 12))
                }
                HStack {
                    Button(action: {
                        SFCertificatePanel.shared().beginSheet(for: window, modalDelegate: nil, didEnd: nil, contextInfo: nil, certificates: [certificate.secCertificate], showGroup: false)
                    }, label: {
                        Text("Details")
                    })
                    if !certificate.trusted {
                        Button(action: {
                            var status = SecItemAdd([
                                kSecClass as String: kSecClassCertificate,
                                kSecValueRef as String: certificate.secCertificate
                            ] as CFDictionary, nil)
                            guard status == errSecSuccess || status == errSecDuplicateItem else {
                                errorTitle = "Failed to add certificate to Keychain"
                                errorMessage = "\(status): " + (SecCopyErrorMessageString(status, nil)! as String)
                                didError = true
                                return
                            }
                            status = SecTrustSettingsSetTrustSettings(certificate.secCertificate, .user, [
                                [kSecTrustSettingsPolicy: SecPolicyCreateBasicX509()],
                                [kSecTrustSettingsPolicy: SecPolicyCreateSSL(true, nil)],
                            ] as CFTypeRef)
                            guard status == errSecSuccess || status == errAuthorizationCanceled else {
                                errorTitle = "Failed to set trust settings for certificate"
                                errorMessage = "\(status): " + (SecCopyErrorMessageString(status, nil)! as String)
                                didError = true
                                return
                            }
                            reload()
                        }, label: {
                            Text("Trust")
                        })
                    }
                }
            }
        }, label: {
            Text("Root Certificate:")
        })
        .alert(
            errorTitle,
            isPresented: $didError,
            presenting: errorMessage
        ) { _ in
            Button("OK") {}
        } message: { message in
            Text(message)
        }
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

struct WindowAccessor: NSViewRepresentable {
    @Binding var window: NSWindow?
    
    func makeNSView(context: Context) -> NSView {
        let view = NSView()
        DispatchQueue.main.async {
            self.window = view.window
        }
        return view
    }
    
    func updateNSView(_ nsView: NSView, context: Context) {}
}

prefix func ! (value: Binding<Bool>) -> Binding<Bool> {
    Binding<Bool>(
        get: { !value.wrappedValue },
        set: { value.wrappedValue = !$0 }
    )
}
