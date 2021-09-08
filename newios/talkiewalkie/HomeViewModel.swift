//
//  HomeViewModel.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import Foundation
import CoreLocation

class HomeViewModel : ObservableObject {
    @Published private(set) var walks: Array<Api.Walk> = []
    @Published var position: CLLocationCoordinate2D?
        
    // MARK: - Intent(s)
    
    func getPage(_ page: Int = 0) -> Void {
    
        Api.walks(offset: page, position: position) { val, err in
            if let walks = val {
                print("fetched \(walks.count) new walks, adding to existing \(self.walks.count)")
                self.walks.append(contentsOf: walks)
            } else {
                print("no walks")
                return
            }
        }
    }
    
    func addWalk(_ w: Api.Walk) -> Void {
        walks.append(w)
    }
    
    func loadMoreIfNeeded(_ w: Api.Walk) -> Void {
        // TODO: include hasNext in api results.
        if w.uuid == walks.last?.uuid {
            getPage(walks.count)
        }
    }
}
