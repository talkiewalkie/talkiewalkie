//
//  ChatView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import SwiftUI

struct ChatView: View {
    let uuid: UUID

    @State var message: String = ""

    @ObservedObject var model: ChatViewModel
    @FetchRequest var conversations: FetchedResults<Conversation>
    @FetchRequest var messages: FetchedResults<Message>
    var conversation: Conversation? { conversations.first }

    init(conversation: Conversation, authed: AuthenticatedState) {
        self.uuid = conversation.uuid!
        self.model = ChatViewModel(authed: authed, uuid: uuid)

        self._conversations = FetchRequest(
            entity: Conversation.entity(),
            sortDescriptors: [],
            predicate: NSPredicate(format: "uuid = %@", self.uuid.uuidString)
        )
        self._messages = FetchRequest(
            entity: Message.entity(),
            sortDescriptors: [NSSortDescriptor(key: "createdAt", ascending: true)],
            // %K thingy from https://code.tutsplus.com/tutorials/core-data-and-swift-relationships-and-more-fetching--cms-25070
            predicate: NSPredicate(format: "%K == %@", "conversation.uuid", conversation)
        )
    }

    var body: some View {
        VStack {
            if model.loading {
                ProgressView()
            } else {
                ScrollView {
                    if let msgs = conversation?.messages2 {
                        VStack(alignment: .leading) {
                            ForEach(msgs, id: \.text) { m in
                                MessageView(message: m)
                            }
                        }
                    }
                    Text("\(messages.count) messages to display:")
                    Text("\(conversation?.messages2.count ?? 0) by the conversation fetch")
                    VStack(alignment: .leading) {
                        ForEach(messages, id: \.text) { m in
                            MessageView(message: m)
                        }
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
        .navigationTitle(Text(conversation?.display ?? "no title"))
    }
}

// struct ChatView_Previews: PreviewProvider {
//    static var previews: some View {
//        let model = ChatViewModel(api: Api(token: ""), uuid: "")
//        ChatView(model: model)
//    }
// }
