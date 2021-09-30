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
    @Environment(\.colorScheme) private var cs: ColorScheme

    var body: some View {
        let _ = print("connect status: \(model.webSocketTask?.connected)")
        ZStack(alignment: .bottomLeading) {
            VStack(spacing: 20) {
                if let ws = self.model.webSocketTask {
                    if ws.connected { Text("Connected").padding() }
                    if !ws.connected { Text("Disconnected").padding() }
                }

                if let groups = model.groups {
                    ScrollView {
                        ForEach(groups, id: \.uuid) { g in
                            NavigationLink(destination: ChatView(model: ChatViewModel(api: model.api, uuid: g.uuid)).environmentObject(auth)) {
                                HStack(spacing: 10) {
                                    Circle().frame(width: 30, height: 30)
                                    Text(g.display)
                                        .lineLimit(1)
                                        .foregroundColor(.black)
                                    Spacer()
                                    Text("(\(g.handles.count))")
                                        .foregroundColor(.gray)
                                }
                                .padding(.horizontal, 10)
                            }
                            Divider().frame(width: UIScreen.main.bounds.width - 40, height: 2, alignment: .center).foregroundColor(.black)
                        }
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
            model.connect()
            model.loadFriends()
            model.loadGroups()
        }
        .onDisappear {
            model.disconnect()
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

struct ConversationView_Previews: PreviewProvider {
    static var previews: some View {
        let model = InboxViewModel(api: Api(token: "XX"))
        InboxView(model: model)
    }
}
