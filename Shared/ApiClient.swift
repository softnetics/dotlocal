//
//  ApiClient.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 21/1/2567 BE.
//

import Foundation
import GRPC
import NIO

func createApiClient(group: EventLoopGroup) throws -> DotLocalAsyncClient {
    let socketPath = URL.init(filePath: "/var/run/dotlocal/api.sock")
    // TODO: try catch
    let channel = try GRPCChannelPool.with(
        target: .unixDomainSocket(socketPath.path(percentEncoded: false)),
        transportSecurity: .plaintext,
        eventLoopGroup: group
    )
    return DotLocalAsyncClient(channel: channel)
}
