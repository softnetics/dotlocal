//
//  ProtoExtensions.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import Foundation

extension Mapping: Identifiable {}

extension Mapping: Decodable {
    public init(from decoder: Decoder) throws {
        do {
            let container = try decoder.singleValueContainer()
            self = try Mapping(serializedData: try container.decode(Data.self))
        } catch {
            print("error decoding: \(error)")
            throw error
        }
    }
}

extension Mapping: Encodable {
    public func encode(to encoder: Encoder) throws {
        do {
            var container = encoder.singleValueContainer()
            try container.encode(try serializedData())
        } catch {
            NSLog("error encoding: \(error)")
            throw error
        }
    }
}

extension Mapping: Comparable {
    public static func < (lhs: Mapping, rhs: Mapping) -> Bool {
        let lhsTitle = "\(lhs.host)\(lhs.pathPrefix)"
        let rhsTitle = "\(rhs.host)\(rhs.pathPrefix)"
        return lhsTitle < rhsTitle
    }
}
