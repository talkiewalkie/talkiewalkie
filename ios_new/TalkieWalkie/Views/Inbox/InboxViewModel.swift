//
//  InboxViewModel.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 16/10/2021.
//

import Foundation
import OSLog

class InboxViewModel: ObservableObject {
    private let authed: AuthState

    @Published var loading = true

    init(_ authed: AuthState) {
        self.authed = authed

        if case .Connected(let api, _) = authed.state {
            api.subscribeIncomingMessages { msg in
                let savedMsg = Message.upsert(msg, context: authed.moc)
                let conversation = Conversation.getByUuidOrCreate(msg.convUuid.uuidOrThrow(), context: authed.moc)

                conversation.addToMessages(savedMsg)

                savedMsg.objectWillChange.send()
                conversation.objectWillChange.send()

                authed.moc.saveOrLogError()
            }
        }
    }

    func syncConversations() {
        DispatchQueue.global().async {
            if case .Connected(let api, _) = self.authed.state {
                self.loading = true
                let (convs, _) = api.listConvs()
                self.loading = false
                os_log(.debug, "loading = false")
                Conversation.dumpFromRemote(convs, context: self.authed.moc)
            }
        }
    }
}
