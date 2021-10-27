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
    private let logger = Logger.withLabel("grpc-status")
    @Published var state = GrpcConnectionState.Disconnected
    
    init(_ url: URL) {
        self.url = url
    }

    func connectivityStateDidChange(from oldState: ConnectivityState, to newState: ConnectivityState) {
        if oldState != .ready, newState == .ready {
            logger.debug("got connected [\(self.url.absoluteString)]")
            state = .Connected
        }
        else if oldState == .ready, newState != .ready {
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

    private let logger = Logger.withLabel("grpc-client")
    private let empty = App_Empty()

    private let group: EventLoopGroup
    private let channel: ClientConnection
    private let token: String
    private let userClient: App_UserServiceClient
    private let convClient: App_ConversationServiceClient
    private let mssgClient: App_MessageServiceClient
    private let persistentContainer: NSPersistentContainer

    init(url: URL, token: String, persistentContainer: NSPersistentContainer) {
        self.url = url
        self.token = token
        self.persistentContainer = persistentContainer
        stateDelegate = GrpcConnectivityState(url)

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

        mssgClient = App_MessageServiceClient(channel: channel)
        mssgClient.defaultCallOptions = token.toCallOption()
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

        let (twCl, error) = userClient.syncContacts(input).waitForOutput()
        if let twCl = twCl {
            persistentContainer.performBackgroundTask { context in
                twCl.users.forEach { u in User.upsert(u, context: context) }
                context.saveOrLogError()
            }
        }

        return (twCl, error)
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

    func sendMessage(text: String, convUuid: UUID) {
        let input = App_MessageSendInput.with {
            $0.content = App_MessageSendInput.OneOf_Content.textMessage(App_TextMessage.with { tm in tm.content = text })
            $0.recipients = App_MessageSendInput.OneOf_Recipients.convUuid(convUuid.uuidString)
        }
        mssgClient.send(input).waitForOutput()
    }

    func subscribeIncomingMessages(completion: @escaping (App_Message) -> Void) {
        let call = mssgClient.incoming(empty, callOptions: token.toStreamingCallOption(.hours(1)), handler: completion)
        _ = call.status.recover { err in
            self.logger.error("received error from incoming messages stream: \(err.localizedDescription)")

            return .processingError
        }
        call.status.whenFailure { err in
            self.logger.error("another callback to notify of failure: \(err.localizedDescription)")
        }
    }
}

extension UnaryCall {
    @discardableResult
    func waitForOutput() -> (ResponsePayload?, Error?) {
        os_log(.debug, "[grpc] \(path) waiting")

        do {
            let st = try status.wait()
            let msg = "\(st.code) - \(st.message)"
            os_log(.debug, "\(msg)")
        } catch {
            os_log(.error, "\(error.localizedDescription)")
        }
        do {
            let res = try response.wait()
            os_log(.debug, "[grpc] \(path) returned")
            return (res, nil)
        } catch {
            os_log(.error, "[grpc] \(path) failed with: \(error.localizedDescription)")
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
