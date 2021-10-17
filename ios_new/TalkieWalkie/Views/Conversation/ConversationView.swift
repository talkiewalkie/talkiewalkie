//
//  DiscussionView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import AVKit
import SwiftUI

struct ConversationView: View {
    var conversation: Conversation
    var namespace: Namespace.ID
    
    let model: ConversationViewModel

    var body: some View {
        VStack {
            if model.loading { ProgressView() } 
            Rectangle().frame(height: 0)

            ScrollViewReader { _ in
                ReversedScrollView(.vertical, showsIndicators: false) {
                    VStack {
//                        ForEach(dummyChatMessages) { message in
//                            ChatBubble(message: message, namespace: namespace)
//                        }
                        Text("\(conversation.messages?.count ?? 0) messages here")
                        ForEach(conversation.messages!.array as! [Message]) { message in
                            switch message.content! {
                            case let tmee as TextMessage:
                                Text(tmee.text ?? "no text")
                            case let vm as VoiceMessage:
                                Text(String(data: vm.rawAudio!, encoding: .utf8) ?? "no audio")
                            default:
                                Text("fatalerrrrrrror")
                            }
                        }
                    }
                    .frame(maxWidth: .infinity, alignment: .trailing)
                    .padding(.vertical)
                    .padding(.horizontal)
                }
                .background(Color(#colorLiteral(red: 0.9024571472, green: 0.9024571472, blue: 0.9024571472, alpha: 1)))

                Rectangle()
                    .foregroundColor(.clear)
                    .frame(height: 50)
            }
        }
        .navigationBarTitleDisplayMode(.inline)
        .navigationBarItems(leading: ConversationBarView(conversation: conversation),
                            trailing: EmptyView())
        .onAppear { model.loadMessages() }
    }
}

struct ChatBubble: View {
    var message: ChatMessage
    var namespace: Namespace.ID

    @EnvironmentObject var messageViewModel: MessageViewModel

    var body: some View {
        /* RoundedRectangle(cornerRadius: 12)
         .foregroundColor(.gray)
         .frame(width: 200, height: 50)
         .padding(.horizontal) */

        HStack {
            if let m = messageViewModel.message, m.id == message.id {
                Group {
                    switch message.type {
                    case let .text(content: text):
                        Audiogram1(text: text)
                            .cornerRadius(20)
                    default:
                        EmptyView()
                    }
                }

            } else {
                Group {
                    switch message.type {
                    case let .text(content: text):
                        Audiogram1(text: text)
                            .matchedGeometryEffect(id: message.id.uuidString, in: namespace, isSource: true)
                            .cornerRadius(20)
                    default:
                        EmptyView()
                    }
                }
                .onTapGesture {
                    messageViewModel.message = message

                    withAnimation {
                        messageViewModel.showDetailView = true
                    }
                }
            }
        }
    }
}

struct ConversationBarView: View {
    var conversation: Conversation

    var body: some View {
        HStack {
            ConversationAvatar(conversation: conversation)

            VStack(alignment: .leading) {
                Text(conversation.title ?? "no conv title")

                Text("Last seen at \(conversation.lastActivity ?? Date(), formatter: Self.dateFormat)")
                    .fontWeight(.regular)
                    .foregroundColor(.secondary)
            }
        }
        .frame(height: DrawingConstraints.height)
    }

    static let dateFormat: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "HH:mm"
        return formatter
    }()

    enum DrawingConstraints {
        static let height: CGFloat = 30
    }
}

// struct DiscussionView_Previews: PreviewProvider {
//    static var previews: some View {
//        TestView()
//    }
//
//    struct TestView: View {
//        @Namespace var namespace
//
//        var body: some View {
//            ConversationView(conversation: dummyDiscussions[0], namespace: namespace)
//        }
//    }
// }

struct ChatMessage: Identifiable {
    let id = UUID()

    var author: String?
    var type: ChatMessageType
    var date: Date

    enum ChatMessageType {
        case text(content: String)
        case audio(url: URL)
        case image(url: URL)
    }
}

let dummyChatMessages = [
    ChatMessage(type: .text(content: "Hello hello hello hello hello hello hello hello hello hello hello hello hello!!!"),
                date: Calendar.current.date(byAdding: .minute, value: -18, to: Date())!),

    ChatMessage(author: "",
                type: .text(content: "How are you are you are you are you are you are you are you?"),
                date: Calendar.current.date(byAdding: .minute, value: -12, to: Date())!),

    ChatMessage(author: "",
                type: .text(content: "Ok cool cool cool cool cool cool cool cool cool cool"),
                date: Calendar.current.date(byAdding: .minute, value: 0, to: Date())!),
]