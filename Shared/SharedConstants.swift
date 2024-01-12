//
//  SharedConstants.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 11/1/2567 BE.
//

import Foundation
import SecureXPC

struct SharedConstants {
    static let installClientRoute = XPCRoute.named("installClient")
    static let uninstallClientRoute = XPCRoute.named("uninstallClient")
    static let exitRoute = XPCRoute.named("exit")
}
