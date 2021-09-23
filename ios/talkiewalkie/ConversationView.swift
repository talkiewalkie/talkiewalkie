//
//  ConversationView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 20/09/2021.
//

import SwiftUI

struct ConversationView: View {
    @State var message: String = ""
    @State var showSheet: Bool = false
    @State var pickedHandles: [String] = []
    @StateObject var model: ConversationModelView
    @EnvironmentObject var auth: UserViewModel

    var body: some View {
        VStack(spacing: 20) {
            if self.model.webSocketTask?.connected ?? false {
                Text("Connecting...").padding()
            }
            if let msg = self.model.newMessage {
                HStack {
                    Text(msg.message).padding()
                }
            }
            if let groups = model.groups {
                ForEach(groups, id: \.uuid) { g in
                    NavigationLink(destination: ChatView(model: ChatViewModel(api: model.api, uuid: g.uuid)).environmentObject(auth)) {
                        HStack {
                            Text(g.display)
                            Spacer()
                            Text("(\(g.handles.count))")
                        }.padding()
                    }
                }
            }
            Spacer()
            HStack {
                TextField("message", text: $message).padding()
                Button("send") {
                    showSheet.toggle()
                }
            }
            .padding(10)
            .background(Color(red: 0, green: 0, blue: 0, opacity: 0.1))
            .foregroundColor(.white)
            .sheet(isPresented: $showSheet) {
                ScrollView {
                    VStack(alignment: .leading) {
                        if let friends = self.model.friends, !self.model.loading {
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
                            Spacer()
                            Button("Send") {
                                self.model.message(text: message, handles: pickedHandles)
                            }
                        } else {
                            ProgressView("Loading friends")
                        }
                    }
                }
            }
        }.onAppear {
            model.connect()
            model.loadFriends()
            model.loadGroups()
        }
        .onDisappear {
            model.disconnect()
        }
    }
}

struct ConversationView_Previews: PreviewProvider {
    static var previews: some View {
        let model = ConversationModelView(api: Api(token: "XX"))
        ConversationView(model: model)
    }
}
