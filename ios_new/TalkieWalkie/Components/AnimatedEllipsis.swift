//
//  AnimatedEllipsis.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct AnimatedEllipsis: View {
    var numberOfDots = 3
    var minOpactiy: Double = 0.2
    var animationDelay: Double = 0.2
    var color = Color.primary
    var fontSize: CGFloat = 20
    
    @State private var opacity: Double = 1
    
    var body: some View {
        HStack(spacing: 0) {
            ForEach(0..<3) { index in
                Text(" .")
                    .font(.system(size: fontSize))
                    .kerning(-(fontSize / 20))
                    .foregroundColor(color)
                    .opacity(opacity)
                    .animation(
                        Animation.easeInOut(duration: animationDelay)
                            .delay(animationDelay * Double(index))
                    )
            }
        }
        .onAppear {
            Timer.scheduledTimer(withTimeInterval: animationDelay * Double(numberOfDots + 1), repeats: true) { _ in
                withAnimation {
                    opacity = 1 - opacity + minOpactiy
                }
            }
        }
    }
}

struct AnimatedEllipsis_Previews: PreviewProvider {
    static var previews: some View {
        AnimatedEllipsis()
    }
}
