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

    let authed: AuthState
    let conversation: Conversation
    init(_ authed: AuthState, conversation: Conversation) {
        self.authed = authed
        self.conversation = conversation
    }

    func loadMessages() {
        DispatchQueue.global().async {
            if case .Connected(let api, _) = self.authed.state, let uuid = self.conversation.uuid {
                self.loading = true
                let (remoteConv, _) = api.convByUuid(uuid)
                self.loading = false
                if let remoteConv = remoteConv {
                    self.authed.backgroundMoc.perform {
                        Conversation.dumpFromRemote([remoteConv], context: self.authed.backgroundMoc)
                    }
                }
            } else {
                os_log(.error, "conv without uuid!!!")
            }
        }
    }
}
