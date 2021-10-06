//
//  AuthViewModel.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 16/09/2021.
//

import Foundation
import SwiftUI

import CoreData
import FirebaseAuth
import OSLog

class RootViewModel: ObservableObject {
    private var coredataCtx: NSManagedObjectContext
    private var auth = Auth.auth()

    @Published var user: FirebaseAuth.User?
    @Published var authed: AuthenticatedState?

    init(_ ctx: NSManagedObjectContext) {
        coredataCtx = ctx
        user = auth.currentUser

        auth.addStateDidChangeListener { _, user in
            self.user = user

            if let user = user {
                AuthenticatedState.build(fbU: user, context: self.coredataCtx) { s in
                    self.authed = s
                }
            }
        }
    }

    var isSignedIn: Bool { user != nil }

    func signIn(_ email: String, _ password: String) {
        auth.signIn(withEmail: email, password: password) { [self] _, error in
            if let err = error {
                print(err)
            } else {
                self.user = auth.currentUser
            }
        }
    }

    func createUser(_ email: String, _ password: String) {
        auth.createUser(withEmail: email, password: password) { [self] _, error in
            if let err = error {
                print(err)
            } else {
                self.user = auth.currentUser
            }
        }
    }

    func signOut() {
        try! auth.signOut()
        user = nil
    }
}

enum AuthError: Error {
    case badBody
}

class AuthenticatedState: ObservableObject {
    var user: FirebaseAuth.User
    var api: Api
    var gApi: AuthedGrpcApi
    var context: NSManagedObjectContext

    var me: Me { Me.fromCache(context: context)! }

    static func build(fbU: FirebaseAuth.User, context: NSManagedObjectContext, completion: @escaping (AuthenticatedState) -> Void) {
        fbU.getIDTokenResult { res, _ in
            if let token = res {
                let api = Api(token: token.token)
                let gApi = AuthedGrpcApi(token: token.token)

                if let me = Me.fromCache(context: context) {
                    // TODO: even in this case we should query the server to update
                    // the local cache, just in case.
                    os_log("loaded user info from cache: \(me)")
                    completion(AuthenticatedState(user: fbU, api: api, gApi: gApi, context: context))
                } else {
                    let (res, _) = gApi.me()
                    if let res = res {
                        let me = Me(context: context)
                        me.firebaseUid = fbU.uid
                        me.someOptions = true
                        me.handle = res.user.handle
                        me.uuid = UUID(uuidString: res.user.uuid)

                        me.objectWillChange.send()
                        context.saveOrLogError()

                        completion(AuthenticatedState(user: fbU, api: api, gApi: gApi, context: context))
                    }
                }
            }
        }
    }

    private init(user: FirebaseAuth.User, api: Api, gApi: AuthedGrpcApi, context: NSManagedObjectContext) {
        self.user = user
        self.api = api
        self.gApi = gApi
        self.context = context
    }
}
