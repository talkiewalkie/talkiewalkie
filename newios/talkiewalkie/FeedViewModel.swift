//
//  HomeViewModel.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import CoreLocation
import Foundation

class FeedViewModel: ObservableObject {
    @Published private(set) var walks: [Api.WalksItem] = []
    @Published private(set) var loading = false // true
    @Published var position: CLLocationCoordinate2D?

    // MARK: - Intent(s)

    func getPage(_ page: Int = 0) {
        loading = true
        Api.walks(offset: page, position: position) { val, _ in
            self.loading = false
            if let walks = val {
                print("fetched \(walks.count) new walks, adding to existing \(self.walks.count)")
                self.walks.append(contentsOf: walks)
            } else {
                print("no walks")
                return
            }
        }
    }

    func addWalk(_ w: Api.WalksItem) {
        walks.append(w)
    }

    func loadMoreIfNeeded(_ w: Api.WalksItem) {
        // TODO: include hasNext in api results.
        if w.uuid == walks.last?.uuid {
            getPage(walks.count)
        }
    }
}
