//
//  CloseButton.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct CloseButtonView: View {
    var body: some View {
        Image(systemName: "xmark")
            .font(.system(size: 11))
            .foregroundColor(.white)
            .padding(7)
            .background(
                Circle()
                    .foregroundColor(.black)
                    .opacity(0.4)
            )
            .shadow(color: .white, radius: 2)
            .padding(5)
            .contentShape(Circle())
    }
}

struct CloseButton: View {
    @Binding var show: Bool
    
    var animation: Animation?

    var body: some View {
        Button(action: {
            guard let animation = animation else { return show = false }
            
            withAnimation(animation) {
                show = false
            }
        }) {
           CloseButtonView()
        }
    }
}

struct CloseButton_Previews: PreviewProvider {
    static var previews: some View {
        ZStack {
            Color.gray

            CloseButton(show: Binding.constant(true))
        }
    }
}
