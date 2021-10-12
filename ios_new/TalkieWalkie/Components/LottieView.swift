//
//  LottieView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import Lottie
import SwiftUI

struct LottieView: UIViewRepresentable {
    var name: String
    var loopMode: LottieLoopMode = .loop
    var contentMode: UIView.ContentMode = .scaleAspectFit

    func makeUIView(context _: Context) -> UIView {
        let view = UIView(frame: .zero)

        let animationView = AnimationView()
        animationView.animation = Animation.named(name)
        animationView.loopMode = loopMode
        animationView.contentMode = contentMode
        animationView.play()
        animationView.translatesAutoresizingMaskIntoConstraints = false

        view.addSubview(animationView)

        NSLayoutConstraint.activate([
            animationView.widthAnchor.constraint(equalTo: view.widthAnchor),
            animationView.heightAnchor.constraint(equalTo: view.heightAnchor),
        ])

        return view
    }

    func updateUIView(_: UIView, context _: Context) {}
}

struct LottieView_Previews: PreviewProvider {
    static var previews: some View {
        LottieView(name: "crown")
    }
}
