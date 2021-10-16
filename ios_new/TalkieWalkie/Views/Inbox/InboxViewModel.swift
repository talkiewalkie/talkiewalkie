//
//  InboxViewModel.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 16/10/2021.
//

import Foundation
import OSLog

class InboxViewModel: ObservableObject {
    private let authed: AuthenticatedState

    @Published var loading = true

    init(_ authed: AuthenticatedState) {
        self.authed = authed
        self.authed.gApi.subscribeIncomingMessages { msg in
            let savedMsg = Message.upsert(msg, context: authed.context)
            let conversation = Conversation.getByUuidOrCreate(msg.convUuid.uuidOrThrow(), context: authed.context)

            conversation.addToMessages(savedMsg)

            savedMsg.objectWillChange.send()
            conversation.objectWillChange.send()

            authed.context.saveOrLogError()
        }
    }

    func syncConversations() {
        DispatchQueue.main.async {
            self.loading = true
            let (convs, _) = self.authed.gApi.listConvs()
            self.loading = false
            os_log(.debug, "loading = false")
            Conversation.dumpFromRemote(convs, context: self.authed.context)
        }
    }
}
