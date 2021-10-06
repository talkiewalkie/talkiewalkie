//
//  ChatViewModel.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import Foundation
import OSLog

class ChatViewModel: ObservableObject {
    let authed: AuthenticatedState
    let uuid: UUID
    init(authed: AuthenticatedState, uuid: UUID) {
        self.authed = authed
        self.uuid = uuid
    }

    @Published var loading = false

    func loadMessages(page: Int) {
        loading = true
        let (conv, _) = authed.gApi.convByUuid(uuid)
        loading = false
        
        if let conv = conv {
            let localConv = Conversation.upsert(conv, context: authed.context)

            localConv.objectWillChange.send()
            authed.context.saveOrLogError()
        }
    }

    func message(text: String) {
        authed.gApi.sendMessage(text: text, convUuid: uuid)
    }
}
