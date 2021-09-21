//
//  ConversationModelView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import Foundation

class ConversationModelView: ObservableObject {
    let api: Api

    init(api: Api) {
        self.api = api
    }

    @Published private(set) var loading = true
    @Published private(set) var friends: Api.Friends?
    @Published private(set) var groups: [Api.GroupsOutputGroup] = []

    func loadFriends() {
        self.api.friends { friends, _ in
            self.loading = false
            self.friends = friends
        }
    }
    
    func message(text: String, handles: Array<String>) {
        self.api.message(text,  handles) { data, _ in }
    }
    
    func loadGroups() {
        loading = true
        self.api.groups { groups, error in
            self.loading = false
            self.groups = groups?.groups ?? []
        }
    }
}
