//
// DO NOT EDIT.
// swift-format-ignore-file
//
// Generated by the protocol buffer compiler.
// Source: proto/dot-local.proto
//
import GRPC
import NIO
import NIOConcurrencyHelpers
import SwiftProtobuf


/// Usage: instantiate `DotLocalClient`, then call methods of this protocol to make API calls.
public protocol DotLocalClientProtocol: GRPCClient {
  var serviceName: String { get }
  var interceptors: DotLocalClientInterceptorFactoryProtocol? { get }

  func createMapping(
    _ request: CreateMappingRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<CreateMappingRequest, SwiftProtobuf.Google_Protobuf_Empty>

  func removeMapping(
    _ request: MappingKey,
    callOptions: CallOptions?
  ) -> UnaryCall<MappingKey, SwiftProtobuf.Google_Protobuf_Empty>

  func listMappings(
    _ request: SwiftProtobuf.Google_Protobuf_Empty,
    callOptions: CallOptions?
  ) -> UnaryCall<SwiftProtobuf.Google_Protobuf_Empty, ListMappingsResponse>
}

extension DotLocalClientProtocol {
  public var serviceName: String {
    return "DotLocal"
  }

  /// Unary call to CreateMapping
  ///
  /// - Parameters:
  ///   - request: Request to send to CreateMapping.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  public func createMapping(
    _ request: CreateMappingRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<CreateMappingRequest, SwiftProtobuf.Google_Protobuf_Empty> {
    return self.makeUnaryCall(
      path: DotLocalClientMetadata.Methods.createMapping.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeCreateMappingInterceptors() ?? []
    )
  }

  /// Unary call to RemoveMapping
  ///
  /// - Parameters:
  ///   - request: Request to send to RemoveMapping.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  public func removeMapping(
    _ request: MappingKey,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<MappingKey, SwiftProtobuf.Google_Protobuf_Empty> {
    return self.makeUnaryCall(
      path: DotLocalClientMetadata.Methods.removeMapping.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeRemoveMappingInterceptors() ?? []
    )
  }

  /// Unary call to ListMappings
  ///
  /// - Parameters:
  ///   - request: Request to send to ListMappings.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  public func listMappings(
    _ request: SwiftProtobuf.Google_Protobuf_Empty,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<SwiftProtobuf.Google_Protobuf_Empty, ListMappingsResponse> {
    return self.makeUnaryCall(
      path: DotLocalClientMetadata.Methods.listMappings.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeListMappingsInterceptors() ?? []
    )
  }
}

@available(*, deprecated)
extension DotLocalClient: @unchecked Sendable {}

@available(*, deprecated, renamed: "DotLocalNIOClient")
public final class DotLocalClient: DotLocalClientProtocol {
  private let lock = Lock()
  private var _defaultCallOptions: CallOptions
  private var _interceptors: DotLocalClientInterceptorFactoryProtocol?
  public let channel: GRPCChannel
  public var defaultCallOptions: CallOptions {
    get { self.lock.withLock { return self._defaultCallOptions } }
    set { self.lock.withLockVoid { self._defaultCallOptions = newValue } }
  }
  public var interceptors: DotLocalClientInterceptorFactoryProtocol? {
    get { self.lock.withLock { return self._interceptors } }
    set { self.lock.withLockVoid { self._interceptors = newValue } }
  }

  /// Creates a client for the DotLocal service.
  ///
  /// - Parameters:
  ///   - channel: `GRPCChannel` to the service host.
  ///   - defaultCallOptions: Options to use for each service call if the user doesn't provide them.
  ///   - interceptors: A factory providing interceptors for each RPC.
  public init(
    channel: GRPCChannel,
    defaultCallOptions: CallOptions = CallOptions(),
    interceptors: DotLocalClientInterceptorFactoryProtocol? = nil
  ) {
    self.channel = channel
    self._defaultCallOptions = defaultCallOptions
    self._interceptors = interceptors
  }
}

public struct DotLocalNIOClient: DotLocalClientProtocol {
  public var channel: GRPCChannel
  public var defaultCallOptions: CallOptions
  public var interceptors: DotLocalClientInterceptorFactoryProtocol?

