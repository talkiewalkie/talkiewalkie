//
//  DiscussionListView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import CoreData
import SwiftUI
import OSLog

class InboxViewModel: ObservableObject {
    private let authed: AuthenticatedState

    @Published var loading = true

    init(_ authed: AuthenticatedState) {
        self.authed = authed
        self.authed.gApi.subscribeIncomingMessages { msg in
            let savedMsg = Message.upsert(msg, context: authed.context)
            let conversation = Conversation.getByUuidOrCreate(msg.convUuid.uuidOrThrow(), context: authed.context)

            conversation.addToMessages(savedMsg)

            savedMsg.objectWillChange.send()
            conversation.objectWillChange.send()

            authed.context.saveOrLogError()
        }
    }

    func syncConversations() {
        loading = true
        let (convs, _) = authed.gApi.listConvs()
        loading = false
        os_log(.debug, "loading = false")
        Conversation.dumpFromRemote(convs, context: authed.context)
    }
}

struct DiscussionListView: View {
    var namespace: Namespace.ID
    @ObservedObject var model: InboxViewModel

    @FetchRequest(
        entity: Conversation.entity(),
        sortDescriptors: []
    ) var conversations: FetchedResults<Conversation>

    var body: some View {
        NavigationView {
            VStack {
                if model.loading { ProgressView("syncing...") }

                Text("\(conversations.count) conv loaded")
                List(conversations) { conv in
                    Text(conv.title ?? "conv without title")
                }

                List(dummyDiscussions) { discussion in
                    NavigationLink(
                        destination: DiscussionView(discussion: discussion, namespace: namespace)
                    ) {
                        DiscussionListItemView(discussion: discussion)
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

struct DiscussionAvatar: View {
    var discussion: DiscussionModel

    var body: some View {
        Group {
            if let image = discussion.image {
                image
                    .resizable()
                    .aspectRatio(1, contentMode: .fit)
            } else {
                let initialLetter = discussion.name.prefix(1)
                let color = generateColorFor(text: discussion.id.uuidString)

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
        }
        .clipShape(Circle())
    }
}

struct DiscussionListItemView: View {
    var discussion: DiscussionModel

    var body: some View {
        HStack(alignment: .top) {
            DiscussionAvatar(discussion: discussion)

            HStack {
                VStack {
                    Text(discussion.name)
                        .fontWeight(.medium)
                }

                Spacer()

                Text("\(discussion.date, formatter: Self.dateFormat)")
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
