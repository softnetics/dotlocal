//
//  ProtoExtensions.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import Foundation

extension Mapping: Identifiable {}

extension Mapping: Comparable {
    public static func < (lhs: Mapping, rhs: Mapping) -> Bool {
        let lhsTitle = "\(lhs.host)\(lhs.pathPrefix)"
        let rhsTitle = "\(rhs.host)\(rhs.pathPrefix)"
        return lhsTitle < rhsTitle
    }
}
