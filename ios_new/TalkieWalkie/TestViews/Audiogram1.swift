//
//  Audiogram1.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 04.10.21.
//

import SwiftUI

struct Audiogram1: View {
    
    var text: String
    
    @State var displayedText: String = ""
    
    var body: some View {
        
        ZStack {
            Image("background-image")
                .resizable()
                .aspectRatio(contentMode: .fill)
                .saturation(0.1)
                .grayscale(0.5)
                .brightness(-0.15)
            
            ZStack(alignment: .topLeading) {
                
                Text(text)
                    .font(.title)
                    .fontWeight(.bold)
                    .colorInvert()
                    .opacity(0.5)
                
                Text(displayedText)
                    .font(.title)
                    .fontWeight(.bold)
                    .colorInvert()
            }
            .padding()
            .onAppear {
                animate()
            }
            .onTapGesture {
                animate()
            }
        }
        .frame(width: 300, height: 300)
        .overlay(
            Image(systemName: "speaker.wave.2.fill")
                .font(.caption)
                .colorInvert()
                .padding(10)
                .background(Circle())
                .padding(),
            alignment: .topTrailing
        )
    }
    
    func animate() {
        self.displayedText = ""
        
        let words = text.components(separatedBy: " ")
        words.enumerated().forEach { index, word in
            DispatchQueue.main.asyncAfter(deadline: .now() + Double(index) * 0.3) {
                self.displayedText += word + " "
                
                if index == words.count - 1 {
                    DispatchQueue.main.asyncAfter(deadline: .now() + 0.5) {
                        self.animate()
                    }
                }
            }
        }
    }
}

struct Audiogram1_Previews: PreviewProvider {
    static var previews: some View {
        Audiogram1(text: "The reason that happens so much is because once you're starting to")
    }
}
