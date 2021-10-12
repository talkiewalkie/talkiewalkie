//
//  MemeFont.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct MemeFontUILabel: UIViewRepresentable {
    var frame: CGRect

    func getStringAttributes() -> [NSAttributedString.Key: Any] {
        let paragraphStyle = NSMutableParagraphStyle()
        paragraphStyle.alignment = .left

        let attributes: [NSAttributedString.Key: Any] = [
            NSAttributedString.Key.font: UIFont(name: "Futura Bold Oblique", size: 14) ?? UIFont.systemFont(ofSize: 50),
            NSAttributedString.Key.paragraphStyle: paragraphStyle,

            NSAttributedString.Key.foregroundColor: UIColor.white,
            NSAttributedString.Key.backgroundColor: UIColor.clear,
            NSAttributedString.Key.strokeColor: UIColor.black,
            NSAttributedString.Key.strokeWidth: NSNumber(value: -3.5),
        ]

        return attributes
    }

    func makeUIView(context _: Context) -> UILabel {
        let label = UILabel(frame: frame)
        label.attributedText = NSAttributedString(string: "Good morning, Starshine! The world says Hello. bla bla bla bla bla bla",
                                                  attributes: getStringAttributes())

        label.numberOfLines = 0
        label.setContentCompressionResistancePriority(.defaultLow, for: .horizontal)
        return label
    }

    func updateUIView(_: UILabel, context _: Context) {}
}

struct MemeFontText: View {
    var body: some View {
        GeometryReader { geom in
            MemeFontUILabel(frame: geom.frame(in: .global))
        }
    }
}

struct MemeFont_Previews: PreviewProvider {
    static var previews: some View {
        ZStack {
            Color.orange
                .opacity(0.8)
                .edgesIgnoringSafeArea(.all)

            MemeFontText()
                .background(Color.red)
        }
    }
}
