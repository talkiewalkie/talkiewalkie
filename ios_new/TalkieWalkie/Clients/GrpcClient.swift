//
//  GrpcClient.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 12/10/2021.
//

import Combine
import Foundation
import GRPC
import NIO
import OSLog
import CoreData

private class GrpcConnectivityState: ConnectivityStateDelegate {
    private let logger = Logger.withLabel("grpc-status")
    func connectivityStateDidChange(from oldState: ConnectivityState, to newState: ConnectivityState) {
        if oldState != .ready, newState == .ready { logger.debug("got connected") }
        else if oldState == .ready, newState != .ready { logger.debug("got disconnected") }
        if newState == .connecting { logger.debug("connecting") }
        if newState == .transientFailure { logger.debug("error (transient failure)") }
    }
}

class AuthedGrpcApi {
    private let url: URL

    private let logger = Logger.withLabel("grpc-client")
    private let stateDelegate = GrpcConnectivityState()
    private let empty = App_Empty()

    private let group: EventLoopGroup
    private let channel: ClientConnection
    private let token: String
    private let userClient: App_UserServiceClient
    private let convClient: App_ConversationServiceClient
    private let mssgClient: App_MessageServiceClient
    private let context: NSManagedObjectContext

    init(url: URL, token: String, context: NSManagedObjectContext) {
        self.url = url
        self.token = token
        self.context = context

        group = PlatformSupport.makeEventLoopGroup(loopCount: 1)

        #if DEBUG
        let channelBuilder = ClientConnection.insecure(group: group)
        #else
        let channelBuilder = ClientConnection.usingPlatformAppropriateTLS(for: group)
        #endif
        channel = channelBuilder.withKeepalive(ClientConnectionKeepalive(interval: TimeAmount.seconds(10), timeout: TimeAmount.seconds(5)))
            .withConnectionReestablishment(enabled: true)
            .withConnectivityStateDelegate(stateDelegate, executingOn: DispatchQueue.main)
            .connect(host: url.host!, port: url.port!)

        let authedOption = CallOptions(customMetadata: ["Authorization": "Bearer \(token)"])

        userClient = App_UserServiceClient(channel: channel)
        userClient.defaultCallOptions = authedOption

        convClient = App_ConversationServiceClient(channel: channel)
        convClient.defaultCallOptions = authedOption

        mssgClient = App_MessageServiceClient(channel: channel)
        mssgClient.defaultCallOptions = authedOption
    }

    deinit {
        do {
            try channel.close().wait()
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
            twCl.users.forEach { u in
                _ = User.upsert(u, context: context)
            }
        }
        
        context.saveOrLogError()
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
        let call = mssgClient.incoming(empty, handler: completion)
        _ = call.status.recover { err in
            self.logger.error("received error from incoming messages stream: \(err.localizedDescription)")

            return .processingError
        }
    }
}

extension UnaryCall {
    @discardableResult
    func waitForOutput() -> (ResponsePayload?, Error?) {
        os_log(.debug, "[grpc] \(path) waiting")

        let res: ResponsePayload?
        do {
            res = try response.wait()
            os_log(.debug, "[grpc] \(path) returned")
            return (res, nil)
        } catch {
            os_log(.debug, "[grpc] \(path) failed with: \(error.localizedDescription)")
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
