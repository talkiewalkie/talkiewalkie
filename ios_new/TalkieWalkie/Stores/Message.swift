//
//  Message.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 06/10/2021.
//

import CoreData
import Foundation
import GRPC
import OSLog

extension Message {
    static func fromProto(_ msg: App_Message, context: NSManagedObjectContext, block: (_ message: Message) -> Void = { _ in }) -> Message {
        let localM = Message.getByUuidOrCreate(msg.uuid.uuidOrThrow(), context: context)
        localM.uuid = msg.uuid.uuidOrThrow()
        // TODO: grpc zeroes out msg.author when author is nil in db, which is a valid state (user deleted their accounts)
        //       this need to be handled in the next line.
        if msg.author.uuid == "" {} else { localM.author = User.fromProto(msg.author, context: context) }
        localM.createdAt = msg.createdAt.date

        switch msg.content {
        case let .textMessage(txt):
            let tm = TextMessage(context: context)
            tm.text = txt.content

            localM.content = tm
        case let .voiceMessage(voice):
            let vm = VoiceMessage(context: context)
            vm.rawAudio = voice.rawContent
            // TODO: change core data model to be transformable and have setter and getters handle things. ref article:
            // TODO: https://medium.com/@rohanbhale/hazards-of-using-mutable-types-as-transformable-attributes-in-core-data-2c95cdc27088
            do { vm.siriTranscript = try voice.siriTranscript.serializedData() }
            catch { os_log(.error, "failed to serialize transcript pb message: \(error.localizedDescription)") }

            localM.content = vm
        default:
            fatalError()
        }

        block(localM)

        return localM
    }

    static func getByLocalUuidOrThrow(_ uuid: UUID, context: NSManagedObjectContext) -> Message {
        let request = Message.fetchRequest()
        request.predicate = NSPredicate(format: "localUuid_ = %@", uuid as NSUUID)

        let result = try! context.fetch(request)
        if result.count != 1 {
            fatalError("this should not happen: \(result.count) messages found for local uuid:[\(uuid)]")
        }

        return result.first!
    }
}
