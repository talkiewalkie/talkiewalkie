//
//  ConversationView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import SwiftUI

struct InboxView: View {
    @State var showQuickSendSheet: Bool = false
    @State var pickedHandles: [String] = []

    @StateObject var model: InboxViewModel
    @EnvironmentObject var auth: AuthenticatedState

    @FetchRequest(
        entity: Conversation.entity(),
        sortDescriptors: [NSSortDescriptor(key: "lastActivityAt", ascending: false)]
    ) var conversations: FetchedResults<Conversation>

    @Environment(\.colorScheme) private var cs: ColorScheme

    var body: some View {
        ZStack(alignment: .bottomLeading) {
            VStack(spacing: 20) {
                if self.model.loading {
                    ProgressView()
                }

                Button("Refresh") { model.syncConversations() }
                    .padding()

                ScrollView {
                    ForEach(conversations, id: \.uuid?.uuidString) { conv in
                        NavigationLink(destination: ChatView(conversation: conv, authed: auth).environmentObject(auth)) {
                            HStack(spacing: 10) {
                                Circle().frame(width: 30, height: 30)
                                Text(conv.display ?? "no title")
                                    .lineLimit(1)
                                    .foregroundColor(.black)
                                Spacer()
                                Text("(\(conv.users?.count ?? 0))")
                                    .foregroundColor(.gray)
                            }
                            .padding(.horizontal, 10)
                        }
                        Divider().frame(width: UIScreen.main.bounds.width - 40, height: 2, alignment: .center).foregroundColor(.black)
                    }
                }

                Spacer()
                HStack(alignment: .center) {
                    Button(action: { showQuickSendSheet.toggle() }) {
                        Image(systemName: "mic")
                            .padding(30)
                    }
                    .foregroundColor(.white)
                    .background(Color.black)
                    .clipShape(Circle())
                    .shadow(radius: 10.0)
                }
            }
        }
        .onAppear {
//            model.connect()
//            model.loadFriends()
            model.syncConversations()
        }
        .sheet(isPresented: $showQuickSendSheet) {
            if let friends = self.model.friends, !self.model.loading {
                VStack(alignment: .leading) {
                    ScrollView {
                        ForEach(friends.friends, id: \.uuid) { g in
                            HStack {
                                let amIPicked: Binding<Bool> = Binding(
                                    get: { self.pickedHandles.contains(g.display) },
                                    set: { if $0 { self.pickedHandles.append(g.display) } else { self.pickedHandles = self.pickedHandles.filter { e in e != g.display } } }
                                )
                                Toggle(isOn: amIPicked) {
                                    Text(g.display)
                                }.padding()
                            }
                        }
                        Divider()
                        ForEach(friends.randoms, id: \.self) { g in
                            HStack {
                                let amIPicked: Binding<Bool> = Binding(
                                    get: { self.pickedHandles.contains(g) },
                                    set: { if $0 { self.pickedHandles.append(g) } else { self.pickedHandles = self.pickedHandles.filter { e in e != g } } }
                                )
                                Toggle(isOn: amIPicked) {
                                    Text(g)
                                }.padding()
                            }
                        }
                    }
                    .padding(.top, 20)

                    Spacer()
                    HStack {
                        Spacer()
                        Button("Send") {
                            model.message(text: "fake message", handles: pickedHandles)
                        }
                        .padding()
                        .foregroundColor(.white)
                        .background(Color.blue)
                        .cornerRadius(10)
                        Spacer()
                    }
                    .padding()
                }
            } else {
                VStack(alignment: .center) {
                    Spacer()
                    ProgressView("Loading friends")
                    Spacer()
                }
            }
        }
    }
}

// struct ConversationView_Previews: PreviewProvider {
//    static var previews: some View {
//        let model = InboxViewModel(api: Api(token: "XX"))
//        InboxView(model: model)
//    }
// }
