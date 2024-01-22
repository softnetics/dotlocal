//
//  DaemonState.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 12/1/2567 BE.
//

import Foundation

enum DaemonState: Codable {
    case unknown
    case stopped
    case starting
    case started(savedState: SavedState)
}
