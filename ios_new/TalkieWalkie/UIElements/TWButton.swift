//
//  TWButton.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct TWButton<Content>: View where Content: View {
    var action: () -> Void
    var primary: Bool = true
    var compact: Bool = true
    
    var font: Font = .callout.weight(.heavy)
    var padding: CGFloat = 15
    
    var content: () -> Content
    
    @State var scale: CGFloat = 1
    
    var body: some View {
        Button(action: action) {
            content()
                .foregroundColor(primary ? .white : Color("DarkBlue"))
                .font(font)
                .frame(maxWidth: compact ? .none : .infinity)
                .padding(padding)
                .background(
                    Group {
                        primary ? Color("Bordeaux") : Color("LightBlue")
                    }.cornerRadius(15))
                
        }
        .buttonStyle(ScaleButtonStyle())
    }
}

struct ScaleButtonStyle: ButtonStyle {
    func makeBody(configuration: Self.Configuration) -> some View {
        configuration.label
            .scaleEffect(configuration.isPressed ? 0.92 : 1)
    }
}

struct TWButton_Previews: PreviewProvider {
    static var previews: some View {
        TWButton(action: {}, primary: false, compact: false) {
            Text("Continue".uppercased())
        }
    }
}
