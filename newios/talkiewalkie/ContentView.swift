//
//  ContentView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import SwiftUI

struct WalkCard: View {
    let walk: Api.Walk
    
    var body: some View {
        VStack(alignment: .leading, spacing: 10) {
            HStack(alignment: .center, spacing: nil) {
                Text("p")
                    .foregroundColor(Color.white)
                    .frame(width: 40.0, height: 40.0)
                    .background(/*@START_MENU_TOKEN@*//*@PLACEHOLDER=View@*/Color.black/*@END_MENU_TOKEN@*/)
                    .clipShape(/*@START_MENU_TOKEN@*/Circle()/*@END_MENU_TOKEN@*/)
                Text(walk.author.handle)
            }
            .padding(.horizontal, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)
            HStack {
                RemoteImage(url: walk.coverUrl)
                    .scaledToFill()
                    .aspectRatio(contentMode: .fill)
                    .frame(width: UIScreen.main.bounds.size.width - 20, height: 200)
                    .clipped()
            }
            HStack {
                Text(walk.title)
                Spacer()
                Text(distance(walk.distanceFromPoint)).foregroundColor(Color(red: 0, green: 0, blue: 0, opacity: 0.4))
            }
            .padding(.horizontal, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)
        }
        .padding(.vertical, /*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)
        .background(/*@START_MENU_TOKEN@*//*@PLACEHOLDER=View@*/Color.white/*@END_MENU_TOKEN@*/)
        .frame(width: UIScreen.main.bounds.size.width - 20)
        .cornerRadius(5)
        .overlay(
            RoundedRectangle(cornerRadius: 5)
                .stroke(Color(red: 0, green: 0, blue: 0, opacity: 0.2), lineWidth: 1)
        )
    }
}

struct ContentView: View {
    @ObservedObject var home: HomeViewModel
    @StateObject var locManager = LocationManager()
    
    var body: some View {
        NavigationView {
            ScrollView {
                LazyVStack(alignment: .center, spacing: 40.0) {
                    ForEach(home.walks, id: \.uuid){ w in
                        WalkCard(walk: w)
                            .onAppear {
                                home.loadMoreIfNeeded(w)
                            }
                    }
                    .frame(width: UIScreen.main.bounds.size.width, alignment: /*@START_MENU_TOKEN@*/.center/*@END_MENU_TOKEN@*/)
                }.background(Color(red: 0, green: 0, blue: 0, opacity: 0.1))
                Button("fetch") {
                    home.getPage(20)
                }
            }.navigationBarTitle("TalkieWalkie")
        }.onAppear() {
            home.position = locManager.lastLocation?.coordinate
            home.getPage(0)
        }
    }
}

struct ContentView_Previews: PreviewProvider {
    
    static var previews: some View {
        let home = HomeViewModel()
        home.addWalk(Api.Walk(title: "Tour of the thing", uuid: "uuid",  coverUrl: "https://picsum.photos/200",  author: Api.WalkAuthor(uuid: "uuid1", handle: "theo"), distanceFromPoint: 300))
        home.addWalk(Api.Walk(title: "Moving to montreal", uuid: "uu2id",  coverUrl: "https://picsum.photos/200",  author: Api.WalkAuthor(uuid: "uui21", handle: "georg"), distanceFromPoint: 10000))
        
        return ContentView(home: home)
    }
}
