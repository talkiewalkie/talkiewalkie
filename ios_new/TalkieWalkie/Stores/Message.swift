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
            let tm = TextMessage(context: context)
            tm.text = txt.content
            
            localM.content = tm
        case .voiceMessage(let voice):
            let vm = VoiceMessage(context: context)
            // TODO: download from bucket now if we don't ship the audio direclty through gRPC
            vm.rawAudio = voice.url.data(using: .utf8)
            
            localM.content = vm
        default:
            fatalError()
        }

        return localM
    }
}
