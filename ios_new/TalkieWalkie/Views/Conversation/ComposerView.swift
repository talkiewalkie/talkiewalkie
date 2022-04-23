//
//  ComposerView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 20.10.21.
//

import GiphyUISDK
import SwiftUI

enum ComposerState {
    case inactive
    case recording
    case recordingFinished
    case canceled

    var isInactive: Bool {
        return self == .inactive || self == .canceled
    }
}

struct ComposerView: View {
    let conversationUuid: UUID
    @Binding var isTextFieldFocused: Bool

    @State var isRecording: Bool = false
    @State var recordState: ComposerState = .inactive

    @State var audioRecorder = AudioRecorder()
    @State var timer = Timer.publish(every: 0.05, on: .main, in: .common).autoconnect()

    @State var audioPowers = [Float]()

    @State var offset: CGFloat = .zero

    @State var text: String = ""

    @State var showAttachementActionSheet: Bool = false

    @EnvironmentObject var authed: AuthState

    var textfield: some View {
        AutoTextField($text, isFocused: $isTextFieldFocused)
            .placeholder("Aa")
            .setMaxHeight(105)
            .maxLines(isTextFieldFocused ? 0 : 1)
            .multilineTextAlignment(.leading)
            .padding(.vertical, 8)
            .padding(.horizontal, 12)
            .background(Color.gray.opacity(0.15))
            .cornerRadius(24)
            .contentShape(Rectangle())
            .onTapGesture {
                isTextFieldFocused = true
            }
            .frame(maxWidth: isTextFieldFocused ? .infinity : 100)
    }

    var gifIcon: some View {
        Text("GIF")
            .font(.system(size: 10, weight: .heavy))
            .padding(.horizontal, 3)
            .padding(.vertical, 2.5)
            .overlay(
                RoundedRectangle(cornerRadius: 5).stroke(lineWidth: 1.5)
                    .foregroundColor(.secondary)
            )
            .foregroundColor(.secondary)
            .padding(.horizontal, 15)
            .padding(.vertical, 18)
            .contentShape(Rectangle())
            .onTapGesture {
                guard let root = UIApplication.shared.windows.last?.rootViewController else { return }

                let giphy = GiphyViewController()
                giphy.mediaTypeConfig = [.gifs, .recents]
                giphy.theme = GPHTheme(type: .lightBlur)

                root.present(giphy, animated: true, completion: nil)
            }
    }

    var stickerIcon: some View {
        Image("sticker")
            .resizable()
            .aspectRatio(contentMode: .fit)
            .frame(height: 26)
            .foregroundColor(.secondary)
            .padding(15)
            .contentShape(Rectangle())
            .onTapGesture {
                guard let root = UIApplication.shared.windows.last?.rootViewController else { return }

                let giphy = GiphyViewController()
                giphy.mediaTypeConfig = [.stickers]
                giphy.theme = GPHTheme(type: .lightBlur)

                root.present(giphy, animated: true, completion: nil)
            }
    }

    var attachmentIcon: some View {
        Button(action: { showAttachementActionSheet = true }) {
            Image(systemName: "paperclip")
                .font(.title3)
                .foregroundColor(.secondary)
                .contentShape(Rectangle().scale(1.75))
                .padding(.horizontal, 15)
                .padding(.vertical, 2)
        }.actionSheet(isPresented: $showAttachementActionSheet) {
            ActionSheet(title: Text("Add attachment"), buttons: [
                .cancel(),
                .default(Text("Photo or Video"), action: {}),
                .default(Text("File"), action: {}),
                .default(Text("Location"), action: {}),
                .default(Text("Contact"), action: {}),
            ])
        }
    }

    var textSendButton: some View {
        Button(action: {
            let messageLocalUuid = UUID()

            // Core Data
            authed.withWriteContext { ctx, me in
                guard let me = me else { fatalError() }
                let message = Message(context: ctx)
                message.localUuid_ = messageLocalUuid
                message.content = {
                    let content = TextMessage(context: ctx)
                    content.text = text
                    return content
                }()
                message.author = me
                message.status_ = 0
                message.conversation = Conversation.getByUuidOrCreate(conversationUuid, context: ctx)
                message.createdAt = Date()
            }

            // Sync
            if case let .Connected(api, _) = authed.state {
                DispatchQueue.global(qos: .background).async {
                    let localMessage = api.sendMessage(text: text, convUuid: conversationUuid)
                    authed.withWriteContext { ctx, _ in
                        Message.fromEventProto(localMessage, context: ctx)
                    }
                }
            }
        }) {
            Image(systemName: "paperplane.fill")
                .padding(8)
                .offset(x: -1, y: 1)
                .foregroundColor(.white)
                .background(
                    Circle().foregroundColor(.accentColor)
                )
                .contentShape(Rectangle().scale(1.5))
                .padding(.horizontal, 12)
        }
    }

