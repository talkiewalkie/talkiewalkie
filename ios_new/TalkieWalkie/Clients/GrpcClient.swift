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

class AuthedGrpcApi {
    #if DEBUG
        // static let url = URL(string: "https://theo.dev.talkiewalkie.app:443")!
        // static let url = URL(string: "http://localhost:8080")!
        static let url = URL(string: "https://7081-2a01-e34-ec46-8190-bd43-eb20-d561-c195.ngrok.io:443")!
    #else
        static let url = URL(string: "https://api.talkiewalkie.app:443")!
    #endif

    private let logger = Logger()
    private let empty = App_Empty()

    private let group: EventLoopGroup
    private let channel: ClientConnection
    private let token: String
    private let userClient: App_UserServiceClient
    private let convClient: App_ConversationServiceClient
    private let mssgClient: App_MessageServiceClient

    init(token: String) {
        self.token = token
        group = PlatformSupport.makeEventLoopGroup(loopCount: 1)

        channel = ClientConnection.insecure(group: group)
            .withKeepalive(ClientConnectionKeepalive(interval: TimeAmount.seconds(10), timeout: TimeAmount.seconds(5)))
            .withConnectionReestablishment(enabled: true)
            .connect(host: AuthedGrpcApi.url.host!, port: AuthedGrpcApi.url.port!)

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
        os_log("[grpc] waiting for \(path)")

        let res: ResponsePayload?
        do {
            res = try response.wait()
            return (res, nil)
        } catch {
            os_log("[grpc] unary call failed with: \(error.localizedDescription)")
            return (nil, error)
        }
    }
}

extension ServerStreamingCall {
    func waitCompletion() -> Error? {
        os_log("[grpc] waiting for \(path)")
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
            os_log("[grpc] \(path) returned with status \(code.description), \(msg) error: \(errMsg)")
            return error
        }

        os_log("[grpc] \(path) completed")
        return nil
    }
}
