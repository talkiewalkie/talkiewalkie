//
//  talkiewalkieApp.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import Firebase
import SwiftUI

@main
struct talkiewalkieApp: App {
    init() {
        FirebaseApp.configure()
    }

    var body: some Scene {
        WindowGroup {
            let auth = AuthViewModel()
            let home = FeedViewModel()

            AuthView<FeedView>(vm: auth) {
                FeedView(model: home)
            }
        }
    }
}