    var body: some View {
        ZStack(alignment: .top) {
            handle
                .opacity(recordState.isInactive ? 0 : 1)

            VStack(spacing: 0) {
                VStack {
                    AudioPowerVisualizer(
                        powers: audioPowers,
                        onTap: self.playHandler
                    )
                }
                .padding(.vertical, 25)
                .frame(height: recordState.isInactive ? 0 : DrawingConstraints.MIN_HEIGHT - min(offset, 0))
                .padding(.horizontal)

                ZStack {
                    HStack(alignment: .bottom, spacing: 0) {
                        attachmentIcon
//                            .background(Color.blue.opacity(0.5))

                        textfield
//                            .background(Color.yellow.opacity(0.5))

                        textSendButton
                            .opacity(isTextFieldFocused ? 1 : 0)
                            .frame(width: isTextFieldFocused ? nil : 0)
//                            .background(Color.purple.opacity(0.5))
                    }
                    .padding(.vertical, 5)
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .opacity(recordState.isInactive ? 1 : 0)

                    HStack(spacing: 0) {
                        gifIcon

                        stickerIcon
                    }
                    .opacity(isTextFieldFocused ? 0 : 1)
                    .offset(x: isTextFieldFocused ? 100 : 0)
                    .frame(maxWidth: .infinity, alignment: .trailing)
                    .opacity(recordState.isInactive ? 1 : 0)

                    RecordButton(isRecording: $isRecording)
                        .padding(5)
                        .opacity(isTextFieldFocused ? 0 : 1)
                        .frame(maxWidth: .infinity, alignment: isTextFieldFocused ? .trailing : .center)
                        .frame(height: isTextFieldFocused ? 0 : nil)

                    HStack(spacing: 10) {
                        cancelButton
                            .opacity(recordState == .recordingFinished ? 1 : 0)

                        audioSendButton
                            .opacity(recordState == .recordingFinished ? 1 : 0)
                    }
                    .frame(maxWidth: .infinity, alignment: .trailing)
                }
            }
        }
        .animation(.easeInOut(duration: 0.15), value: isTextFieldFocused)
        .frame(maxWidth: .infinity)
        .background(
            VisualEffectView(effect: UIBlurEffect(
                style: recordState.isInactive ? .systemThickMaterial : .systemUltraThinMaterial)
            )
            .opacity(recordState.isInactive ? 1 : 0.95)
            .cornerRadius(20, corners: [.topLeft, .topRight])
            .edgesIgnoringSafeArea(.vertical)
        )
        .shadow(color: .black.opacity(0.1), radius: 10)
        .offset(y: max(offset, 0))
        .gesture(
            DragGesture(coordinateSpace: .global)
                .onChanged { drag in
                    guard !recordState.isInactive else { return }
                    offset = max(drag.translation.height, -20)
                }
                .onEnded { drag in
                    withAnimation(.easeInOut) {
                        offset = .zero

                        if drag.predictedEndTranslation.height > DrawingConstraints.SWIPE_TO_DISCARD_THRESHOLD {
                            isRecording = false

                            audioRecorder.stopRecording()
                            audioPowers.removeAll()

                            recordState = .canceled
                        }
                    }
                }
        )
        .onChange(of: isRecording) { _ in
            if isRecording {
                audioRecorder.record()

                withAnimation(.easeInOut) {
                    recordState = .recording
                }

            } else {
                if recordState == .canceled {
                    withAnimation(.easeInOut) {
                        recordState = .inactive
                    }
                    return
                }

                audioRecorder.stopRecording()

                withAnimation(.easeInOut) {
                    recordState = .recordingFinished
                }
            }
        }
        .onReceive(timer) { _ in
            if isRecording {
                audioRecorder.recorder.updateMeters()

                let power = audioRecorder.recorder.averagePower(forChannel: 0)
                let scaledPower = sigmoid((power + 17) / 1.3) / 0.96 + 0.04

                audioPowers.append(scaledPower)
            }
        }
    }

    func sigmoid(_ z: Float) -> Float {
        return 1.0 / (1.0 + exp(-z))
    }

    func playHandler() {
        if audioRecorder.state == .stopped {
            audioRecorder.play()
        } else if audioRecorder.state == .playing {
            audioRecorder.stopPlayback()
        }
    }

    var handle: some View {
        Capsule().foregroundColor(.gray.opacity(0.75))
            .frame(width: 45, height: 6)
            .padding(8)
            .frame(maxWidth: .infinity)
    }

    var audioSendButton: some View {
        Button(action: {
            audioPowers.removeAll()

            withAnimation(.easeInOut) {
                recordState = .inactive
            }
        }) {
            Image(systemName: "checkmark")
                .font(.title3)
                .foregroundColor(.white)
                .padding(10)
                .background(Circle().foregroundColor(.red))
                .padding(5)
        }
    }

    var cancelButton: some View {
        Button(action: {
            audioPowers.removeAll()

            withAnimation(.easeInOut) {
                recordState = .inactive
            }
        }) {
            Image(systemName: "delete.left.fill")
                .font(.system(size: 30))
                .foregroundColor(.gray)
                .opacity(0.7)
                .padding(5)
        }
    }

    enum DrawingConstraints {
        static let MIN_HEIGHT: CGFloat = 100
        static let SWIPE_TO_DISCARD_THRESHOLD: CGFloat = 200
    }
}

struct AudioWaveVisualizer: View {
    var power: Float
    var isPlaying: Bool

    var onTap: (() -> Void)?

    var body: some View {
        AudioWaveView(amplitude: Double(power), isPlaying: isPlaying)
    }
}

struct AudioPowerVisualizer: View {
    var powers: [Float]
    var onTap: (() -> Void)?

    var body: some View {
        GeometryReader { geom in
            HStack(spacing: 3.5) {
                ForEach(Array(powers.enumerated()), id: \.element) { _, power in
                    Capsule().foregroundColor(.red)
                        .frame(width: 1.4, height: CGFloat(power) * geom.size.height)
                }
            }
            .frame(width: geom.size.width, height: geom.size.height, alignment: .trailing)
        }
        .padding(.vertical, 10)
        .contentShape(Rectangle())
        .onTapGesture { onTap?() }
    }
}

struct ComposerView_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
    }

    struct TestView: View {
        @State var isTextFieldFocused: Bool = false

        var body: some View {
            ZStack(alignment: .bottom) {
                List(1 ..< 20) { i in
                    Text("\(i)")
                }

                ComposerView(
                    conversationUuid: UUID(),
                    isTextFieldFocused: $isTextFieldFocused,
                    recordState: .inactive, audioPowers: [0.2, 1.0, 0.5]
                )
            }
        }
    }
}
