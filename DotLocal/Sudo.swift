// https://stackoverflow.com/a/69788418

import Foundation
import Security

public struct Sudo {

    private typealias AuthorizationExecuteWithPrivilegesImpl = @convention(c) (
        AuthorizationRef,
        UnsafePointer<CChar>, // path
        AuthorizationFlags,
        UnsafePointer<UnsafeMutablePointer<CChar>?>, // args
        UnsafeMutablePointer<UnsafeMutablePointer<FILE>>?
    ) -> OSStatus

    /// This wraps the deprecated AuthorizationExecuteWithPrivileges
    /// and makes it accessible by Swift
    ///
    /// - Parameters:
    ///   - path: The executable path
    ///   - arguments: The executable arguments
    /// - Returns: `errAuthorizationSuccess` or an error code
    public static func run(path: String, arguments: [String]) async -> (Bool, String) {
        var authRef: AuthorizationRef!
        var status = AuthorizationCreate(nil, nil, [], &authRef)

        guard status == errAuthorizationSuccess else { return (false, "") }
        defer { AuthorizationFree(authRef, [.destroyRights]) }

        var item = kAuthorizationRightExecute.withCString { name in
            AuthorizationItem(name: name, valueLength: 0, value: nil, flags: 0)
        }
        var rights = withUnsafeMutablePointer(to: &item) { ptr in
            AuthorizationRights(count: 1, items: ptr)
        }

        status = AuthorizationCopyRights(authRef, &rights, nil, [.interactionAllowed, .preAuthorize, .extendRights], nil)

        guard status == errAuthorizationSuccess else { return (false, "") }

        let (osStatus, stdout) = await executeWithPrivileges(authorization: authRef, path: path, arguments: arguments)

        return (status == errAuthorizationSuccess, stdout)
    }

    private static func executeWithPrivileges(authorization: AuthorizationRef,
                                              path: String,
                                              arguments: [String]) async -> (OSStatus, String) {
        let RTLD_DEFAULT = dlopen(nil, RTLD_NOW)
        guard let funcPtr = dlsym(RTLD_DEFAULT, "AuthorizationExecuteWithPrivileges") else { return (-1, "") }
        let args = arguments.map { strdup($0) }
        defer { args.forEach { free($0) }}
        let impl = unsafeBitCast(funcPtr, to: AuthorizationExecuteWithPrivilegesImpl.self)
        var communicationsPipe = UnsafeMutablePointer<FILE>.allocate(capacity: 1)
        let osStatus = impl(authorization, path, [], args, &communicationsPipe)
        let _file = communicationsPipe.pointee
        return await withCheckedContinuation { continuation in
            DispatchQueue.global().async {
                var fileContent = ""
                defer {
                    continuation.resume(returning: (osStatus, fileContent))
                }
                
                var file = _file
                guard file._read != nil else { return }
                
                let bufferSize = 1024
                var buffer = [UInt8](repeating: 0, count: bufferSize)

                // Read data from the file
                var bytesRead = fread(&buffer, 1, bufferSize, &file)

                var content = Data(buffer.prefix(bytesRead))

                // Continue reading until the end of the file
                while bytesRead > 0 {
                    bytesRead = fread(&buffer, 1, bufferSize, &file)
                    content.append(contentsOf: buffer.prefix(bytesRead))
                }

                // Convert the data to a string (adjust the encoding as needed)
                fileContent = String(data: content, encoding: .utf8) ?? ""
    //            let bufferSize = 1024
    //            print("1")
    //            guard let _read = file._read else {
    //                print("???")
    //                return
    //            }
    //            print("2")
    //            var buffer = [CChar](repeating: 0, count: Int(bufferSize))
    //            _read(&file, &buffer, Int32(bufferSize))
    //            print("got \(buffer)")
            }
        }
    }
}
