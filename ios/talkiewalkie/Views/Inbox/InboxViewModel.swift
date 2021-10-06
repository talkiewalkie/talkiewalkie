//
//  ConversationModelView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import Foundation
import OSLog

class InboxViewModel: ObservableObject {
    let authed: AuthenticatedState

    init(authed: AuthenticatedState) {
        self.authed = authed

        self.authed.gApi.subscribeIncomingMessages { msg in
            let savedMsg = Message.upsert(msg, context: authed.context)
            let conversation = Conversation.getByUuidOrCreate(msg.convUuid.uuidOrThrow(), context: authed.context)

            savedMsg.addToConversation(conversation)

            savedMsg.objectWillChange.send()
            conversation.objectWillChange.send()

            authed.context.saveOrLogError()
        }
    }

    // MARK: - INBOX

    @Published private(set) var loading = true

    func message(text: String, handles: [String]) {
        authed.api.message(text, handles) { _, _ in }
    }

    func syncConversations() {
        loading = true
        let (convs, _) = authed.gApi.listConvs()
        loading = false
        Conversation.dumpFromRemote(convs, context: authed.context)
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
}
