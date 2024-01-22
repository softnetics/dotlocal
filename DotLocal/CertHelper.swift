//
//  CertHelper.swift
//  DotLocal
//
//  Created by Suphon Thanakornpakapong on 21/1/2567 BE.
//

import Foundation
import SecurityInterface
import SwiftProtobuf

struct CertHelper {
    static func getRootCertificate() async throws -> CertHelper.Certificate {
        let res = try await DaemonManager.shared.apiClient.getRootCertificate(Google_Protobuf_Empty())
        return try await CertHelper.Certificate(res: res)
    }
    
    static func rootCertificateLogo() -> NSImage {
        let bundle = Bundle(for: SFCertificateView.self)
        return bundle.image(forResource: "CertLargeRoot")!
    }
    
    struct Certificate {
        let secCertificate: SecCertificate
        let commonName: String
        let notBefore: Date
        let notAfter: Date
        let trusted: Bool
        
        init(res: GetRootCertificateResponse) async throws {
            guard let certificate = SecCertificateCreateWithData(nil, res.certificate as CFData) else {
                throw CertHelperError.invalidCertificate
            }
            secCertificate = certificate
            let tmp = UnsafeMutablePointer<CFString?>.allocate(capacity: 1)
            SecCertificateCopyCommonName(certificate, tmp)
            commonName = tmp.pointee! as String
            notBefore = res.notBefore.date
            notAfter = res.notAfter.date
            trusted = try await evaluateTrust(certificate: secCertificate)
        }
    }
}

enum CertHelperError: Error {
    case invalidCertificate
}

fileprivate func evaluateTrust(certificate: SecCertificate) async throws -> Bool {
    var secTrust: SecTrust?
    if SecTrustCreateWithCertificates(certificate, SecPolicyCreateBasicX509(), &secTrust) == errSecSuccess, let trust = secTrust {
        let error = UnsafeMutablePointer<CFError?>.allocate(capacity: 1)
        let result = SecTrustEvaluateWithError(trust, error)
        if error.pointee != nil {
            return false
        }
        return result
    } else {
        return false
    }
}
