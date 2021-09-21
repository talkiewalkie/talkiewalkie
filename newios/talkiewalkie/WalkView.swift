//
//  WalkView.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 13/09/2021.
//

import AudioKit
import AVFoundation
import SwiftUI

struct WalkView: View {
    @ObservedObject var model: WalkViewModel
    @StateObject var locManager = LocationManager()
    @State var player: AVPlayer?

    @ViewBuilder var imgOverlay: some View {
        if model.loading {
            Image(systemName: "loading")
        } else if player != nil {
            Image(systemName: "play")
                .opacity(0.9)
                .frame(width: 200, height: 200, alignment: .center)
        } else {
            EmptyView()
        }
    }

    var body: some View {
        if let w = model.walk {
            ScrollView {
                Button(action: {
                    if let pl = player, !model.loading {
                        print("player is not null!, volume at \(pl.volume)")
                        if pl.isPlaying { pl.pause() }
                        else {
                            pl.play()
                        }
                    } else {
                        print("no player, \(w.audioUrl)")
                    }
                }) { RemoteImage(url: w.coverUrl)
                    .scaledToFill()
                    .aspectRatio(contentMode: .fill)
                    .frame(width: UIScreen.main.bounds.size.width, height: 200)
                    .clipped()
                }
                .foregroundColor(.black)
                .overlay(imgOverlay)
                .onTapGesture {
                    print("poop")
                }
                VStack(alignment: .leading, spacing: 5) {
                    Text(w.title).font(.title3)
                    Text(w.description).font(.body).foregroundColor(.gray)
                }
                .padding(.horizontal, 10)
                Spacer()
            }
            .onAppear {
                //            do {
                //                try AVAudioSession.sharedInstance().setCategory(.playback, mode: .default, options: [])
                //            } catch {
                //                print("Setting category to AVAudioSessionCategoryPlayback failed.")
                //            }
                model.getWalk()
                print("walk received")
            }
            .onChange(of: model.loading, perform: { _ in
                print("model loading chg \(model.loading)")
                if let w = model.walk, w.audioUrl != "fakeurl" {
                    print("audio url: \(w.audioUrl)")
                    let playerItem = AVPlayerItem(url: URL(string: w.audioUrl)!)
                    self.player = AVPlayer(playerItem: playerItem)
                }
            })
        } else {
            ProgressView("Walk")
        }
    }
}

struct WalkView_Previews: PreviewProvider {
    static var previews: some View {
        let walk = Api.WalksItem(title: "Tour of the thing", description: "I did a thing that's really great yeah i was there a few times in November last year blablabla.", uuid: "uuid", coverUrl: "https://picsum.photos/200", author: Api.WalkAuthor(uuid: "uuid1", handle: "theo"), distanceFromPoint: 300)
        let vm = WalkViewModel(input: WalkViewModel.Input.FeedWalk(walk), api: Api(token: "xxx"))
        let vm2 = WalkViewModel(input: WalkViewModel.Input.Uuid("ieijie"), api: Api(token: "xxx"))

        return Group {
            WalkView(model: vm)
            WalkView(model: vm2)
        }
    }
}

extension AVPlayer {
    var isPlaying: Bool {
        return rate != 0 && error == nil
    }
}
