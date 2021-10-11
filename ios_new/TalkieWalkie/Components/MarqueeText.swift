//
//  MarqueeText.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

extension String {
   func widthOfString(usingFont font: UIFont) -> CGFloat {
        let fontAttributes = [NSAttributedString.Key.font: font]
        let size = self.size(withAttributes: fontAttributes)
        return size.width
    }
}

struct MarqueeText: View {
    var text: String
    var spacing: CGFloat = 30
    var maxWidth: CGFloat = 80
    var textSize: CGFloat = 18

    @State var offset: CGFloat = .zero

    @State var timer : Timer?

    var textWidth: CGFloat {
        text.widthOfString(usingFont: UIFont.systemFont(ofSize: textSize))
    }

    var body: some View {
        let marqueeWidth = min(maxWidth, textWidth + spacing)
        
        ZStack {
            Text(text).font(.system(size: textSize)).fixedSize()
                .offset(x: offset)

            Text(text).font(.system(size: textSize)).fixedSize()
                .offset(x: offset + textWidth + spacing)
        }

        .frame(width: marqueeWidth, alignment: .leading)
        .clipped()
        .onAppear {
            startAnimation()
        }
        .onChange(of: text, perform: { _ in
            startAnimation()
        })
    }

    func animate() {
        withAnimation(Animation.linear(duration: 2.5)) {
            offset = -(textWidth + spacing)
        }
    }

    func startAnimation() {
        timer?.invalidate()

        offset = .zero
        animate()

        timer = Timer.scheduledTimer(withTimeInterval: 2.5, repeats: true) { timer in
            offset = .zero
            animate()
        }
    }
}

struct MarqueeText_Previews: PreviewProvider {
    static var previews: some View {
        MarqueeText(text: "Bla bla bla")
    }
}
