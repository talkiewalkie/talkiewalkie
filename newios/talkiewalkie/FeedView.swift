//
//  ContentView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import SwiftUI

struct WalkCard: View {
    let walk: Api.WalksItem

    var body: some View {
        NavigationLink(destination: WalkView(model: WalkViewModel(input: WalkViewModel.Input.FeedWalk(walk)))) {
            VStack(alignment: .leading, spacing: 10) {
                HStack(alignment: .center, spacing: nil) {
                    Text("p")
                        .foregroundColor(Color.white)
                        .frame(width: 40.0, height: 40.0)
                        .background(/*@START_MENU_TOKEN@*//*@PLACEHOLDER=View@*/Color.black/*@END_MENU_TOKEN@*/)
                        .clipShape(/*@START_MENU_TOKEN@*/Circle()/*@END_MENU_TOKEN@*/)
                    Text(walk.author.handle)
                    Spacer()
                    Button("...") {}
                }
                .padding(.horizontal, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)
                HStack {
                    RemoteImage(url: walk.coverUrl)
                        .scaledToFill()
                        .aspectRatio(contentMode: .fill)
                        .frame(width: UIScreen.main.bounds.size.width, height: 200)
                        .clipped()
                }
                HStack {
                    Text(walk.title)
                    Spacer()
                    Text(distance(walk.distanceFromPoint)).foregroundColor(Color(red: 0, green: 0, blue: 0, opacity: 0.4))
                }
                .padding(.horizontal, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)
                VStack {
                    Text(walk.description).foregroundColor(.gray).font(.footnote).lineLimit(2).truncationMode(.tail)
                }.padding(.horizontal, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)
            }
            .padding(.vertical, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)
            .background(/*@START_MENU_TOKEN@*//*@PLACEHOLDER=View@*/Color.white/*@END_MENU_TOKEN@*/)
            .frame(width: UIScreen.main.bounds.size.width)
        }.foregroundColor(.black)
    }
}

struct FeedView: View {
    @ObservedObject var model: FeedViewModel
    @State private var showingSheet = false
    @EnvironmentObject var auth: UserViewModel

    var body: some View {
        ScrollView {
            if model.loading {
                // TODO: Spacers not working, the loading text is not centered vertically
                Spacer()
                Text("loading").foregroundColor(.gray)
                Spacer()
            } else {
                LazyVStack(alignment: .center, spacing: 20.0) {
                    ForEach(model.walks, id: \.uuid) { w in
                        WalkCard(walk: w)
                            .padding(.top, model.walks.first?.uuid == w.uuid ? 20 : 0)
                            .onAppear {
                                model.loadMoreIfNeeded(w)
                            }
                    }
                    .frame(width: UIScreen.main.bounds.size.width, alignment: /*@START_MENU_TOKEN@*/ .center/*@END_MENU_TOKEN@*/)
                }
            }
        }
        .onAppear {
            model.getPage(0)
        }
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        let home = FeedViewModel()
        home.addWalk(Api.WalksItem(title: "Tour of the thing", description: "I did a thing that's really great yeah i was there a few times in November last year blablabla.", uuid: "uuid", coverUrl: "https://picsum.photos/200", author: Api.WalkAuthor(uuid: "uuid1", handle: "theo"), distanceFromPoint: 300))
        home.addWalk(Api.WalksItem(title: "Moving to montreal", description: "I did a thing that's really great yeah i was there a few times blablabla.", uuid: "uu2id", coverUrl: "https://picsum.photos/200", author: Api.WalkAuthor(uuid: "uui21", handle: "georg"), distanceFromPoint: 10000))

        return FeedView(model: home)
    }
}
