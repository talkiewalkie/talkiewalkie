//
//  ChatViewModel.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import Foundation

class ChatViewModel: ObservableObject {
    let authed: AuthenticatedState
    let uuid: String
    init(authed: AuthenticatedState, uuid: String) {
        self.authed = authed
        self.uuid = uuid
    }

    @Published var loading = false
    @Published var conversation: Api.ConversationOutput?
    @Published var messages: [Api.ConversationOutputMessage] = []

    func loadMessages(page: Int) {
        authed.api.conversation(uuid, offset: page) { g, _ in
            self.conversation = g
            self.messages.append(contentsOf: g?.messages ?? [])
        }
    }

    func message(text: String) {
        authed.api.message(text, conversationUuid: uuid) { _, _ in }
    }
}
