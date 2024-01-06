//
//  MappingListViewModel.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import Foundation
import Combine

@MainActor class MappingListViewModel: ObservableObject {
    @Published var loading = true
    @Published var mappings = [Mapping]()
    
    private var subscriptions = Set<AnyCancellable>()
    
    init() {
        DaemonManager.shared.updates()
            .flatMap { _ in
                Future<[Mapping], Never> { promise in
                    Task {
                        guard let apiClient = DaemonManager.shared.apiClient else {
                            promise(.success([]))
                            return
                        }
                        let res = try await apiClient.listMappings(.with({_ in}))
                        promise(.success(res.mappings.sorted()))
                    }
                }
            }.sink { mappings in
                self.loading = false
                self.mappings = mappings
            }.store(in: &subscriptions)
    }
}
