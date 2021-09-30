//
//  AuthViewModel.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 16/09/2021.
//

import FirebaseAuth
import Foundation

class RootViewModel: ObservableObject {
    private var auth = Auth.auth()
    
    @Published var user: FirebaseAuth.User?
    @Published var me: Api.MeOutput?
    @Published var api: Api?
    
    init() {
        self.user = auth.currentUser
        auth.addStateDidChangeListener { _, user in
            self.user = user
            
            if let user = user {
                user.getIDTokenResult { result, _ in
                    if let result = result {
                        self.api = Api(token: result.token)
                        self.api?.me { res, _ in
                            if let data = res { self.me = data }
                        }
                    }
                }
            }
        }
    }
    
    func authenticatedModel() -> AuthenticatedState? {
        guard let u = user, let a = api, let m = me else { return nil }
        return AuthenticatedState(user: u, me: m, api: a)
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
    var me: Api.MeOutput
    var api: Api
    
    init(user: FirebaseAuth.User, me: Api.MeOutput, api: Api) {
        self.user = user
        self.me = me
        self.api = api
    }
}
