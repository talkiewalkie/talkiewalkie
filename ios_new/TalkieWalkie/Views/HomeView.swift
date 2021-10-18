//
//  HomeView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import CoreData
import FirebaseAuth
import FirebaseMessaging
import OSLog
import SwiftUI

class MessageViewModel: ObservableObject {
    @Published var message: ChatMessage?
    @Published var showDetailView: Bool = false
}

struct HomeView: View {
    @EnvironmentObject var authed: AuthState

    @AppStorage("name") var name: String = ""
    @AppStorage("showOnboarding") var showOnboarding: Bool = false

    var body: some View {
        if showOnboarding {
            OnboardingView(onboardingDone: {
                self.showOnboarding = false

                DispatchQueue.global().async {
                    if case .Connected(let api, _) = authed.state {
                        _ = api.onboardingComplete(displayName: name, locales: ["fr"])
                    }
                }
            })
        } else if authed.firebaseUser == nil {
            ProgressView().onAppear {
                sleep(1)
                self.showOnboarding = true
            }
        } else {
            InboxView(model: InboxViewModel(authed))
        }
    }
}
