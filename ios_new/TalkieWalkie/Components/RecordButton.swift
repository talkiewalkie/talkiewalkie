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

    var animation: Animation?
    var buttonAnimation: Animation = .easeInOut(duration: 0.3)

    func toggleRecording() {
        withAnimation(animation) {
            isRecording.toggle()
        }
    }

    var body: some View {
        ZStack {
            Rectangle()
                .fill(Color.red)
                .cornerRadius(isRecording ? DrawingConstraints.innerCircleRadius : DrawingConstraints.innerCircleSize / 2)
                .frame(width: DrawingConstraints.innerCircleSize, height: DrawingConstraints.innerCircleSize)
                .scaleEffect(isRecording ? DrawingConstraints.innerCircleScale : 1.0)
                .animation(buttonAnimation, value: isRecording)

            Circle()
                .strokeBorder(Color.red,
                              style: StrokeStyle(lineWidth: !isRecording ? 6 : isOuterCircleThin ? 4 : 7,
                                                 lineCap: .round, lineJoin: .round))
                .animation(.easeOut(duration: 0.4), value: isOuterCircleThin)
                .opacity(DrawingConstraints.outerCircleOpacity)
                .frame(width: DrawingConstraints.outerCircleSize, height: DrawingConstraints.outerCircleSize)
                .scaleEffect(isRecording ? DrawingConstraints.outerCircleScale : 1.0)
                .animation(buttonAnimation, value: isRecording)
        }
        .padding(5)
        .contentShape(Circle().scale(DrawingConstraints.tapAreaScale))
        .onLongPressGesture(minimumDuration: 10000, maximumDistance: 1000, perform: {}, onPressingChanged: { _ in
            toggleRecording()
        })
        .simultaneousGesture(
            TapGesture()
                .onEnded { _ in
                    toggleRecording()
                }
        )
        .onChange(of: isRecording) { _ in
            timer?.invalidate()

            if !isRecording {
                isOuterCircleThin = true
            } else {
                isOuterCircleThin.toggle()
                timer = Timer.scheduledTimer(withTimeInterval: 0.4, repeats: true) { _ in
                    isOuterCircleThin.toggle()
                }
            }
        }
    }

    enum DrawingConstraints {
        static let innerCircleSize: CGFloat = 60
        static let innerCircleRadius: CGFloat = 10
        static let innerCircleScale: CGFloat = 0.7

        static let outerCircleOpacity: Double = 0.6
        static let outerCircleSize: CGFloat = 75
        static let outerCircleScale: CGFloat = 1.75

        static let tapAreaScale: CGFloat = 1.75
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
