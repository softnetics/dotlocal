//
//  MappingList.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import SwiftUI

struct MappingList: View {
    @StateObject var clientManager = ClientManager.shared
    @StateObject var daemonManager = DaemonManager.shared
    
    var body: some View {
        let mappings = daemonManager.savedState.mappings
        List(mappings) { mapping in
            HStack(spacing: 12) {
                VStack(alignment: .leading, spacing: 4) {
                    Text("\(mapping.host)\(mapping.pathPrefix)")
                    Text("\(mapping.target)").foregroundStyle(.secondary)
                }
                Spacer()
                let href = "http://\(mapping.host)\(mapping.pathPrefix)"
                Link(destination: URL(string: href)!) {
                    Image(systemName: "link")
                }
                .buttonStyle(.borderless)
            }
            .padding(.vertical, 4)
        }
        .if(!mappings.isEmpty) {
            if mappings.count > 1 {
                $0.navigationSubtitle("\(mappings.count) routes")
            } else {
                $0.navigationSubtitle("1 route")
            }
        }
        .overlay {
            if mappings.isEmpty {
                if #available(macOS 14.0, *) {
                    ContentUnavailableView {
                        Label("No Routes", systemImage: "arrow.triangle.swap")
                    } description: {
                        hintView()
                    }
                } else {
                    VStack(spacing: 8) {
                        Image(systemName: "arrow.triangle.swap").font(.system(size: 48)).foregroundStyle(.tertiary).padding(.bottom, 8)
                        Text("No Routes").font(.title).fontWeight(.bold)
                        hintView()
                    }
                    .foregroundStyle(.secondary)
                }
            }
        }
    }
    
    private func getMappings(state: DaemonState) -> [Mapping] {
        if case .started(let savedState) = state {
            return savedState.mappings
        }
        return []
    }
    
    @ViewBuilder
    private func hintView() -> some View {
        VStack(spacing: 4) {
            if clientManager.installed {
                Text("Try creating a new route")
                    .font(.system(size: 18, weight: .bold))
                    .bold()
                Text("dotlocal --host test.local pnpm dev").monospaced()
            } else {
                Text("Install \(getDotLocalLabel()) Command Line Tool to create routes")
                Button(action: {
                    Task {
                        await clientManager.installCli()
                    }
                }, label: {
                    Text("Install")
                })
            }
        }
        .onAppear {
            clientManager.checkInstalled()
        }
    }
    
    private func getDotLocalLabel() -> AttributedString {
        var label = AttributedString("dotlocal")
        label.font = .body.monospaced()
        return label
    }
}
