//
//  ConversationModelView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import Foundation

class ConversationModelView: ObservableObject {
    let api: Api

    init(api: Api) {
        self.api = api
    }

    @Published private(set) var loading = true
    @Published private(set) var friends: Api.Friends?
    @Published private(set) var groups: [Api.GroupsOutputGroup] = []

    func loadFriends() {
        api.friends { friends, _ in
            self.loading = false
            self.friends = friends
        }
    }

    func message(text: String, handles: [String]) {
        api.message(text, handles) { _, _ in }
    }

    func loadGroups() {
        loading = true
        api.groups { groups, _ in
            self.loading = false
            self.groups = groups?.groups ?? []
        }
    }

    // MARK: - Connection

    @Published var webSocketTask: MyWsDelegate?
    @Published var newMessage: Api.GroupWsMessage?

    private func onReceive(result: Result<URLSessionWebSocketTask.Message, Error>) {
        switch result {
        case .success(let m):
            switch m {
            case .string(let content):
                DispatchQueue.main.async {
                    self.newMessage = try! JSONDecoder().decode(Api.GroupWsMessage.self, from: content.data(using: .utf8)!)
                }
            default:
                print("\(Date()) ws connection only handles text messages, message ignored")
            }
        case .failure(let err):
            let nsErr = err as NSError
            if nsErr.domain == NSPOSIXErrorDomain, nsErr.code == 57 {
                print("\(Date()) unwanted disconnection from ws '\(webSocketTask?.url)'")
                webSocketTask?.reinit()
            } else {
                print("\(Date()) ws reception error: \(nsErr)")
                webSocketTask?.connected = false
            }
        }
    }

    func connect() {
        webSocketTask = api.groupWs(onReceive: onReceive)
    }

    func disconnect() {
        webSocketTask?.disconnect()
    }

    deinit { // 9
        disconnect()
    }
}
