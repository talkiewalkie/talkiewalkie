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
            os_log(.debug, "subscribing...")
            api.subscribeIncomingMessages { msg in
                os_log("hello i received a new message")
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
        self.loading = true
        DispatchQueue.global(qos: .background).async {
            if case .Connected(let api, _) = self.authed.state {
                let (convs, _) = api.listConvs()
                self.authed.persistentContainer.performBackgroundTask { context in
                    Conversation.dumpFromRemote(convs, context: context)
                }
            }
            DispatchQueue.main.async { self.loading = false }
        }
    }
}
