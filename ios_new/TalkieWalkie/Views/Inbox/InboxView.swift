//
//  DiscussionListView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import CoreData
import OSLog
import SwiftUI

struct InboxView: View {
    @Namespace var namespace
    @ObservedObject var model: InboxViewModel
    @EnvironmentObject var authed: AuthState

    @State var guideState = false
    @State var isRecording = false

    @FetchRequest(
        entity: Conversation.entity(),
        sortDescriptors: []
    ) var conversations: FetchedResults<Conversation>

    var body: some View {
        NavigationView {
            ZStack {
                VStack {
                    List(conversations) { conversation in
                        NavigationLink(
                            destination: ConversationView(conversation: conversation, namespace: namespace, model: ConversationViewModel(authed, conversation: conversation))
                        ) {
                            ConversationListItemView(conversation: conversation)
                        }
                    }
                    .listStyle(.plain)
                }
            }
            .navigationTitle("Chats")
            .navigationBarItems(
                leading: HeaderSettingsView(),
                trailing: Button(action: model.syncConversations) { Image(systemName: "arrow.clockwise") }
            )
            .toolbar {
                ToolbarItem(placement: .principal) {
                    // TODO: Not working
                    if case .Connected(let api, _) = authed.state {
                        switch api.stateDelegate.state {
                        case .Disconnected:
                            Text("disconnected")
                        case .Connected:
                            Text("connected")
                        case .Connecting:
                            Text("connecting...")
                        }
                    } else {
                        EmptyView()
                    }
                }
            }
            .onAppear { model.syncConversations() }
        }
    }
}

struct ConversationAvatar: View {
    @ObservedObject var conversation: Conversation
    @EnvironmentObject var authed: AuthState

    var body: some View {
        Group {
            let initialLetter = conversation.firstParticipant(thatIsNot: authed.me)?.displayName?.prefix(1) ?? "T"
            let color = generateColorFor(text: conversation.uuid?.uuidString ?? UUID().uuidString)

            ZStack {
                Color(color)
                    .brightness(-0.1)

                Text(initialLetter.count > 0 ? initialLetter : "A")
                    .font(.title2)
                    .fontWeight(.medium)
                    .foregroundColor(.white)
            }
            .aspectRatio(1, contentMode: .fit)
        }
        .clipShape(Circle())
    }
}

struct ConversationListItemView: View {
    @ObservedObject var conversation: Conversation

    var convPreview: String {
        switch conversation.lastMessage?.content {
        case nil:
            return "No messages yet!"
        case let tm as TextMessage:
            return tm.text ?? "weird"
        case _ as VoiceMessage:
            return "audio message"
        default:
            return "new message!"
        }
    }

    var body: some View {
        HStack(alignment: .top) {
            ConversationAvatar(conversation: conversation)

            HStack {
                VStack(alignment: .leading) {
                    Text(conversation.title ?? "new conv")
                        .font(.body)
                        .fontWeight(.medium)
                    Spacer()
                    Text(convPreview)
                        .font(.callout)
                        .foregroundColor(.gray)
                }

                Spacer()

                Text("\(conversation.lastActivity ?? Date(), formatter: Self.dateFormat)")
                    .foregroundColor(.secondary)
            }.padding(.vertical, 5)
        }
        .frame(height: DrawingConstraints.height)
    }

    static let dateFormat: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "HH:mm"
        return formatter
    }()

    enum DrawingConstraints {
        static let height: CGFloat = 60
    }
}

// struct DiscussionListView_Previews: PreviewProvider {
//    static var previews: some View {
//        TestView()
//    }
//
//    struct TestView: View {
//        @Namespace var namespace
//
//        let vm =AuthenticatedState.dummy()
//        var body: some View {
//            DiscussionListView(namespace: namespace, model: vm)
//        }
//    }
// }

struct DiscussionModel: Identifiable {
    let id = UUID()

    var image: Image?
    var name: String
    var date: Date
}

let dummyImages = [
    Image(uiImage: #imageLiteral(resourceName: "profile4")),
    Image(uiImage: #imageLiteral(resourceName: "profile1")),
    Image(uiImage: #imageLiteral(resourceName: "profile2")),
    Image(uiImage: #imageLiteral(resourceName: "profile3")),
]

let dummyDiscussions = [
    DiscussionModel(image: Image(uiImage: #imageLiteral(resourceName: "profile4")),
                    name: "Maxime",
                    date: Calendar.current.date(byAdding: .hour, value: 0, to: Date())!),

    DiscussionModel(image: Image(uiImage: #imageLiteral(resourceName: "profile2")),
                    name: "Nina",
                    date: Calendar.current.date(byAdding: .hour, value: -1, to: Date())!),

    DiscussionModel(image: Image(uiImage: #imageLiteral(resourceName: "profile1")),
                    name: "Nicolas",
                    date: Calendar.current.date(byAdding: .hour, value: -3, to: Date())!),

    DiscussionModel(name: "Marie",
                    date: Calendar.current.date(byAdding: .hour, value: -4, to: Date())!),

    DiscussionModel(image: Image(uiImage: #imageLiteral(resourceName: "profile3")),
                    name: "Laura",
                    date: Calendar.current.date(byAdding: .hour, value: -6, to: Date())!),

    DiscussionModel(name: "Julien",
                    date: Calendar.current.date(byAdding: .hour, value: -7, to: Date())!),
]
