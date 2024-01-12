//
//  Utils.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 12/1/2567 BE.
//

import Foundation

func parentAppURL() throws -> URL {
    return try parentAppURL(Bundle.main.bundleURL.pathComponents)
}

func parentAppURL(_ bundleURL: URL) throws -> URL {
    return try parentAppURL(bundleURL.pathComponents)
}

func parentAppURL(_ components: [String]) throws -> URL {
    guard let contentsIndex = components.lastIndex(of: "Contents"),
          components[components.index(before: contentsIndex)].hasSuffix(".app") else {
        throw MyError.runtimeError("""
        Parent bundle could not be found.
        Path:\(Bundle.main.bundleURL)
        """)
    }
    
    return URL(fileURLWithPath: "/" + components[1..<contentsIndex].joined(separator: "/"))
}
