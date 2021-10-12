//
//  CirclesBackground.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 04.10.21.
//

import SwiftUI

struct CirclesBackground: View {
    var body: some View {
        ZStack {
            MovingCircle(color: Color(#colorLiteral(red: 0.9529411793, green: 0.6862745285, blue: 0.1333333403, alpha: 1)), offset: .init(width: -110, height: +100), scale: 0.9, duration: 4)
                .scaleEffect(0.65)
                .offset(x: 20, y: 10)

            MovingCircle(color: Color(#colorLiteral(red: 0.9372549057, green: 0.3490196168, blue: 0.1921568662, alpha: 1)), offset: .init(width: -100, height: -150), scale: 0.9, duration: 2)
                .scaleEffect(0.5)
                .offset(x: -25, y: 20)

            MovingCircle(color: Color(#colorLiteral(red: 0.4666666687, green: 0.7647058964, blue: 0.2666666806, alpha: 1)), offset: .init(width: 50, height: 100), scale: 0.9, duration: 3)
                .scaleEffect(0.4)
                .offset(x: 10, y: -20)

            Text("Text")
                .font(.largeTitle)
                .fontWeight(.medium)
                .scaleEffect(2.0)
                .colorInvert()
                .shadow(color: .black, radius: 25)
        }
    }
}

struct MovingCircle: View {
    var color: Color
    var offset: CGSize
    var scale: CGFloat
    var duration: Double

    @State var currOffset: CGSize = .zero
    @State var currScale: CGFloat = 1

    var body: some View {
        Circle()
            .foregroundColor(color)
            .scaleEffect(currScale)
            .offset(currOffset)
            .blur(radius: 35)
            .animation(.easeInOut(duration: duration).repeatForever(autoreverses: true))
            .onAppear {
                currOffset = offset
                currScale = scale
            }
    }
}

struct CirclesBackground_Previews: PreviewProvider {
    static var previews: some View {
        CirclesBackground()
    }
}
