//
//  HomeView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import CoreData
import FirebaseAuth
import OSLog
import SwiftUI

class MessageViewModel: ObservableObject {
    @Published var message: ChatMessage?
    @Published var showDetailView: Bool = false
}

class HomeViewModel: ObservableObject {
    private var auth = Auth.auth()
    private var coredataCtx: NSManagedObjectContext

    @AppStorage("showOnboarding") var showOnboarding: Bool = false {
        willSet { objectWillChange.send() }
    }

    @Published var user: FirebaseAuth.User?
    @Published var authed: AuthenticatedState?

    init(_ ctx: NSManagedObjectContext, config: Config) {
        coredataCtx = ctx
        user = auth.currentUser

        auth.addStateDidChangeListener { _, user in
            self.user = user

            if let user = user, self.showOnboarding {
                AuthenticatedState.build(config, fbU: user, context: self.coredataCtx) { s in
                    self.authed = s
                }
            } else {
                self.authed = nil
                self.showOnboarding = true
            }
        }
    }
}

class AuthenticatedState: ObservableObject {
    var user: FirebaseAuth.User
    var gApi: AuthedGrpcApi
    var context: NSManagedObjectContext

    var me: Me { Me.fromCache(context: context)! }

    static func build(_ config: Config, fbU: FirebaseAuth.User, context: NSManagedObjectContext, completion: @escaping (AuthenticatedState) -> Void) {
        fbU.getIDTokenResult { res, error in
            if let token = res {
                let gApi = AuthedGrpcApi(url: config.apiUrl, token: token.token)

                if let me = Me.fromCache(context: context) {
                    // TODO: even in this case we should query the server to update
                    // the local cache, just in case.
                    os_log("loaded user info from cache: \(me)")
                    completion(AuthenticatedState(user: fbU, gApi: gApi, context: context))
                } else {
                    // TODO: add a timeout to api calls.
                    DispatchQueue.main.async {
                        let (res, _) = gApi.me()
                        if let res = res {
                            let me = Me(context: context)
                            // TODO: Me.upsert??
                            me.uuid = UUID(uuidString: res.user.uuid)
                            me.displayName = res.user.displayName

                            me.firebaseUid = fbU.uid

                            me.objectWillChange.send()
                            context.saveOrLogError()

                            completion(AuthenticatedState(user: fbU, gApi: gApi, context: context))
                        }
                    }
                }
            } else if let error = error {
                os_log(.error, "could not get firebase token: \(error.localizedDescription)")
            }
        }
    }

    private init(user: FirebaseAuth.User, gApi: AuthedGrpcApi, context: NSManagedObjectContext) {
        self.user = user
        self.gApi = gApi
        self.context = context
    }
}

struct HomeView: View {
    @StateObject var messageViewModel = MessageViewModel()
    @ObservedObject var homeViewModel: HomeViewModel

    @AppStorage("name") var name: String = ""

    var body: some View {
        Group {
            if homeViewModel.showOnboarding {
                OnboardingView {
                    homeViewModel.showOnboarding = false
                    guard let authState = homeViewModel.authed else {
                        fatalError("on onboarding finish we cannot be in a state without a token")
                    }

                    _ = authState.gApi.onboardingComplete(displayName: name, locales: ["fr"])
                }
            } else if let authState = homeViewModel.authed {
                AuthedView()
                    .environmentObject(authState)
            } else {
                ProgressView()
                    .onAppear {
                        sleep(1)
                        homeViewModel.showOnboarding = true
                    }
            }
        }
        .environmentObject(messageViewModel)
    }
}

struct HomeView_Previews: PreviewProvider {
    static var previews: some View {
        HomeView(homeViewModel: HomeViewModel(NSManagedObjectContext(concurrencyType: .mainQueueConcurrencyType), config: Config.load(version: "dev")))
            .withDummyVariables()
    }
}
