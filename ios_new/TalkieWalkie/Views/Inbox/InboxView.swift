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
    var namespace: Namespace.ID
    @ObservedObject var model: InboxViewModel
    @EnvironmentObject var authed: AuthenticatedState

    @FetchRequest(
        entity: Conversation.entity(),
        sortDescriptors: []
    ) var conversations: FetchedResults<Conversation>

    var body: some View {
        NavigationView {
            VStack {
                if model.loading { ProgressView("syncing...") }
                List(conversations) { conversation in
                    NavigationLink(
                        destination: ConversationView(conversation: conversation, namespace: namespace, model: ConversationViewModel(authed, conversation: conversation))
                    ) {
                        ConversationListItemView(conversation: conversation)
                    }
                }
                .listStyle(.plain)
            }
            .navigationTitle("Chats")
            .navigationBarItems(leading: HeaderSettingsView())
            .toolbar {
                ToolbarItem(placement: .principal) {
                    Text("TalkieWalkie")
                }
            }
            .onAppear { model.syncConversations() }
        }
    }
}

struct ConversationAvatar: View {
    var conversation: Conversation
    @EnvironmentObject var authed: AuthenticatedState

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
    var conversation: Conversation

    var body: some View {
        HStack(alignment: .top) {
            ConversationAvatar(conversation: conversation)

            HStack {
                VStack {
                    Text(conversation.title ?? "new conv")
                        .fontWeight(.medium)
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
