//
//  Audiogram2.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 04.10.21.
//

import SwiftUI

struct Audiogram2: View {
    var text: String

    @State var displayedText: String = ""

    var body: some View {
        ZStack {
            Image("chuck")
                .resizable()
                .aspectRatio(contentMode: .fill)

            ZStack {
                ScalingShape(color: .red, scale: 0.95, duration: 0.3) {
                    RoundedRectangle(cornerRadius: 80)
                        .aspectRatio(1, contentMode: .fit)
                        .rotationEffect(.init(degrees: 30))
                }
                .scaleEffect(0.6)

                ScalingShape(color: .red, scale: 0.95, duration: 0.4) {
                    RoundedRectangle(cornerRadius: 80)
                        .aspectRatio(1, contentMode: .fit)
                        .rotationEffect(.init(degrees: 10))
                }
                .scaleEffect(0.6)
            }
            .opacity(0.3)

            Text(displayedText)
                .font(.largeTitle)
                .fontWeight(.medium)
                .colorInvert()
        }
        .frame(width: 300, height: 300)
        .clipped()
        .onAppear {
            animate()
        }
        .onTapGesture {
            animate()
        }
    }

    func animate() {
        let words = text.components(separatedBy: " ")
        words.enumerated().forEach { index, word in
            DispatchQueue.main.asyncAfter(deadline: .now() + Double(index) * 0.3) {
                self.displayedText = word
            }
        }
    }
}

struct ScalingShape<ShapeView>: View where ShapeView: View {
    var color: Color
    var scale: CGFloat
    var duration: Double
    var shape: () -> ShapeView

    @State var currScale: CGFloat = 1

    var body: some View {
        shape()
            .foregroundColor(color)
            .scaleEffect(currScale)
            .animation(.easeInOut(duration: duration).repeatForever(autoreverses: true))
            .onAppear {
                currScale = scale
            }
    }
}

struct Audiogram2_Previews: PreviewProvider {
    static var previews: some View {
        Audiogram2(text: "Time waits for no man; Unless that man is Chuck Norris. Chuck Norris does not sleep; He waits. Chuck Norris can dribble a bowling ball.")
    }
}
