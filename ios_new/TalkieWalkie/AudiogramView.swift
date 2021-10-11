//
//  AudiogramView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 07.10.21.
//

import SwiftUI

struct AudiogramView: View {
    var transcript: Transcript
    
    @State var displayedText: String = ""
    
    var body: some View {
        ZStack {
            Image("background-image")
                .resizable()
                .aspectRatio(contentMode: .fill)
                .saturation(0.1)
                .grayscale(0.5)
                .brightness(-0.15)
            
            ZStack {
                ZStack {
                    ScalingShape(color: .red, scale: 0.95, duration: 0.3) {
                        RoundedRectangle(cornerRadius: 80)
                            .aspectRatio(1, contentMode: .fit)
                            .rotationEffect(.init(degrees: 30))
                    }
                    .scaleEffect(0.55)
                    
                    ScalingShape(color: .red, scale: 0.95, duration: 0.35) {
                        RoundedRectangle(cornerRadius: 80)
                            .aspectRatio(1, contentMode: .fit)
                            .rotationEffect(.init(degrees: 10))
                    }
                    .scaleEffect(0.6)
                }
                .opacity(0.4)
                
                Text(displayedText)
                    .font(.largeTitle)
                    .fontWeight(.bold)
                    .colorInvert()
            }
            .padding()
            .onAppear {
                animate()
            }
        }
        .frame(width: 300, height: 300)
        .cornerRadius(15)
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
        
        transcript.segments.forEach { segment in
            DispatchQueue.main.asyncAfter(deadline: .now() + Double(segment.timestamp)) {
                self.displayedText = segment.string
            }
            
            DispatchQueue.main.asyncAfter(deadline: .now() + transcript.totalDuration, execute: animate)
        }
    }
}

struct AudiogramView_Previews: PreviewProvider {
    static var previews: some View {
        AudiogramView(transcript: dummyTranscript)
    }
}