  /// Creates a client for the DotLocal service.
  ///
  /// - Parameters:
  ///   - channel: `GRPCChannel` to the service host.
  ///   - defaultCallOptions: Options to use for each service call if the user doesn't provide them.
  ///   - interceptors: A factory providing interceptors for each RPC.
  public init(
    channel: GRPCChannel,
    defaultCallOptions: CallOptions = CallOptions(),
    interceptors: DotLocalClientInterceptorFactoryProtocol? = nil
  ) {
    self.channel = channel
    self.defaultCallOptions = defaultCallOptions
    self.interceptors = interceptors
  }
}

@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
public protocol DotLocalAsyncClientProtocol: GRPCClient {
  static var serviceDescriptor: GRPCServiceDescriptor { get }
  var interceptors: DotLocalClientInterceptorFactoryProtocol? { get }

  func makeCreateMappingCall(
    _ request: CreateMappingRequest,
    callOptions: CallOptions?
  ) -> GRPCAsyncUnaryCall<CreateMappingRequest, SwiftProtobuf.Google_Protobuf_Empty>

  func makeRemoveMappingCall(
    _ request: MappingKey,
    callOptions: CallOptions?
  ) -> GRPCAsyncUnaryCall<MappingKey, SwiftProtobuf.Google_Protobuf_Empty>

  func makeListMappingsCall(
    _ request: SwiftProtobuf.Google_Protobuf_Empty,
    callOptions: CallOptions?
  ) -> GRPCAsyncUnaryCall<SwiftProtobuf.Google_Protobuf_Empty, ListMappingsResponse>
}

@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
extension DotLocalAsyncClientProtocol {
  public static var serviceDescriptor: GRPCServiceDescriptor {
    return DotLocalClientMetadata.serviceDescriptor
  }

  public var interceptors: DotLocalClientInterceptorFactoryProtocol? {
    return nil
  }

  public func makeCreateMappingCall(
    _ request: CreateMappingRequest,
    callOptions: CallOptions? = nil
  ) -> GRPCAsyncUnaryCall<CreateMappingRequest, SwiftProtobuf.Google_Protobuf_Empty> {
    return self.makeAsyncUnaryCall(
      path: DotLocalClientMetadata.Methods.createMapping.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeCreateMappingInterceptors() ?? []
    )
  }

  public func makeRemoveMappingCall(
    _ request: MappingKey,
    callOptions: CallOptions? = nil
  ) -> GRPCAsyncUnaryCall<MappingKey, SwiftProtobuf.Google_Protobuf_Empty> {
    return self.makeAsyncUnaryCall(
      path: DotLocalClientMetadata.Methods.removeMapping.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeRemoveMappingInterceptors() ?? []
    )
  }

  public func makeListMappingsCall(
    _ request: SwiftProtobuf.Google_Protobuf_Empty,
    callOptions: CallOptions? = nil
  ) -> GRPCAsyncUnaryCall<SwiftProtobuf.Google_Protobuf_Empty, ListMappingsResponse> {
    return self.makeAsyncUnaryCall(
      path: DotLocalClientMetadata.Methods.listMappings.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeListMappingsInterceptors() ?? []
    )
  }
}

@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
extension DotLocalAsyncClientProtocol {
  public func createMapping(
    _ request: CreateMappingRequest,
    callOptions: CallOptions? = nil
  ) async throws -> SwiftProtobuf.Google_Protobuf_Empty {
    return try await self.performAsyncUnaryCall(
      path: DotLocalClientMetadata.Methods.createMapping.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeCreateMappingInterceptors() ?? []
    )
  }

  public func removeMapping(
    _ request: MappingKey,
    callOptions: CallOptions? = nil
  ) async throws -> SwiftProtobuf.Google_Protobuf_Empty {
    return try await self.performAsyncUnaryCall(
      path: DotLocalClientMetadata.Methods.removeMapping.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeRemoveMappingInterceptors() ?? []
    )
  }

  public func listMappings(
    _ request: SwiftProtobuf.Google_Protobuf_Empty,
    callOptions: CallOptions? = nil
  ) async throws -> ListMappingsResponse {
    return try await self.performAsyncUnaryCall(
      path: DotLocalClientMetadata.Methods.listMappings.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeListMappingsInterceptors() ?? []
    )
  }
}

@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
public struct DotLocalAsyncClient: DotLocalAsyncClientProtocol {
  public var channel: GRPCChannel
  public var defaultCallOptions: CallOptions
  public var interceptors: DotLocalClientInterceptorFactoryProtocol?

