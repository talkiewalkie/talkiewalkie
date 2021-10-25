//
//  AudioWaveView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 25.10.21.
//

import SwiftUI

struct AudioWaveView: View {
    var amplitude: Double
    var isPlaying: Bool
    
    @State var phase: Double = 0
    @State var timer: Timer?
    
    var body: some View {
        GeometryReader { geom in
            ZStack(alignment: .bottom)  {
                AudioWave2(frequency: 10, phase: phase)
                    .fill()
                    .foregroundColor(.red)
                    .opacity(0.6)
                    .frame(height: amplitude * geom.size.height)
                
                AudioWave2(frequency: 15, phase: phase + 40)
                    .fill()
                    .foregroundColor(.red)
                    .opacity(0.2)
                    .frame(height: amplitude * geom.size.height / 2)
                
                AudioWave2(frequency: 5, phase: phase + 20)
                    .fill()
                    .foregroundColor(.red)
                    .opacity(0.2)
                    .frame(height: amplitude * geom.size.height / 3)
            }
            .animation(.linear(duration: 1), value: phase)
            .frame(height: geom.size.height, alignment: .bottom)
        }
        .padding(.bottom, 10)
        .onChange(of: isPlaying) { newValue in
            timer?.invalidate()
            
            if newValue {
                phase += .pi * 2
                timer = Timer.scheduledTimer(withTimeInterval: 1, repeats: true) { _ in
                    phase += .pi * 2
                }
            }
        }
    }
}


struct AudioWave2: Shape {
    var frequency: Double
    var phase: Double
    
    var animatableData: Double {
        get { phase }
        set { phase = newValue }
    }
    
    func path(in rect: CGRect) -> Path {
        var path = Path()
        
        let width = Double(rect.width)
        let height = Double(rect.height)
        let midHeight = height / 2
        
        let wavelength = width / frequency
        
        for x in stride(from: 0, to: width, by: 1) {
            let relativeX = x / wavelength
            let sine = sin(relativeX + phase)
            let y = midHeight * (sine + 1)
            
            if x == 0 {
                path.move(to: CGPoint(x: x, y: y))
            } else {
                path.addLine(to: CGPoint(x: x, y: y))
            }
        }
        
        path.addLine(to: CGPoint(x: width, y: height + 10))
        path.addLine(to: CGPoint(x: 0, y: height + 10))
        path.closeSubpath()
        
        return path
    }
}


struct AudioWaveView_Previews: PreviewProvider {
    static var previews: some View {
        TestView2()
    }
    
    struct TestView2: View {
        @State var isPlaying: Bool = false
        @State private var amplitude: Double = 1
        
        var body: some View {
            ZStack(alignment: .bottom) {
                Color.clear
                
                AudioWaveView(amplitude: amplitude, isPlaying: isPlaying)
            }
            .frame(height: 100, alignment: .bottom)
            
            .onAppear {
                withAnimation(Animation.easeInOut(duration: 1).repeatForever(autoreverses: true)) {
                    amplitude = 0.2
                }
                
                isPlaying = true
            }
            .ignoresSafeArea()
        }
    }
    
//    struct TestView: View {
//        @State private var amplitude: Double = 100
//        @State private var phase: Double = 0
//
//        var body: some View {
//            ZStack {
//                Color.clear
//
//                VStack {
//                    ZStack(alignment: .bottom) {
//                        AudioWave2(frequency: 10, phase: phase)
//                            .fill()
//                            .foregroundColor(.red)
//                            .opacity(0.6)
//                            .frame(height: amplitude)
//
//                        AudioWave2(frequency: 15, phase: phase + 40)
//                            .fill()
//                            .foregroundColor(.red)
//                            .opacity(0.2)
//                            .frame(height: amplitude / 2)
//
//                        AudioWave2(frequency: 5, phase: phase + 20)
//                            .fill()
//                            .foregroundColor(.red)
//                            .opacity(0.2)
//                            .frame(height: amplitude / 3)
//                    }
//                    .padding(.vertical)
//                    .frame(height: 200, alignment: .bottom)
//
//                }
//            }
//            .onAppear {
//                withAnimation(Animation.easeInOut(duration: 1).repeatForever(autoreverses: true)) {
//                    amplitude = 25
//                }
//
//                withAnimation(Animation.linear(duration: 1).repeatForever(autoreverses: false)) {
//                    phase = .pi * 2
//                }
//            }
//            .ignoresSafeArea()
//        }
//    }
}
