//
//  AuthView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 16/09/2021.
//

import FirebaseAuth
import SwiftUI

struct AuthView<Label>: View where Label: View {
    @ObservedObject var vm: AuthViewModel
    let signedInView: () -> Label
    @State private var showingSheet = false
    @State private var signInFlow = false

    @State var email: String = ""
    @State var password: String = ""

    var body: some View {
        if let u = vm.user {
            NavigationView {
                signedInView()
                    .environmentObject(UserViewModel(user: u))
                    .sheet(isPresented: $showingSheet) {
                        VStack {
                            Text(u.email ?? "no email")
                            Button("sign out") {
                                vm.signOut()
                                showingSheet = false
                            }.padding()
                            Button("go back") { print("sign out") }
                        }
                    }
                    .navigationBarTitle("TalkieWalkie")
                    .navigationBarItems(leading: Text(""), trailing: Button("account") { showingSheet = true })
            }
        } else {
            VStack(alignment: .leading) {
                HStack {
                    Text(signInFlow ? "Sign In" : "Create your account!").font(.title)
                    Spacer()
                    Button(signInFlow ? "Create" : "Sign in") { signInFlow.toggle() }
                }

                if signInFlow {
                    Spacer()
                    Text("email").fontWeight(.bold)
                    TextField("email@example.com", text: $email)
                        .autocapitalization(/*@START_MENU_TOKEN@*/ .none/*@END_MENU_TOKEN@*/)
                        .disableAutocorrection(true)
                        .padding(/*@START_MENU_TOKEN@*/ .all/*@END_MENU_TOKEN@*/, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)

                    Text("password").fontWeight(/*@START_MENU_TOKEN@*/ .bold/*@END_MENU_TOKEN@*/)
                    SecureField("*****", text: $password)
                        .autocapitalization(/*@START_MENU_TOKEN@*/ .none/*@END_MENU_TOKEN@*/)
                        .disableAutocorrection(true)
                        .padding(/*@START_MENU_TOKEN@*/ .all/*@END_MENU_TOKEN@*/, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)

                    Button("Let's goooo!") {
                        vm.signIn(email, password)
                    }
                    .padding()
                    .background(Color.blue)
                    .foregroundColor(.white)

                    Spacer().frame(height: 200, alignment: /*@START_MENU_TOKEN@*/ .center/*@END_MENU_TOKEN@*/)
                } else {
                    Spacer()

                    Text("email").fontWeight(.bold)
                    TextField("email@example.com", text: $email)
                        .autocapitalization(/*@START_MENU_TOKEN@*/ .none/*@END_MENU_TOKEN@*/)
                        .disableAutocorrection(true)
                        .padding(/*@START_MENU_TOKEN@*/ .all/*@END_MENU_TOKEN@*/, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)

                    Text("password").fontWeight(/*@START_MENU_TOKEN@*/ .bold/*@END_MENU_TOKEN@*/)
                    SecureField("*****", text: $password)
                        .autocapitalization(/*@START_MENU_TOKEN@*/ .none/*@END_MENU_TOKEN@*/)
                        .disableAutocorrection(true)
                        .padding(/*@START_MENU_TOKEN@*/ .all/*@END_MENU_TOKEN@*/, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)

                    Button("create") {
                        vm.createUser(email, password)
                    }
                    .padding()
                    .background(Color.blue)
                    .foregroundColor(.white)

                    Spacer().frame(height: 200, alignment: /*@START_MENU_TOKEN@*/ .center/*@END_MENU_TOKEN@*/)
                }
            }.padding(.horizontal, 20).cornerRadius(/*@START_MENU_TOKEN@*/3.0/*@END_MENU_TOKEN@*/)
        }
    }
}

struct AuthView_Previews: PreviewProvider {
    static var previews: some View {
        let vm = AuthViewModel()

        return AuthView(vm: vm) { EmptyView() }
    }
}
