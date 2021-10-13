//
//  Api.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import Combine
import Foundation
import GRPC
import NIO
import OSLog

struct Api {
    private let session = URLSession.shared
    private let logger = Logger(subsystem: Bundle.main.bundleIdentifier!, category: "Network")

    #if DEBUG
        static let url = "https://theo.dev.talkiewalkie.app"
    #else
        static let url = "https://api.talkiewalkie.app"
    #endif

    #if DEBUG
        static let wsUrl = "wss://theo.dev.talkiewalkie.app/ws"
    #else
        static let wsUrl = "wss://api.talkiewalkie.app/ws"
    #endif

    var token: String

    private func get<T>(_ url: URL, completion: @escaping (T?, Error?) -> Void) where T: Codable {
        var request = URLRequest(url: url)
        request.addValue("Bearer \(token)", forHTTPHeaderField: "X-TalkieWalkie-Auth")

        logger.debug("[\(url)] GET")
        session.dataTask(with: request) { data, res, httpErr in
            if let err = httpErr { logger.error("[\(url)] error in GET request: \(err.localizedDescription)"); return }
            if let httpResponse = res as? HTTPURLResponse {
                if httpResponse.statusCode > 299 {
                    let body = data != nil ? String(data: data!, encoding: .utf8)! : "empty body"
                    logger.error("GET '\(url)' -> \(httpResponse.statusCode): '\(body)'")

                    return
                }
            }
            guard let data = data else {
                logger.debug("[\(url)] GET -> no data")

                return
            }

            DispatchQueue.main.async {
                do {
                    let output = try JSONDecoder().decode(T.self, from: data)

                    return completion(output, nil)
                } catch let err {
                    logger.error("[\(url)] GET - Could not decode JSON: \(err.localizedDescription)")

                    return completion(nil, err)
                }
            }
        }.resume()
    }

    private func post<I, T>(_ url: URL, payload: I, completion: @escaping (T?, Error?) -> Void) where T: Codable, I: Codable {
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.addValue("Bearer \(token)", forHTTPHeaderField: "X-TalkieWalkie-Auth")
        let jsonData = try? JSONEncoder().encode(payload)
        request.httpBody = jsonData

        logger.debug("[\(url)] POST")
        session.dataTask(with: request) { data, res, httpErr in
            if let err = httpErr { logger.error("[\(url)] error in POST request: \(err.localizedDescription)"); return }
            if let httpResponse = res as? HTTPURLResponse {
                if httpResponse.statusCode > 299 {
                    let body = data != nil ? String(data: data!, encoding: .utf8)! : "empty body"
                    logger.error("[\(url)] POST -> \(httpResponse.statusCode): '\(body)'")

                    return
                }
            }
            guard let data = data else {
                logger.debug("[\(url)] POST -> no data")

                return
            }

            DispatchQueue.main.async {
                do {
                    let output = try JSONDecoder().decode(T.self, from: data)

                    return completion(output, nil)
                } catch let err {
                    logger.error("[\(url)] POST - Could not decode JSON: \(err.localizedDescription)")

                    return completion(nil, err)
                }
            }
        }.resume()
    }

    func ws(path: String) -> WebSocketTaskConnection {
        var request = URLRequest(url: URL(string: "\(Api.wsUrl)/\(path)")!)
        request.httpMethod = "GET"
        request.addValue("Bearer \(token)", forHTTPHeaderField: "X-TalkieWalkie-Auth")

        return WebSocketTaskConnection(request: request)
    }

    // MARK: - CONVERSATION WS

    struct ConversationWsMessage: Codable {
        let message: String
        let authorHandle: String
        let conversationUuid: String
        let uuid: String
    }

    // MARK: - FRIENDS

    struct ConversationInfo: Codable {
        let uuid: String
        let display: String
        let handles: [String]
    }

    struct Friends: Codable {
        let friends: [ConversationInfo]
        let randoms: [String]
    }

    func friends(completion: @escaping (Api.Friends?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/me/friends")!
        get(url.url!, completion: completion)
    }

    // MARK: - MESSAGE

    struct MessageInput: Codable {
        let text: String
        let handles: [String]
    }

    struct MessageInput2: Codable {
        let text: String
        let conversationUuid: String
    }

    struct MessageOutput: Codable {}

    func message(_ text: String, _ recipients: [String], completion: @escaping (Api.MessageOutput?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/message")!
        post(url.url!, payload: MessageInput(text: text, handles: recipients), completion: completion)
    }

    func message(_ text: String, conversationUuid: String, completion: @escaping (Api.MessageOutput?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/message")!
        post(url.url!, payload: MessageInput2(text: text, conversationUuid: conversationUuid), completion: completion)
    }

    // MARK: - CONVERSATIONS

    struct ConversationsOutput: Codable {
        let conversations: [ConversationsOutputConversation]
    }

    struct ConversationsOutputConversation: Codable {
        let uuid: UUID
        let display: String
        let handles: [String]
    }

    func conversations(completion: @escaping (Api.ConversationsOutput?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/me/conversations")!
        get(url.url!, completion: completion)
    }

    // MARK: - CHAT

    struct ConversationOutput: Codable {
        let uuid: String
        let display: String
        let handles: [String]
        let messages: [ConversationOutputMessage]
    }

    struct ConversationOutputMessage: Codable {
        let text: String
        let createdAt: String
        let authorHandle: String
    }

    func conversation(_ uuid: String, offset: Int, completion: @escaping (Api.ConversationOutput?, Error?) -> Void) {
        var url = URLComponents(string: "\(Api.url)/conversation/\(uuid)")!
        url.queryItems = [URLQueryItem(name: "offset", value: String(offset))]

        get(url.url!, completion: completion)
    }

    // MARK: - ME

    struct MeOutput: Codable {
        let uuid: UUID
        let handle: String
    }

    func me(completion: @escaping (Api.MeOutput?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/me")!

        get(url.url!, completion: completion)
    }
}

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

        if let code = st?.code, code != .ok {
            os_log("[grpc] \(path) returned with status \(code.description), \((st?.message).repr) error: \((error?.localizedDescription).repr)")
            return error
        }

        os_log("[grpc] \(path) completed")
        return nil
    }
}
