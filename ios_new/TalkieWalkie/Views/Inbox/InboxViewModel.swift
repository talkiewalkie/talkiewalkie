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

                authed.withWriteContext { ctx, _ in
                    let savedMsg = Message.fromProto(msg, context: ctx)
                    let conversation = Conversation.getByUuidOrCreate(msg.convUuid.uuidOrThrow(), context: ctx)

                    conversation.addToMessages_(savedMsg)

                    DispatchQueue.main.async {
                        savedMsg.objectWillChange.send()
                        conversation.objectWillChange.send()
                    }
                }
            }
        }
    }

    func syncConversations() {
        self.loading = true
        DispatchQueue.global(qos: .background).async {
            if case .Connected(let api, _) = self.authed.state {
                let (convs, _) = api.listConvs()
                self.authed.withWriteContext { ctx, _ in
                    convs.forEach { remoteConv in Conversation.fromProto(remoteConv, context: ctx) }
                }
            }
            DispatchQueue.main.async { self.loading = false }
        }
    }
}
