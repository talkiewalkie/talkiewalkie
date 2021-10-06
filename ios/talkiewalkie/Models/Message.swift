//
//  Message.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 06/10/2021.
//

import CoreData
import Foundation

extension Message {
    static func upsert(_ msg: App_Message, context: NSManagedObjectContext) -> Message {
        let localM = Message.getByUuidOrCreate(msg.uuid.uuidOrThrow(), context: context)
        localM.uuid = msg.uuid.uuidOrThrow()
        localM.author = User.upsert(msg.author, context: context)
        localM.createdAt = msg.createdAt.date
        switch msg.content {
        case .textMessage(let txt):
            localM.text = txt.content
        case .voiceMessage(_):
            localM.text = "{{ THIS IS A VOICE MESSAGE }}"
        default:
            fatalError()
        }

        return localM
    }
}
