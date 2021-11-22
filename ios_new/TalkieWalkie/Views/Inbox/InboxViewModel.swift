//
//  InboxViewModel.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 16/10/2021.
//

import Foundation
import OSLog

class InboxViewModel: ObservableObject {
    private let authed: AuthState

    @Published var loading = false

    init(_ authed: AuthState) {
        self.authed = authed
    }
    
    func sync() {
        if case .Connected(let api, _) = authed.state {
            let (down, _) = api.sync()
            if let down = down {
                authed.withWriteContext {ctx, _ in
                    down.events.forEach { LoadEventToCoreData($0, ctx: ctx) }
                }
            }
        }
    }
}
