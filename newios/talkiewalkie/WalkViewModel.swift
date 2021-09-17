//
//  WalkViewModel.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 13/09/2021.
//

import Foundation

class WalkViewModel: ObservableObject {
    enum Input {
        case FeedWalk(Api.WalksItem)
        case Uuid(String)
    }

    let input: WalkViewModel.Input

    init(input: WalkViewModel.Input) {
        self.input = input

        // TODO: would be nice to have an instant page show when coming from the feed
        switch input {
        case let .FeedWalk(fw):
            walk = Api.Walk(title: fw.title, description: fw.description, uuid: fw.uuid, coverUrl: fw.coverUrl, audioUrl: "fakeurl", author: fw.author)
        default:
            break
        }
    }

    @Published private(set) var walk: Api.Walk?
    @Published private(set) var loading = true

    // MARK: - Intent(s)

    func getWalk() {
        var uuid: String

        switch input {
        case let .FeedWalk(fw):
            uuid = fw.uuid
        case let .Uuid(uuid_):
            uuid = uuid_
        }

        loading = true

        Api.walk(uuid) { val, err in
            self.loading = false
            print("finished getting walk")
            if let w = val {
                print("updating walk with \(w.audioUrl)")
                self.walk = w
            } else if let error = err {
                print(error)
            }
        }
    }

    func setWalk(_ walk: Api.Walk) {
        self.walk = walk
    }
}
