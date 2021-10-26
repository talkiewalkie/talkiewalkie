//
//  Conversation.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 04/10/2021.
//

import CoreData
import Foundation
import OSLog

extension Conversation {
    @discardableResult
    static func upsert(_ conv: App_Conversation, context: NSManagedObjectContext) -> Conversation {
        let localC = Conversation.getByUuidOrCreate(conv.uuid.uuidOrThrow(), context: context)
        localC.uuid = conv.uuid.uuidOrThrow()
        localC.title = conv.title

        let messages = conv.messages.map { Message.upsert($0, context: context) }
        localC.addToMessages(NSOrderedSet(array: messages))

        let participants = conv.participants.map { User.upsert($0, context: context) }
        localC.addToUsers(NSSet(array: participants))

        // TODO: server should store and retrieve these
        // localC.lastMemberReadUntil = Date().addingTimeInterval()
        localC.createdAt = Date()

        return localC
    }

    static func dumpFromRemote(_ convs: [App_Conversation], context: NSManagedObjectContext) {
        convs.forEach { Conversation.upsert($0, context: context) }
        context.saveOrLogError()
    }

    var lastActivity: Date? {
        (self.messages?.array as? [Message] ?? []).last?.createdAt
    }

    func firstParticipant(thatIsNot user: User?) -> User? {
        let allUsers: Set<User> = (self.users as? Set<User>) ?? Set()
        let others: [User] = allUsers.filter { $0.uuid != user?.uuid }
        return others.sorted(by: { $0.uuid?.uuidString ?? "a" < $1.uuid?.uuidString ?? "b" }).first
    }
}
