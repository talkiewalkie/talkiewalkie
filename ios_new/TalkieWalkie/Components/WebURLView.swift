//
//  WebURLView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct WebURLView: View {
    var url: URL?
    var displayURL: String?

    init(urlString: String) {
        guard let url = URL(string: urlString), let host = url.host else { return }
        self.url = url
        displayURL = host
    }

    var body: some View {
        if url != nil && displayURL != nil {
            Link(destination: url!, label: {
                Text(displayURL!)
                    .font(.callout)
                    .padding(.horizontal, 20).padding(.vertical, 8)
                    .background(Capsule().stroke(Color.secondary))
            })
        }
    }
}

struct WebURLView_Previews: PreviewProvider {
    static var previews: some View {
        ZStack {
            Color.black

            WebURLView(urlString: "https://www.mirror.co.uk/tv/gallery/harry-potter-through-the-years-3827675")
                .environment(\.colorScheme, .dark)
        }
    }
}
