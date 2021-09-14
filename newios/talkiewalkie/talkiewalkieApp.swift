//
//  talkiewalkieApp.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import SwiftUI

@main
struct talkiewalkieApp: App {
    var home = FeedViewModel()

    var body: some Scene {
        WindowGroup {
            FeedView(model: home)
        }
    }
}
