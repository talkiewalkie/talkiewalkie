//
//  RecordSheetView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 20.10.21.
//

import SwiftUI

struct RecordSheetView: View {
    @State var isRecording: Bool = false
    @State var showHandle: Bool = false
    
    @State var offset: CGFloat = .zero
    
    let MIN_HEIGHT: CGFloat = 100
    
    var body: some View {
        ZStack {
            Group {
                ZStack(alignment: .top) {
                    handle
                        .opacity(isRecording ? 1 : 0)

                    VStack {
                        VStack {
                            waveform
                                .padding(.vertical, 15)
                                .frame(maxHeight: 20 + 2 * 15)
                        }
                        .frame(maxHeight: isRecording ? MIN_HEIGHT - min(offset, 0) : 0)
                        
                        RecordButton(isRecording: $isRecording, animation: .easeInOut)
                            .padding(5)
                    }
                    .padding()
                }
                .frame(maxWidth: .infinity)
                .background(
                    VisualEffectView(effect: UIBlurEffect(style: isRecording ? .systemUltraThinMaterial : .systemThickMaterial))
                        .opacity(isRecording ? 0.95 : 1)
                        .cornerRadius(20, corners: [.topLeft, .topRight])
                        .edgesIgnoringSafeArea(.vertical)
                )
                .shadow(color: .black.opacity(0.2), radius: 30, y: 0)
                .contentShape(Rectangle())
                .offset(y: max(offset, 0))
                .gesture(
                    DragGesture(coordinateSpace: .global)
                        .onChanged { drag in
                            guard isRecording else { return }
                            offset = max(min(drag.translation.height, 20), -20)
                        }
                        .onEnded { _ in
                            withAnimation(.easeInOut) {
                                offset = .zero
                            }
                        }
                )
                
            }
        }
    }
    
    
    var handle: some View {
        Capsule().foregroundColor(.gray.opacity(0.75))
            .frame(width: 45, height: 6)
            .padding(8)
            .frame(maxWidth: .infinity)
    }
    
    var waveform: some View {
        HStack(spacing: 5) {
            ForEach(1..<35) { _ in
                Capsule().foregroundColor(.red)
                    .frame(width: 4)
            }
        }
    }
}

struct RecordSheetView_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
    }
    
    struct TestView: View {
        var body: some View {
            ZStack(alignment: .bottom) {
                List(1..<20) { i in
                    Text("\(i)")
                }
                
                
                RecordSheetView()
            }
        }
    }
}
