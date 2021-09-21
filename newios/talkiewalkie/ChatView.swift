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
    @EnvironmentObject var auth: UserViewModel

    var body: some View {
        VStack {
            ScrollView {
                if let msgs = model.messages {
                    VStack(alignment: .leading) { ForEach(msgs, id: \.text) { m in
                        HStack {
                            if m.authorHandle == auth.user.displayName { Spacer() }
                            Text(m.text).padding().background(Color.gray).foregroundColor(.white)
                        }
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
        }.onAppear {
            model.loadMessages(page: 0)
        }
    }
}

struct ChatView_Previews: PreviewProvider {
    static var previews: some View {
        let model = ChatViewModel(api: Api(token: ""), uuid: "")
        ChatView(model: model)
    }
}
