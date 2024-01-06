//
//  MappingList.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import SwiftUI

struct MappingList: View {
    @StateObject var vm = MappingListViewModel()
    
    var body: some View {
        List(vm.mappings) { mapping in
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
        .overlay {
            if vm.loading {
                ProgressView()
            } else if vm.mappings.isEmpty {
                ContentUnavailableView {
                    Label("No Routes", systemImage: "arrow.triangle.swap")
                } description: {
                    VStack(spacing: 4) {
                        Text("Try creating a new route")
                            .font(.system(size: 18, weight: .bold))
                            .bold()
                        Text("dotlocal -n test.local pnpm dev").monospaced()
                    }
                }
            }
        }
    }
}
