//
//  ConversationViewModel.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 16/10/2021.
//

import Foundation
import OSLog

class ConversationViewModel: ObservableObject {
    @Published var loading = false

    let authed: AuthenticatedState
    let conversation: Conversation
    init(_ authed: AuthenticatedState, conversation: Conversation) {
        self.authed = authed
        self.conversation = conversation
    }

    func loadMessages() {
        DispatchQueue.main.async {
            if let uuid = self.conversation.uuid {
                self.loading = true
                let (remoteConv, _) = self.authed.gApi.convByUuid(uuid)
                self.loading = false
                if let remoteConv = remoteConv { Conversation.dumpFromRemote([remoteConv], context: self.authed.context) }
            } else {
                os_log(.error, "conv without uuid!!!")
            }
        }
    }
}
