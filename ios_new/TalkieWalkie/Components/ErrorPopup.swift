//
//  ErrorPopup.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct ErrorPopup: View {
    @Binding var show: Bool

    var message: LocalizedStringKey

    var body: some View {
        if show {
            Text(message)
                .foregroundColor(.white)
                .padding(.vertical, 10).padding(.horizontal)
                .frame(maxWidth: .infinity)
                .multilineTextAlignment(.center)
                .background(Color.red)
                .opacity(0.9)

                .transition(AnyTransition.move(edge: .top).combined(with: .opacity))

                .onTapGesture {
                    show = false
                }
        }
    }
}

struct ErrorPopup_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
    }

    struct TestView: View {
        @State var show = false

        var body: some View {
            ZStack {
                VStack {
                    ErrorPopup(show: $show, message: "Error message")

                    Spacer()
                }
                .padding()

                Button(action: {
                    withAnimation(.default) {
                        show = true
                    }
                }) {
                    Text("Click")
                }
            }
        }
    }
}
