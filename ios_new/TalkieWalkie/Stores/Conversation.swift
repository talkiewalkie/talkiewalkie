//
//  Conversation.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 04/10/2021.
//

import CoreData
import Foundation

extension Conversation {
    static func upsert(_ conv: App_Conversation, context: NSManagedObjectContext) -> Conversation {
        let localC = Conversation.getByUuidOrCreate(conv.uuid.uuidOrThrow(), context: context)
        localC.uuid = conv.uuid.uuidOrThrow()
        localC.title = conv.title

        let messages = conv.messages.map { Message.upsert($0, context: context) }
        localC.addToMessages(NSOrderedSet(array: messages))

        // TODO: server should store and retrieve these
        // localC.lastMemberReadUntil = Date().addingTimeInterval()
        localC.createdAt = Date()

        return localC
    }

    static func dumpFromRemote(_ convs: [App_Conversation], context: NSManagedObjectContext) {
        let localConvs: [Conversation] = convs.map { Conversation.upsert($0, context: context) }
        localConvs.forEach { $0.objectWillChange.send() }
        context.saveOrLogError()
    }
}
