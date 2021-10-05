//
//  Refreshable.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 17/09/2021.
//

import SwiftUI

// https://dev.to/gualtierofr/pull-down-to-refresh-in-swiftui-4j26
struct RefreshableScrollView<Content: View>: View {
    init(action: @escaping () -> Void, @ViewBuilder content: @escaping () -> Content) {
        self.content = content
        refreshAction = action
    }

    var body: some View {
        GeometryReader { geometry in
            ScrollView {
                content()
                    .anchorPreference(key: OffsetPreferenceKey.self, value: .top) {
                        geometry[$0].y
                    }
            }
            .onPreferenceChange(OffsetPreferenceKey.self) { offset in
                if offset > threshold, Date() > lastActionOccuredAt.addingTimeInterval(cancellationWindow / 1000) {
                    refreshAction()
                    lastActionOccuredAt = Date()
                }
            }
        }
    }

    // MARK: - Private

    private var content: () -> Content
    private var refreshAction: () -> Void
    private let threshold: CGFloat = 50.0
    private let cancellationWindow: Double = 2000
    @State private var lastActionOccuredAt = Date()
}

private struct OffsetPreferenceKey: PreferenceKey {
    static var defaultValue: CGFloat = 0

    static func reduce(value: inout CGFloat, nextValue: () -> CGFloat) {
        value = nextValue()
    }
}
