//
//  WrapStack.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct WrapHStack<Content: View>: View {
    var spacing: CGFloat = 4
    var content: [Content]

    @State private var size: CGSize = .zero

    var body: some View {
        GeometryReader { geom in
            var width = CGFloat.zero
            var height = CGFloat.zero

            ZStack(alignment: .topLeading) {
                ForEach(content.indices, id: \.self) { index in
                    content[index]
                        .padding(spacing)
                        .alignmentGuide(.leading) { d in
                            if abs(width - d.width) > geom.size.width {
                                width = 0
                                height -= d.height
                            }
                            let result = width
                            if index == content.indices.count - 1 {
                                width = 0
                            } else {
                                width -= d.width
                            }
                            return result
                        }
                        .alignmentGuide(.top) { _ in
                            let result = height
                            if index == content.indices.count - 1 {
                                height = 0
                            }
                            return result
                        }
                }
            }
            .background(sizeReader($size))
        }
        .frame(height: size.height)
    }
}

struct WrapStack_Previews: PreviewProvider {
    static var previews: some View {
        ScrollView {
            WrapHStack(content: [
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
                WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"), WrapTestView(text: "Lorem Ipsum"),
            ])
        }
    }

    struct WrapTestView: View {
        var text: String

        var body: some View {
            Text(text).padding()
                .background(RoundedRectangle(cornerRadius: 8).fill(Color.red))
                .foregroundColor(.white)
        }
    }
}
