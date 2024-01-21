//
//  ProtoExtensions.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 6/1/2567 BE.
//

import Foundation
import SwiftProtobuf

protocol ProtoCodable {
    init(
        serializedData data: Data,
        extensions: ExtensionMap?,
        partial: Bool,
        options: BinaryDecodingOptions
    ) throws
    func serializedData(partial: Bool) throws -> Data
}
extension ProtoCodable {
    public init(from decoder: Swift.Decoder) throws {
        do {
            let container = try decoder.singleValueContainer()
            self = try Self(
                serializedData: try container.decode(Data.self),
                extensions: nil,
                partial: true,
                options: BinaryDecodingOptions()
            )
        } catch {
            print("error decoding: \(error)")
            throw error
        }
    }
    
    public func encode(to encoder: Encoder) throws {
        do {
            var container = encoder.singleValueContainer()
            try container.encode(try serializedData(partial: true))
        } catch {
            NSLog("error encoding: \(error)")
            throw error
        }
    }
}

extension SavedState: ProtoCodable, Codable {}
extension Mapping: ProtoCodable, Codable {}
extension Preferences: ProtoCodable, Codable {}

extension Mapping: Identifiable {}
extension Mapping: Comparable {
    public static func < (lhs: Mapping, rhs: Mapping) -> Bool {
        let lhsTitle = "\(lhs.host)\(lhs.pathPrefix)"
        let rhsTitle = "\(rhs.host)\(rhs.pathPrefix)"
        return lhsTitle < rhsTitle
    }
}
