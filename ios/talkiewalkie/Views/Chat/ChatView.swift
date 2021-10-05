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

    var body: some View {
        VStack {
            ScrollView {
                if let msgs = model.messages {
                    VStack(alignment: .leading) { ForEach(msgs, id: \.text) { m in
                        MessageView(message: m)
                    }}
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
