//
//  MessageView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 20.10.21.
//

import SwiftUI
import CoreData

struct MessageView: View {
    var message: Message
    
    var body: some View {
        HStack {
            VStack {
                Text(message.author?.displayName ?? "")
                
                content
                
                if let date = message.createdAt {
                    Text(Self.dateFormatter.string(from: date))
                }
            }
        }
    }
    
    var content: some View {
        Group {
            switch message.content {
                case let tm as TextMessage:
                    Text(tm.text ?? "no text")
                case let vm as VoiceMessage:
                Text(String(data: vm.rawAudio ?? Data(), encoding: .utf8) ?? "no audio")
                default:
                    Text("")
            }
        }
    }
    
    static let dateFormatter: DateFormatter = {
        let dateFormatter = DateFormatter()
        dateFormatter.dateStyle = .medium
        dateFormatter.timeStyle = .none
        dateFormatter.locale = Locale.current
        return dateFormatter
    }()
}

struct BubbleView_Previews: PreviewProvider {
    static let persistentContainer: NSPersistentContainer = {
        let container = NSPersistentContainer(name: "LocalModels")
        
        container.loadPersistentStores { description, error in
            guard let error = error else { return }
            fatalError("Core Data error: '\(error.localizedDescription)'.")
        }
        
        return container
    }()
    
    static var previews: some View {
        TestView()
            .environment(\.managedObjectContext, persistentContainer.viewContext)
    }
    
    struct TestView: View {
        @Environment(\.managedObjectContext) var moc
        
        var message: Message {
            let user = App_User.with { u in
                u.uuid = UUID().uuidString
                u.displayName = "Lisa"
            }
            
            let textMessage = App_TextMessage.with { tm in
                tm.content = "Hello there"
            }
            
            let message = App_Message.with { m in
                m.uuid = UUID().uuidString
                m.createdAt = .init(date: Date())
                m.author = user
                m.content = App_Message.OneOf_Content.textMessage(textMessage)
            }
            
            return Message.upsert(message, context: moc)
        }
        
        var body: some View {
            MessageView(message: message)
        }
    }
}
