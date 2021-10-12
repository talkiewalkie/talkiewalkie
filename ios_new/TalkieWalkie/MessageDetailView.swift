//
//  MessageDetailView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 06.10.21.
//

import SwiftUI

struct MessageDetailView: View {
    var namespace: Namespace.ID

    @EnvironmentObject var messageViewModel: MessageViewModel

    @State var offsetY: CGFloat = .zero

    var body: some View {
        ZStack {
            Color.black
                .ignoresSafeArea()
                .opacity(max(0.7 - 0.2 * abs(offsetY) / 300, 0)) // TODO:

            if let message = messageViewModel.message {
                Group {
                    switch message.type {
                    case let .text(content: text):
                        Audiogram1(text: text)
                            .matchedGeometryEffect(id: message.id.uuidString, in: namespace)
                    default: EmptyView()
                    }
                }
                .frame(maxWidth: .infinity, maxHeight: 300)
                .background(Color.gray)
                .offset(x: 0, y: offsetY)
            }

            if offsetY == .zero {
                CloseButton(show: $messageViewModel.showDetailView) // , animation: .spring())
                    .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .topLeading)
                    .padding()
            }
        }
        .gesture(DragGesture()
            .onChanged { gesture in
                offsetY = gesture.translation.height
            }
            .onEnded { gesture in
                let predictedY = gesture.predictedEndLocation.y
                let screenHeight = UIScreen.main.bounds.size.height

                if predictedY < 150 || predictedY > screenHeight - 150 {
                    let targetOffset = predictedY < 150 ? -screenHeight : screenHeight

                    messageViewModel.message = nil

                    withAnimation(.spring()) {
                        offsetY = targetOffset
                        messageViewModel.showDetailView = false
                    }

                    DispatchQueue.main.asyncAfter(deadline: .now() + 0.2) {
                        offsetY = .zero
                    }
                } else {
                    withAnimation(.spring()) {
                        offsetY = .zero
                    }
                }
            }
        )
    }
}

struct MessageDetailView_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
    }

    struct TestView: View {
        @Namespace var namespace

        var body: some View {
            MessageDetailView(namespace: namespace)
                .withDummyVariables()
        }
    }
}
