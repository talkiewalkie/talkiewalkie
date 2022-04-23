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
        DispatchQueue.global(qos: .background).async {
            if case let .Connected(api, _) = self.authed.state, let uuid = self.conversation.uuid {
                self.loading = true
                let (remoteConv, _) = api.convByUuid(uuid)
                self.loading = false
                if let remoteConv = remoteConv {
                    self.authed.withWriteContext { ctx, _ in
                        Conversation.fromProto(remoteConv, context: ctx)
                    }
                }
            } else {
                os_log(.error, "conv without uuid!!!")
            }
        }
    }
}
