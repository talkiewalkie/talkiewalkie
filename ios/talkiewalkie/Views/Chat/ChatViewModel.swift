//
//  ChatViewModel.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import Foundation

class ChatViewModel: ObservableObject {
    let api: Api
    let uuid: String
    init(api: Api, uuid: String) {
        self.api = api
        self.uuid = uuid
    }
    
    @Published var loading = false
    @Published var group: Api.GroupOutput?
    @Published var messages: [Api.GroupOutputMessage] = []
    
    func loadMessages(page: Int) {
        self.api.group(uuid, offset: page) { g, _ in
            self.group = g
            self.messages.append(contentsOf: g?.messages ?? [])
        }
    }
    
    
    func message(text: String) {
        self.api.message(text, groupUuid: uuid) { data, _ in }
    }
    
}
