//
//  PageCounter.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct PageCounter: View {
    @Binding var page: Int
    var pages: Int

    var size: CGFloat = 15
    var spacing: CGFloat = 15

    var body: some View {
        ZStack(alignment: .leading) {
            HStack(spacing: spacing) {
                ForEach(0 ..< pages) { _ in
                    Circle()
                        .fill(Color(#colorLiteral(red: 0.8943424931, green: 0.8687898505, blue: 0.8687898505, alpha: 1)))
                        .frame(width: size, height: size)
                        .shadow(color: Color.black.opacity(0.1), radius: 1)
                }
            }

            let offset = (size + spacing) * CGFloat(page)
            Circle()
                .fill(Color.white)
                .frame(width: size, height: size)
                .shadow(radius: 1)
                .offset(x: offset)
                .animation(.default)
        }
    }
}

struct PageCounter_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
    }

    struct TestView: View {
        @State var page = 0

        var body: some View {
            VStack {
                PageCounter(page: $page, pages: 3)

                Button("Click") {
                    page = (page + 1) % 3
                }
            }
        }
    }
}
