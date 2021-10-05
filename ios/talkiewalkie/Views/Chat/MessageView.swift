//
//  MessageView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 28/09/2021.
//

import SwiftUI

struct MessageView: View {
    let message: Api.GroupOutputMessage
    @EnvironmentObject var authState: AuthenticatedState

    var isSentByMe: Bool {
        return message.authorHandle == authState.me.handle
    }

    var body: some View {
        HStack {
            if isSentByMe { Spacer() }
            Text(message.text)
                .padding(10)
                .foregroundColor(isSentByMe ? Color.white : Color.black)
                .background(isSentByMe ? Color.blue : Color(UIColor(red: 240/255, green: 240/255, blue: 240/255, alpha: 1.0)))
                .cornerRadius(10)
        }
        .padding(10)
        .frame(width: .infinity)
    }
}

struct MessageView_Previews: PreviewProvider {
    static var previews: some View {
        let msg = Api.GroupOutputMessage(text: "hello", createdAt: ISO8601DateFormatter().string(from: Date()), authorHandle: "toto")
        MessageView(message: msg)
    }
}