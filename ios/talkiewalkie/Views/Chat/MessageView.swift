//
//  MessageView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 28/09/2021.
//

import OSLog
import SwiftUI

struct MessageView: View {
    let message: Message
    @EnvironmentObject var authState: AuthenticatedState

    var isSentByMe: Bool {
        return message.author?.uuid == authState.me.uuid
    }

    var body: some View {
        HStack {
            if isSentByMe { Spacer() }
            Text(message.text ?? "empty msg???")
                .padding(10)
                .foregroundColor(isSentByMe ? Color.white : Color.black)
                .background(isSentByMe ? Color.blue : Color(UIColor(red: 240 / 255, green: 240 / 255, blue: 240 / 255, alpha: 1.0)))
                .cornerRadius(10)
        }
        .padding(10)
    }
}
//
//struct MessageView_Previews: PreviewProvider {
//    static var previews: some View {
//        let msg = Api.ConversationOutputMessage(text: "hello", createdAt: ISO8601DateFormatter().string(from: Date()), authorHandle: "toto")
//        MessageView(message: msg)
//    }
//}
