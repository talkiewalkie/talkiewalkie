//
//  GrpcClient.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 12/10/2021.
//

import Combine
import CoreData
import Foundation
import GRPC
import NIO
import OSLog

enum GrpcConnectionState {
    case Disconnected
    case Connecting
    case Connected
}

class GrpcConnectivityState: ConnectivityStateDelegate, ObservableObject {
    private let url: URL
    var onReconnection: () -> Void

    private let logger = Logger.withLabel("grpc-status")
    @Published var state = GrpcConnectionState.Disconnected

    init(_ url: URL, onReconnection: @escaping () -> Void) {
        self.url = url
        self.onReconnection = onReconnection
    }

    func connectivityStateDidChange(from oldState: ConnectivityState, to newState: ConnectivityState) {
        if oldState != .ready, newState == .ready {
            logger.debug("got connected [\(self.url.absoluteString)]")
            state = .Connected
            onReconnection()
        } else if oldState == .ready, newState != .ready {
            logger.debug("got disconnected (\(String(describing: newState))) [\(self.url.absoluteString)]")
            state = .Disconnected
        }
        if newState == .connecting {
            logger.debug("connecting [\(self.url.absoluteString)]")
            state = .Connecting
        }
        if newState == .transientFailure {
            logger.debug("error (transient failure) [\(self.url.absoluteString)]")
            state = .Disconnected
        }
    }
}

private extension String {
    func toCallOption() -> CallOptions {
        return CallOptions(customMetadata: ["Authorization": "Bearer \(self)"], timeLimit: .timeout(.seconds(3)))
    }

    func toStreamingCallOption(_ timeout: TimeAmount) -> CallOptions {
        return CallOptions(customMetadata: ["Authorization": "Bearer \(self)"], timeLimit: .timeout(timeout))
    }
}

class AuthedGrpcApi {
    let stateDelegate: GrpcConnectivityState

    private let url: URL
    private var stream: BidirectionalStreamingCall<App_Event, App_Event>?

    private let logger = Logger.withLabel("grpc-client")
    private let empty = App_Empty()
    private let eventQueue = [App_Event]()

    private let group: EventLoopGroup
    private let channel: ClientConnection
    private let token: String

    private let userClient: App_UserServiceClient
    private let convClient: App_ConversationServiceClient
    private let eventClient: App_EventServiceClient

    static func withUrlAndToken(url: URL, token: String, writer: @escaping (_: @escaping (_ context: NSManagedObjectContext, _ me: Me?) -> Void) -> Void) -> AuthedGrpcApi {
        let api = AuthedGrpcApi(url: url, token: token)
        api.stateDelegate.onReconnection = {
            let (down, _) = api.sync()
            if let down = down {
                down.events.forEach { e in writer { ctx, _ in LoadEventToCoreData(e, ctx: ctx) } }
            }

            api.logger.debug("reconnecting to stream")
            api.connect { newEvent in
                api.logger.debug("handling new message from stream")
                writer { ctx, _ in LoadEventToCoreData(newEvent, ctx: ctx) }
            }
        }

        return api
    }

    private init(url: URL, token: String) {
        self.url = url
        self.token = token
        stateDelegate = GrpcConnectivityState(url, onReconnection: {})

        group = PlatformSupport.makeEventLoopGroup(loopCount: 1)

        #if DEBUG
            let channelBuilder = ClientConnection.insecure(group: group)
        #else
            let channelBuilder = ClientConnection.usingPlatformAppropriateTLS(for: group)
        #endif
        channel = channelBuilder
            .withConnectionReestablishment(enabled: true)
            .withCallStartBehavior(.fastFailure)
            .withConnectionBackoff(initial: .seconds(1))
            .withConnectionBackoff(maximum: .seconds(30))
            .withConnectionBackoff(multiplier: 1)
            .withConnectionTimeout(minimum: .seconds(5))
            .withConnectivityStateDelegate(stateDelegate, executingOn: DispatchQueue.global(qos: .background))
            .connect(host: url.host!, port: url.port!)

        userClient = App_UserServiceClient(channel: channel)
        userClient.defaultCallOptions = token.toCallOption()

        convClient = App_ConversationServiceClient(channel: channel)
        convClient.defaultCallOptions = token.toCallOption()

        eventClient = App_EventServiceClient(channel: channel)
        eventClient.defaultCallOptions = token.toCallOption()
    }

    deinit {
        do {
            try? channel.close().wait()
            try group.syncShutdownGracefully()
        } catch { logger.error("could not end connection: \(error.localizedDescription)") }
    }