  public init(
    channel: GRPCChannel,
    defaultCallOptions: CallOptions = CallOptions(),
    interceptors: DotLocalClientInterceptorFactoryProtocol? = nil
  ) {
    self.channel = channel
    self.defaultCallOptions = defaultCallOptions
    self.interceptors = interceptors
  }
}

public protocol DotLocalClientInterceptorFactoryProtocol: Sendable {

  /// - Returns: Interceptors to use when invoking 'createMapping'.
  func makeCreateMappingInterceptors() -> [ClientInterceptor<CreateMappingRequest, SwiftProtobuf.Google_Protobuf_Empty>]

  /// - Returns: Interceptors to use when invoking 'removeMapping'.
  func makeRemoveMappingInterceptors() -> [ClientInterceptor<MappingKey, SwiftProtobuf.Google_Protobuf_Empty>]

  /// - Returns: Interceptors to use when invoking 'listMappings'.
  func makeListMappingsInterceptors() -> [ClientInterceptor<SwiftProtobuf.Google_Protobuf_Empty, ListMappingsResponse>]
}

public enum DotLocalClientMetadata {
  public static let serviceDescriptor = GRPCServiceDescriptor(
    name: "DotLocal",
    fullName: "DotLocal",
    methods: [
      DotLocalClientMetadata.Methods.createMapping,
      DotLocalClientMetadata.Methods.removeMapping,
      DotLocalClientMetadata.Methods.listMappings,
    ]
  )

  public enum Methods {
    public static let createMapping = GRPCMethodDescriptor(
      name: "CreateMapping",
      path: "/DotLocal/CreateMapping",
      type: GRPCCallType.unary
    )

    public static let removeMapping = GRPCMethodDescriptor(
      name: "RemoveMapping",
      path: "/DotLocal/RemoveMapping",
      type: GRPCCallType.unary
    )

    public static let listMappings = GRPCMethodDescriptor(
      name: "ListMappings",
      path: "/DotLocal/ListMappings",
      type: GRPCCallType.unary
    )
  }
}

/// To build a server, implement a class that conforms to this protocol.
public protocol DotLocalProvider: CallHandlerProvider {
  var interceptors: DotLocalServerInterceptorFactoryProtocol? { get }

  func createMapping(request: CreateMappingRequest, context: StatusOnlyCallContext) -> EventLoopFuture<SwiftProtobuf.Google_Protobuf_Empty>

  func removeMapping(request: MappingKey, context: StatusOnlyCallContext) -> EventLoopFuture<SwiftProtobuf.Google_Protobuf_Empty>

  func listMappings(request: SwiftProtobuf.Google_Protobuf_Empty, context: StatusOnlyCallContext) -> EventLoopFuture<ListMappingsResponse>
}

extension DotLocalProvider {
  public var serviceName: Substring {
    return DotLocalServerMetadata.serviceDescriptor.fullName[...]
  }

  /// Determines, calls and returns the appropriate request handler, depending on the request's method.
  /// Returns nil for methods not handled by this service.
  public func handle(
    method name: Substring,
    context: CallHandlerContext
  ) -> GRPCServerHandlerProtocol? {
    switch name {
    case "CreateMapping":
      return UnaryServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<CreateMappingRequest>(),
        responseSerializer: ProtobufSerializer<SwiftProtobuf.Google_Protobuf_Empty>(),
        interceptors: self.interceptors?.makeCreateMappingInterceptors() ?? [],
        userFunction: self.createMapping(request:context:)
      )

    case "RemoveMapping":
      return UnaryServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<MappingKey>(),
        responseSerializer: ProtobufSerializer<SwiftProtobuf.Google_Protobuf_Empty>(),
        interceptors: self.interceptors?.makeRemoveMappingInterceptors() ?? [],
        userFunction: self.removeMapping(request:context:)
      )

    case "ListMappings":
      return UnaryServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<SwiftProtobuf.Google_Protobuf_Empty>(),
        responseSerializer: ProtobufSerializer<ListMappingsResponse>(),
        interceptors: self.interceptors?.makeListMappingsInterceptors() ?? [],
        userFunction: self.listMappings(request:context:)
      )

    default:
      return nil
    }
  }
}

