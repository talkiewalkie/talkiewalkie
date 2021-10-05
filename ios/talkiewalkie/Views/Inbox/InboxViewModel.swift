//
//  ConversationModelView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import Foundation

class InboxViewModel: ObservableObject {
    let authed: AuthenticatedState

    private var socketTask: WebSocketTaskConnection
    private let decoder = JSONDecoder()
    private let encoder = JSONEncoder()
    @Published var connected = false

    init(authed: AuthenticatedState) {
        self.authed = authed

        socketTask = authed.api.ws(path: "conversations")
        socketTask.delegate = self
        socketTask.connect()
    }

    // MARK: - INBOX

    @Published private(set) var loading = true

    func message(text: String, handles: [String]) {
        authed.api.message(text, handles) { _, _ in }
    }

    func syncConversations() {
        loading = true
        authed.api.conversations { conversations, _ in
            self.loading = false
            if let convs = conversations {
                Conversation.dumpFromRemote(convs, context: self.authed.context)
            }
        }
    }

    // MARK: - QUICK SEND

    @Published private(set) var loadingFriends: Bool = true
    @Published private(set) var friends: Api.Friends?

    func loadFriends() {
        loadingFriends = true
        authed.api.friends { friends, _ in
            self.loadingFriends = false
            self.friends = friends
        }
    }

    deinit {
        socketTask.disconnect()
    }
}

extension InboxViewModel: WebSocketConnectionDelegate {
    func onConnected(connection _: WebSocketConnection) {
        print("\(Date()) ws connected")
        DispatchQueue.main.async {
            self.connected = true
        }
    }

    func onDisconnected(connection _: WebSocketConnection, error _: Error?) {
        print("\(Date()) ws disconnected")
        DispatchQueue.main.async {
            self.connected = false
        }
    }

    func onError(connection _: WebSocketConnection, error: Error) {
        print("\(Date()) ws connection err: \(error)")
        socketTask.disconnect()
        socketTask.numReconnects += 1
        socketTask.connect()
    }

    func onMessage(connection _: WebSocketConnection, text: String) {
        do {
            let newMsg = try decoder.decode(Api.ConversationWsMessage.self, from: text.data(using: .utf8)!)
//            let convs = me.conversations?.array as! [Conversation]
//            let conv = convs.first { $0.uuid?.uuidString == newMsg.conversationUuid }
//            guard var conv = conv else {
//                // TODO: load the new conversation data, if the call errs then we err in the client side too.
//                return
//            }

//            conv.messages.append(Message(createdAt: Date(), text: newMsg.message))
            print("\(Date()) ws msg: \(text)")
        } catch { print(error) }
    }

    func onMessage(connection _: WebSocketConnection, data _: Data) {
        print("\(Date()) received byte message from ws connection, unhandled")
    }
}
