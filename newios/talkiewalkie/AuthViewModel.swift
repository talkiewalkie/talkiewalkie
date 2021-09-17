//
//  AuthViewModel.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 16/09/2021.
//

import FirebaseAuth
import Foundation

class AuthViewModel: ObservableObject {
    private var auth = Auth.auth()

    @Published var user: User?
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

extension User {
    func token() {
        // get jwt
    }
}

class UserViewModel: ObservableObject {
    var user: User

    init(user: User) {
        self.user = user
    }
}
