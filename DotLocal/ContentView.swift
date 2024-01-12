//
//  ContentView.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 5/1/2567 BE.
//

import SwiftUI
import ServiceManagement

struct ContentView: View {
    @StateObject var daemonManager = DaemonManager.shared
    @StateObject var helperManager = HelperManager.shared
    
    var body: some View {
        switch helperManager.status {
        case .requiresApproval:
            RequiresApprovalView()
        case .enabled:
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
        default:
            Text("Unexpected state: \(helperManager.status.rawValue)")
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

struct RequiresApprovalView: View {
    @State var openedSettings = false
    
    var body: some View {
        VStack(spacing: 8) {
            Text("Helper Not Enabled").font(.title).fontWeight(.bold)
            Text("Please enable DotLocal in the \"Allow in the Background\" section")
            Button(action: {
                print("previous status: \(HelperManager.shared.status)")
                HelperManager.shared.checkStatus()
                print("status: \(HelperManager.shared.status)")
                if HelperManager.shared.status == .requiresApproval {
                    openedSettings = true
                    SMAppService.openSystemSettingsLoginItems()
                }
            }, label: {
                if openedSettings {
                    Text("Continue")
                } else {
                    Text("Open System Settings")
                }
            })
        }.foregroundStyle(.secondary)
    }
}

//#Preview {
//    ContentView()
//}

#Preview {
    RequiresApprovalView()
}