/// To implement a server, implement an object which conforms to this protocol.
@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
public protocol DotLocalAsyncProvider: CallHandlerProvider, Sendable {
  static var serviceDescriptor: GRPCServiceDescriptor { get }
  var interceptors: DotLocalServerInterceptorFactoryProtocol? { get }

  func createMapping(
    request: CreateMappingRequest,
    context: GRPCAsyncServerCallContext
  ) async throws -> SwiftProtobuf.Google_Protobuf_Empty

  func removeMapping(
    request: MappingKey,
    context: GRPCAsyncServerCallContext
  ) async throws -> SwiftProtobuf.Google_Protobuf_Empty

  func listMappings(
    request: SwiftProtobuf.Google_Protobuf_Empty,
    context: GRPCAsyncServerCallContext
  ) async throws -> ListMappingsResponse
}

@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
extension DotLocalAsyncProvider {
  public static var serviceDescriptor: GRPCServiceDescriptor {
    return DotLocalServerMetadata.serviceDescriptor
  }

  public var serviceName: Substring {
    return DotLocalServerMetadata.serviceDescriptor.fullName[...]
  }

  public var interceptors: DotLocalServerInterceptorFactoryProtocol? {
    return nil
  }

  public func handle(
    method name: Substring,
    context: CallHandlerContext
  ) -> GRPCServerHandlerProtocol? {
    switch name {
    case "CreateMapping":
      return GRPCAsyncServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<CreateMappingRequest>(),
        responseSerializer: ProtobufSerializer<SwiftProtobuf.Google_Protobuf_Empty>(),
        interceptors: self.interceptors?.makeCreateMappingInterceptors() ?? [],
        wrapping: { try await self.createMapping(request: $0, context: $1) }
      )

    case "RemoveMapping":
      return GRPCAsyncServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<MappingKey>(),
        responseSerializer: ProtobufSerializer<SwiftProtobuf.Google_Protobuf_Empty>(),
        interceptors: self.interceptors?.makeRemoveMappingInterceptors() ?? [],
        wrapping: { try await self.removeMapping(request: $0, context: $1) }
      )

    case "ListMappings":
      return GRPCAsyncServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<SwiftProtobuf.Google_Protobuf_Empty>(),
        responseSerializer: ProtobufSerializer<ListMappingsResponse>(),
        interceptors: self.interceptors?.makeListMappingsInterceptors() ?? [],
        wrapping: { try await self.listMappings(request: $0, context: $1) }
      )

    default:
      return nil
    }
  }
}

public protocol DotLocalServerInterceptorFactoryProtocol: Sendable {

  /// - Returns: Interceptors to use when handling 'createMapping'.
  ///   Defaults to calling `self.makeInterceptors()`.
  func makeCreateMappingInterceptors() -> [ServerInterceptor<CreateMappingRequest, SwiftProtobuf.Google_Protobuf_Empty>]

  /// - Returns: Interceptors to use when handling 'removeMapping'.
  ///   Defaults to calling `self.makeInterceptors()`.
  func makeRemoveMappingInterceptors() -> [ServerInterceptor<MappingKey, SwiftProtobuf.Google_Protobuf_Empty>]

  /// - Returns: Interceptors to use when handling 'listMappings'.
  ///   Defaults to calling `self.makeInterceptors()`.
  func makeListMappingsInterceptors() -> [ServerInterceptor<SwiftProtobuf.Google_Protobuf_Empty, ListMappingsResponse>]
}

public enum DotLocalServerMetadata {
  public static let serviceDescriptor = GRPCServiceDescriptor(
    name: "DotLocal",
    fullName: "DotLocal",
    methods: [
      DotLocalServerMetadata.Methods.createMapping,
      DotLocalServerMetadata.Methods.removeMapping,
      DotLocalServerMetadata.Methods.listMappings,
    ]
  )

  public enum Methods {
    public static let createMapping = GRPCMethodDescriptor(
      name: "CreateMapping",
      path: "/DotLocal/CreateMapping",
      type: GRPCCallType.unary
    )

    public static let removeMapping = GRPCMethodDescriptor(
      name: "RemoveMapping",
      path: "/DotLocal/RemoveMapping",
      type: GRPCCallType.unary
    )

    public static let listMappings = GRPCMethodDescriptor(
      name: "ListMappings",
      path: "/DotLocal/ListMappings",
      type: GRPCCallType.unary
    )
  }
}
