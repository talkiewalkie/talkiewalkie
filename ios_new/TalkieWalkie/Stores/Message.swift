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

    static func getByLocalUuid(_ uuid: UUID, context: NSManagedObjectContext) -> Message? {
        let request = Message.fetchRequest()
        request.predicate = NSPredicate(format: "localUuid_ = %@", uuid as NSUUID)

        let result = try! context.fetch(request)
        if result.count != 1 {
            fatalError("this should not happen: \(result.count) messages found for local uuid:[\(uuid)]")
        }

        return result.first
    }

    @discardableResult
    static func fromEventProto(_ event: App_Event, context: NSManagedObjectContext, block _: (_ event: App_Event) -> Void = { _ in }) -> Message {
        let conv: Conversation
        let message: Message
        switch event.content {
        case let .some(.receivedNewMessage(rnm)):
            if event.localUuid != "", let localMessage = Message.getByLocalUuid(event.localUuid.uuidOrThrow(), context: context) {
                message = localMessage
                message.status_ = 1
            } else {
                if rnm.hasConversation {
                    os_log(.debug, "message from new conv!: \(rnm.conversation.uuid)")
                    conv = Conversation.fromProto(rnm.conversation, context: context)
                } else {
                    conv = Conversation.getByUuidOrCreate(rnm.message.convUuid.uuidOrThrow(), context: context)
                }

                message = Message.fromProto(rnm.message, context: context)
                conv.addToMessages_(message)
            }

        case let .some(.sentNewMessage(snm)):
            switch snm.conversation {
            case let .some(.convUuid(convUuid)):
                conv = Conversation.getByUuidOrCreate(convUuid.uuidOrThrow(), context: context)

            case let .some(.newConversation(convInput)):
                conv = Conversation(context: context)
                conv.title = convInput.title
                let users: [UserConversation] = convInput.userUuids.map { uuid in
                    let u = User.getByUuidOrCreate(uuid.uuidOrThrow(), context: context)
                    let uc = UserConversation(context: context)
                    uc.user = u
                    uc.conversation = conv
                    uc.readUntil = Date()

                    return uc
                }
                conv.addToUsers_(NSSet(array: users))

            default:
                fatalError()
            }

            let content: MessageContent
            message = Message(context: context)
            message.localUuid_ = event.localUuid.uuidOrThrow()
            switch snm.message.content {
            case .none:
                fatalError()
            case let .some(.textMessage(tm)):
                let ttm = TextMessage(context: context)
                ttm.text = tm.content

                content = ttm

            case let .some(.voiceMessage(vm)):
                let ctt = VoiceMessage(context: context)
                ctt.processedAudio = vm.rawContent
                ctt.rawAudio = vm.rawContent
                ctt.siriTranscript = try! vm.siriTranscript.serializedData()

                content = ctt
            }
            message.content = content

        default:
            // TODO: find a better way and fail at compile time.
            fatalError()
        }

        return message
    }
}
