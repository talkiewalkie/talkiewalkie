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
    @Published var conversation: Api.ConversationOutput?
    @Published var messages: [Api.ConversationOutputMessage] = []

    func loadMessages(page: Int) {
        let (conv, _) = authed.gApi.convByUuid(uuid)
        if let conv = conv {
            let localConv = Conversation.getByUuidOrCreate(uuid, context: authed.context)
            localConv.uuid = conv.uuid.uuidOrThrow()
            localConv.display = conv.title
            let localMsgs: [Message] = conv.messages.map { m in
                let localMsg = Message(context: authed.context)
                localMsg.conversationUuid = uuid
                
                let author = User.getByUuidOrCreate(m.authorUuid.uuidOrThrow(), context: authed.context)
                localMsg.author = author
                localMsg.createdAt = m.createdAt.date
                
                switch m.content {
                case .textMessage(let content):
                    localMsg.text = content.content
                default:
                    os_log("error unkown message type")
                }

                localMsg.addToConversation(localConv)
                localMsg.objectWillChange.send()
                authed.context.saveOrLogError()
                
                return localMsg
            }

            localConv.addToMessages(NSSet(array: localMsgs))
            localConv.objectWillChange.send()
            authed.context.saveOrLogError()
        }
    }

    func message(text: String) {
        authed.gApi.sendMessage(text: text, convUuid: uuid)
    }
}
