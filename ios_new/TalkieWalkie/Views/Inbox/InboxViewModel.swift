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

    @Published var loading = false

    init(_ authed: AuthState) {
        self.authed = authed

        if case .Connected(let api, _) = authed.state {
            os_log(.debug, "subscribing to grpc inbox stream...")
            api.subscribeIncomingMessages { msg in
                os_log("hello i received a new message")
                let savedMsg = Message.upsert(msg, context: authed.moc)
                let conversation = Conversation.getByUuidOrCreate(msg.convUuid.uuidOrThrow(), context: authed.moc)

                conversation.addToMessages_(savedMsg)

                authed.moc.saveOrLogError()

                DispatchQueue.main.async {
                    savedMsg.objectWillChange.send()
                    conversation.objectWillChange.send()
                }
            }
        }
    }

    func syncConversations() {
        self.loading = true
        self.authed.backgroundMoc.perform {
            if case .Connected(let api, _) = self.authed.state {
                let (convs, _) = api.listConvs()
                Conversation.dumpFromRemote(convs, context: self.authed.backgroundMoc)
            }
            DispatchQueue.main.async { self.loading = false }
        }
    }
}
