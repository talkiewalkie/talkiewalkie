//
//  EnabledButton.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct EnabledView<Content: View>: View {
    var enabled: Bool = true
    var disabledAction: (() -> Void)?
    let content: () -> Content

    var body: some View {
        ZStack {
            content()
                .allowsHitTesting(enabled)
        }
        .contentShape(Rectangle())
        .onTapGesture {
            if !enabled {
                disabledAction?()
            }
        }
    }
}


struct EnabledButton<Label>: View where Label: View {
    var enabled: Bool = true
    var disabledAction: (() -> Void)?
    var action: () -> Void
    var opacity: Double = 0.4
    var label: () -> Label

    var body: some View {
        Button(action: {
            if enabled {
                action()
            } else {
                disabledAction?()
            }
        }, label: label)
        .opacity(enabled ? 1 : opacity)
    }
}

struct EnabledButton_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
    }

    struct TestView: View {
        @State var enabled = true

        @State private var showAlert = false

        var body: some View {
            VStack {
                Toggle("Enabled", isOn: $enabled)

                EnabledButton(enabled: enabled, disabledAction: {
                    showAlert = true
                }, action: {}) {
                    Text("Lorem Ipsum")
                }
                .alert(isPresented: $showAlert) {
                    Alert(title: Text("Lorem Ipsum"))
                }
            }
        }
    }
}
