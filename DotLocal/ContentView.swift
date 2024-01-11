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
    
    var body: some View {
        VStack {
            Button(action: {
                let service = SMAppService.daemon(plistName: "helper.plist")
                do {
                    print("status: \(service.status.rawValue)")
                    if service.status == .enabled {
                        print("will unregister")
                        try service.unregister()
                        print("did unregister")
                    } else {
                        print("will register")
                        try service.register()
                        print("did register")
                    }
                } catch {
                    print("error: \(error)")
                }
            }, label: {
                Text("Test")
            })
            switch daemonManager.state {
            case .stopped:
                Text("DotLocal is not running")
            case .starting:
                ProgressView()
            case .started:
                MappingList()
            }
        }
        .frame(maxWidth: /*@START_MENU_TOKEN@*/.infinity/*@END_MENU_TOKEN@*/, maxHeight: .infinity)
        .toolbar() {
            StartStopButton(state: daemonManager.state, onStart: { daemonManager.start() }, onStop: { daemonManager.stop() })
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
        case .starting:
            ProgressView().controlSize(.small)
        case .started:
            Button(action: onStop) {
                Label("Stop", systemImage: "stop.fill")
            }
        }
    }
}

#Preview {
    ContentView()
}
