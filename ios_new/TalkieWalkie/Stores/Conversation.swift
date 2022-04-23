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
    static func fromProto(
        _ conv: App_Conversation,
        context: NSManagedObjectContext,
        block: (_ conv: Conversation) -> Void = { _ in }
    ) -> Conversation {
        let localC = Conversation.getByUuidOrCreate(conv.uuid.uuidOrThrow(), context: context)
        localC.uuid = conv.uuid.uuidOrThrow()
        localC.title = conv.title

        localC.lastActivity_ = conv.messages.last?.createdAt.date

        let messages = conv.messages.map { Message.fromProto($0, context: context) }
        localC.addToMessages_(NSOrderedSet(array: messages))

        let participants: [UserConversation] = conv.participants.map {
            let u = User.fromProto($0.user, context: context)
            let uc = UserConversation(context: context)
            uc.readUntil = $0.readUntil.date
            uc.user = u
            uc.conversation = localC
            return uc
        }
        localC.addToUsers_(NSSet(array: participants))

        block(localC)

        return localC
    }

    var lastMessage: Message? { messages.last }

    var lastActivity: Date? { lastMessage?.createdAt }

    var messages: [Message] {
        let msgs: [Message] = messages_?.array as? [Message] ?? []
        return msgs.sorted(by: {
            guard let tsA = $0.createdAt, let tsB = $1.createdAt else { return true }
            return tsA < tsB
        })
    }

    func seenMessages(for me: Me) -> [Message] {
        guard let myUC: UserConversation = users.first(where: { $0.user?.uuid == me.uuid }) else {
            os_log(.debug, "no uc found...")
            return messages
        }

        return messages.filter {
            guard let ts = $0.createdAt else {
                os_log(.debug, "message without creation date...")
                return true
            }
            return ts < (myUC.readUntil ?? Date())
        }
    }

    func unseenMessages(for me: Me) -> [Message] {
        guard let myUC: UserConversation = users.first(where: { $0.user?.uuid == me.uuid }) else {
            os_log(.debug, "no uc found...")
            return []
        }

        return messages.filter {
            guard let ts = $0.createdAt else {
                os_log(.debug, "message without creation date...")
                return false
            }
            return ts > (myUC.readUntil ?? Date())
        }
    }

    var users: [UserConversation] { Array(users_ as? Set<UserConversation> ?? Set()) }

    func firstParticipant(thatIsNot me: Me) -> User? {
        let others: [User] = users.filter { $0.user?.uuid != me.uuid }.map { $0.user }.compactMap { $0 }
        return others.sorted(by: { $0.uuid?.uuidString ?? "a" < $1.uuid?.uuidString ?? "b" }).first
    }
}
