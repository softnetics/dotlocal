//
//  ContentView.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 5/1/2567 BE.
//

import SwiftUI
import ServiceManagement
import Blessed
import Authorized
import SecureXPC

struct ContentView: View {
    @StateObject var daemonManager = DaemonManager.shared
    @StateObject var helperManager = HelperManager.shared
    
    var body: some View {
        let status = helperManager.installationStatus
        VStack {
            if status.isReady {
                VStack {
                    switch daemonManager.state {
                    case .stopped:
                        Text("DotLocal is not running")
                    case .starting, .unknown:
                        ProgressView()
                    case .started:
                        MappingList()
                    }
                }.toolbar() {
                    StartStopButton(state: daemonManager.state, onStart: {
                        Task {
                            await daemonManager.start()
                        }
                    }, onStop: {
                        Task {
                            await daemonManager.stop()
                        }
                    })
                }
            } else {
                RequiresHelperView()
            }
        }
    }
}

struct StartStopButton: View {
    var state: DaemonState
    var onStart: () -> Void
    var onStop: () -> Void
    
    var body: some View {
        switch state {
        case .stopped:
            Button(action: onStart) {
                Label("Start", systemImage: "play.fill")
            }
        case .starting, .unknown:
            ProgressView().controlSize(.small)
        case .started:
            Button(action: onStop) {
                Label("Stop", systemImage: "stop.fill")
            }
        }
    }
}

struct RequiresHelperView: View {
    @State private var didError = false
    @State private var errorMessage = ""
    
    var body: some View {
        VStack(spacing: 8) {
            Text("Helper Not Installed").font(.title).fontWeight(.bold)
            Text("Please install the helper in order to use DotLocal")
            Button(action: {
                do {
                    try PrivilegedHelperManager.shared
                        .authorizeAndBless(message: nil)
                } catch AuthorizationError.canceled {
                    // No user feedback needed, user canceled
                } catch {
                    errorMessage = error.localizedDescription
                    didError = true
                }
            }, label: {
                Text("Install Helper")
            })
        }
        .foregroundStyle(.secondary)
        .alert(
            "Install failed",
            isPresented: $didError,
            presenting: errorMessage
        ) { _ in
            Button("OK") {}
        } message: { message in
            Text(message)
        }
    }
}

#Preview {
    ContentView()
}
