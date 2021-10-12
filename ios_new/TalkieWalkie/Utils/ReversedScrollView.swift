//
//  ReversedScrollView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 06.10.21.
//

import SwiftUI

/// Reference:  https://www.thirdrocktechkno.com/blog/implementing-reversed-scrolling-behaviour-in-swiftui/

struct ReversedScrollView<Content: View>: View {
    var axis: Axis.Set
    var leadingSpace: CGFloat
    var showsIndicators: Bool = false
    var content: Content

    init(_ axis: Axis.Set = .horizontal, leadingSpace: CGFloat = 10, showsIndicators: Bool = false, @ViewBuilder builder: () -> Content) {
        self.axis = axis
        self.leadingSpace = leadingSpace
        self.showsIndicators = showsIndicators
        content = builder()
    }

    var body: some View {
        GeometryReader { proxy in
            ScrollView(axis, showsIndicators: showsIndicators) {
                Stack(axis) {
                    Spacer(minLength: leadingSpace)
                    content
                }
                .frame(
                    minWidth: minWidth(in: proxy, for: axis),
                    minHeight: minHeight(in: proxy, for: axis)
                )
            }
        }
    }

    func minWidth(in proxy: GeometryProxy, for axis: Axis.Set) -> CGFloat? {
        axis.contains(.horizontal) ? proxy.size.width : nil
    }

    func minHeight(in proxy: GeometryProxy, for axis: Axis.Set) -> CGFloat? {
        axis.contains(.vertical) ? proxy.size.height : nil
    }

    struct Stack<Content: View>: View {
        var axis: Axis.Set
        var content: Content

        init(_ axis: Axis.Set = .vertical, @ViewBuilder builder: () -> Content) {
            self.axis = axis
            content = builder()
        }

        var body: some View {
            switch axis {
            case .horizontal:
                HStack {
                    content
                }
            case .vertical:
                VStack {
                    content
                }
            default:
                VStack {
                    content
                }
            }
        }
    }
}

struct ReversedScrollView_Previews: PreviewProvider {
    static var previews: some View {
        ReversedScrollView(.vertical, leadingSpace: 50) {
            ForEach(0 ..< 12) { item in
                Text("\(item)")
                    .padding()
                    .background(Color.gray.opacity(0.5))
                    .cornerRadius(6)
            }
        }
    }
}
