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
    
    var showAuthor: Bool = true
    
    var isMe: Bool { // TODO: wrong logic
        return message.author == nil
    }
    
    var body: some View {
        HStack {
            VStack(alignment: .leading) {
                if showAuthor, let displayName = message.author?.displayName {
                    Text(displayName)
                        .foregroundColor(.init(generateColorFor(text: message.author?.displayName ?? "")))
                }
                
                content
            }
            .frame(minWidth: .zero, maxWidth: DrawingConstraints.maxWidth, alignment: .leading)
            .padding(.bottom, 0)
            .padding(.trailing, 2)
            .overlay(
                HStack(spacing: 4) {
                    if let date = message.createdAt {
                        Text(Self.dateFormatter.string(from: date))
                            .font(.footnote)
                            .foregroundColor(.init(UIColor.tertiaryLabel))
                    }
                    
                    if isMe {
                        checkStatus
                    }
                   
                }, alignment: .bottomTrailing
            )
            .padding(8)
        }
        .background(isMe ? Color("LightGreen") : .white)
        .cornerRadius(10)
        .clipped()
        .shadow(color: .black.opacity(0.1), radius: 2)
        .frame(maxWidth: .infinity, alignment: isMe ? .trailing : .leading)
    }
    
    var checkStatus: some View {
        let checkmark = Image(systemName: "checkmark")
            .font(.footnote)
            .foregroundColor(.accentColor)
        
        return HStack(spacing: -9) {
            checkmark
            checkmark
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
    
    struct DrawingConstraints {
        static let maxWidth: CGFloat = 315
    }
    
    static let dateFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "HH:mm"
        return formatter
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
                tm.content = "Hello there Hello there Hello there Hello there Hello there Hello there"
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
