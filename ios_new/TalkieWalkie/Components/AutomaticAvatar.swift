//
//  AutomaticAvatar.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 17/10/2021.
//

import SwiftUI

struct AutomaticAvatar: View {
    var letter: String = "?"
    var color: UIColor?
    var size: CGFloat

    init(_ letter: String, size: CGFloat = 48, color: UIColor? = nil) {
        self.letter = letter
        self.color = color
        self.size = size
    }

    var body: some View {
        Group {
            ZStack {
                Color(color ?? generateColorFor(text: letter))
                    .brightness(-0.1)

                Text(letter)
                    .font(.title2)
                    .fontWeight(.medium)
                    .foregroundColor(.white)
            }
            .aspectRatio(1, contentMode: .fit)
            .frame(width: size)
        }

        .clipShape(Circle())
    }
}

struct AutomaticAvatar_Previews: PreviewProvider {
    static var previews: some View {
        AutomaticAvatar("T", size: 24)
    }
}
