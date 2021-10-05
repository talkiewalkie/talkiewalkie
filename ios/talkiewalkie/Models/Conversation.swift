//
//  Conversation.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 04/10/2021.
//

import CoreData
import Foundation

extension Conversation {
    static func dumpFromRemote(_ convs: Api.ConversationsOutput, context: NSManagedObjectContext) {
        let localConvs: [Conversation] = convs.conversations.map { conv in
            let localConv = Conversation.getByUuidOrCreate(conv.uuid, context: context)
            localConv.uuid = conv.uuid
            localConv.display = conv.display

            // TODO: server should store and retrieve these
            localConv.lastActivityAt = Date()
            localConv.createdAt = Date()

            return localConv
        }

        localConvs.forEach { $0.objectWillChange.send() }
        context.saveOrLogError()
    }

    func loadMessages(_: UUID, page _: Int = 0) {
//        self.addToMessages(<#T##value: Message##Message#>)
    }
}
