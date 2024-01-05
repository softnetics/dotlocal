//
//  MappingListViewModel.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import Foundation

@MainActor class MappingListViewModel: ObservableObject {
    @Published var loading = true
    @Published var mappings = [Mapping]()
    
    func fetchMappings() async {
        guard let apiClient = DaemonManager.shared.apiClient else {return}
        do {
            let res = try await apiClient.listMappings(.with({_ in }))
            self.mappings = res.mappings.sorted()
        } catch {
            print(error)
        }
        loading = false
    }
}
