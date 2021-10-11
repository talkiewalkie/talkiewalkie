//
//  RecordButton.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct RecordButton: View {
    @Binding var isRecording: Bool
    
    @State var isOuterCircleThin = true
    @State var timer: Timer?
    
    var body: some View {
        ZStack {
            Rectangle()
                .fill(Color.red)
                .cornerRadius(isRecording ? 10 : 30)
                .frame(width: 60, height: 60)
                .scaleEffect(isRecording ? 0.7 : 1.0)
                .animation(.easeInOut(duration: 0.3), value: isRecording)
            
            Circle()
                .stroke(Color.red,
                        style: StrokeStyle(lineWidth: !isRecording ? 6 : isOuterCircleThin ? 4 : 7,
                                           lineCap: .round, lineJoin: .round))
                .animation(.easeOut(duration: 0.45), value: isOuterCircleThin)
                .opacity(0.7)
                .frame(width: 70, height: 70)
                .scaleEffect(isRecording ? 1.3 : 1.0)
                .animation(.easeInOut(duration: 0.3), value: isRecording)

        }
        .onTapGesture {
//            let currentClipLength = recordViewModel.sequenceRecorder.currentClipLength
//
//            if isRecording && currentClipLength < 0.25 {
//                DispatchQueue.main.asyncAfter(deadline: .now() + 0.25 - currentClipLength) {
//                    isRecording.toggle()
//                }
//                return
//            }
            
            isRecording.toggle()
        }
        .onChange(of: isRecording) { value in
            timer?.invalidate()
            
            if !isRecording {
                isOuterCircleThin = true
            } else {
                isOuterCircleThin.toggle()
                timer = Timer.scheduledTimer(withTimeInterval: 0.45, repeats: true) { _ in
                    isOuterCircleThin.toggle()
                }
            }
        }
    }
}

struct RecordButton_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
    }
    
    struct TestView: View {
        @State var isRecording = false
        
        var body: some View {
            RecordButton(isRecording: $isRecording)
        }
    }
}
