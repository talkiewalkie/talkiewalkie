//
//  ChatView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import SwiftUI

struct ChatView: View {
    @State var message: String = ""
    @ObservedObject var model: ChatViewModel

    @FetchRequest var conversations: FetchedResults<Conversation>
    @FetchRequest var messages: FetchedResults<Message>

    init(uuid: UUID, authed: AuthenticatedState) {
        self._conversations = FetchRequest(
            entity: Conversation.entity(),
            sortDescriptors: [],
            predicate: NSPredicate(format: "uuid = %@", uuid.uuidString)
        )
        self._messages = FetchRequest(
            entity: Message.entity(),
            sortDescriptors: [NSSortDescriptor(key: "createdAt", ascending: true)],
            predicate: NSPredicate(format: "conversationUuid = %@", uuid.uuidString)
        )
        self.model = ChatViewModel(authed: authed, uuid: uuid)
    }

    var body: some View {
        VStack {
            ScrollView {
//                if let msgs = conversations.first?.messages?.allObjects as? Array<Message> {
//                    VStack(alignment: .leading) { ForEach(msgs, id: \.text) { m in
//                        MessageView(message: m)
//                    }}
//                }

                VStack(alignment: .leading) {
                    ForEach(messages, id: \.text) { m in
                        MessageView(message: m)
                    }
                }
            }
            HStack {
                TextField("Message", text: $message)
                    .padding()
                Button("send") {
                    self.model.message(text: message)
                }
                .padding()
            }.padding()
        }
        .onAppear {
            model.loadMessages(page: 0)
        }
        .navigationTitle(Text(model.conversation?.display ?? "Loading conversation..."))
    }
}

// struct ChatView_Previews: PreviewProvider {
//    static var previews: some View {
//        let model = ChatViewModel(api: Api(token: ""), uuid: "")
//        ChatView(model: model)
//    }
// }
