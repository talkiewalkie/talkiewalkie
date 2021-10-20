//
//  RecordButton.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI
import Lottie

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
        // .contentShape(Rectangle())
        .onLongPressGesture(minimumDuration: 10000, maximumDistance: 1000, perform: { }, onPressingChanged: { pressing in
            isRecording.toggle()
        })
        .simultaneousGesture(
            TapGesture()
                .onEnded { _ in
                    isRecording.toggle()
                }
        )
        .onChange(of: isRecording) { _ in
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
        
        // .padding()
        .frame()
        .contentShape(Rectangle())
        .overlay(Rectangle().foregroundColor(.yellow).opacity(0.2))
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
