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
        // TODO: grpc zeroes out msg.author when author is nil in db, which is a valid state (user deleted their accounts)
        //       this need to be handled in the next line.
        localM.author = User.upsert(msg.author, context: context)
        localM.createdAt = msg.createdAt.date
        switch msg.content {
        case let .textMessage(txt):
            let tm = TextMessage(context: context)
            tm.text = txt.content

            localM.content = tm
        case let .voiceMessage(voice):
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