    func userByUuid(_ uuid: UUID) -> (App_User?, Error?) {
        let input = App_UserGetInput.with { $0.uuid = uuid.uuidString }

        return userClient.get(input).waitForOutput()
    }

    func me() -> (App_MeUser?, Error?) {
        return userClient.me(empty).waitForOutput()
    }

    func onboardingComplete(displayName: String, locales: [String]) -> (App_MeUser?, Error?) {
        let input = App_OnboardingInput.with {
            $0.displayName = displayName
            $0.locales = locales
        }

        return userClient.onboarding(input).waitForOutput()
    }

    func syncContactList(phones: [String]) -> (App_SyncContactsOutput?, Error?) {
        let input = App_SyncContactsInput.with { $0.phoneNumbers = phones }

        return userClient.syncContacts(input).waitForOutput()
    }

    func listConvs() -> ([App_Conversation], Error?) {
        let input = App_ConversationListInput.with { $0.page = 0 }
        var convs: [App_Conversation] = []
        let call = convClient.list(input) { convs.append($0) }

        let error = call.waitCompletion()
        return (convs, error)
    }

    func convByUuid(_ uuid: UUID) -> (App_Conversation?, Error?) {
        let input = App_ConversationGetInput.with { $0.uuid = uuid.uuidString }

        return convClient.get(input).waitForOutput()
    }

    func sync() -> (App_DownSync?, Error?) {
        let input = App_UpSync.with { this in
            this.events = eventQueue
            this.lastEventUuid = UserDefaults.standard.value(forKey: "lastEventUuid") as! String
        }

        let (down, error) = eventClient.sync(input).waitForOutput()
        if let down = down {
            UserDefaults.standard.set(down.lastEventUuid, forKey: "lastEventUuid")
        }

        return (down, error)
    }

    func connect(completion: @escaping (App_Event) -> Void) {
        let stream = eventClient.connect(callOptions: token.toStreamingCallOption(.hours(1))) { newEvent in
            UserDefaults.standard.set(newEvent.uuid, forKey: "lastEventUuid")
            completion(newEvent)
        }

        self.stream = stream
    }

    func sendMessage(text: String, convUuid: UUID) -> App_Event {
        let message = App_Event.with { this in
            this.localUuid = UUID().uuidString
            this.content = .sentNewMessage(App_Event.SentNewMessage.with { msg in
                msg.message = App_MessageSendInput.with { input in
                    input.content = App_MessageSendInput.OneOf_Content.textMessage(App_TextMessage.with { $0.content = text })
                }

                msg.conversation = .convUuid(convUuid.uuidString)
            })
        }
        stream?.sendMessage(message)
            .whenFailure { err in
                self.logger.error("could not send message: \(err.localizedDescription)")
            }
        return message
    }
}

extension UnaryCall {
    @discardableResult
    func waitForOutput() -> (ResponsePayload?, Error?) {
        os_log(.debug, "[grpc] \(path) waiting")

        var msg = "[grpc] \(path)"
        do {
            let st = try status.wait()
            msg += " status:(\(st.code.description))"
            if let m = st.message, !m.isEmpty { msg += " message:(\(m))" }
        } catch {
            os_log(.error, "\(error.localizedDescription)")
            return (nil, error)
        }

        do {
            let res = try response.wait()
            os_log(.debug, "\(msg)")
            return (res, nil)
        } catch {
            os_log(.error, "\(msg) failed with: \(error.localizedDescription)")
            return (nil, error)
        }
    }
}

extension ServerStreamingCall {
    func waitCompletion() -> Error? {
        os_log(.debug, "[grpc] \(path) waiting")
        var error: Error?
        var st: GRPCStatus?

        _ = try? status
            .map { s in st = s }
            .recover { err in
                error = err
                st = .processingError
            }
            .wait()

        if let code = st?.code, let msg = st?.message, code != .ok {
            let errMsg: String = error?.localizedDescription ?? "[did not catch error]"
            os_log(.debug, "[grpc] \(path) returned with status \(code.description), \(msg) error: \(errMsg)")
            return error
        }

        os_log(.debug, "[grpc] \(path) completed")
        return nil
    }
}

func LoadEventToCoreData(_ event: App_Event, ctx: NSManagedObjectContext) {
    switch event.content {
    case .some(.receivedNewMessage(_)):
        Message.fromEventProto(event, context: ctx)

    case .some(.deletedMessage(_)):
        break

    case .some(.changedPicture(_)):
        break

    case .some(.joinedConversation(_)):
        break

    case .some(.leftConversation(_)):
        break

    case .some(.conversationTitleChanged(_)):
        break

    case .some(.sentNewMessage(_)):
        fatalError()

    case .none:
        fatalError()
    }
}
