//
//  AudioRecorderView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 06.10.21.
//

import SwiftUI
import AVFoundation


class AudioRecordViewModel: ObservableObject {
    @Published var transcript: Transcript?
    
    var player: AVPlayer { AVPlayer.sharedDingPlayer }
    
    var speechRecognizer = SpeechRecognizer()
    
    init() {
        speechRecognizer = SpeechRecognizer(onTranscriptReady: onTranscriptReady)
    }
    
    func onTranscriptReady(transcript: Transcript) {
        self.transcript = transcript
    }
    
    func startRecording() {
        // Play sound
        player.seek(to: .zero); player.play()
        
        speechRecognizer.record()
    }
    
    func endRecording() {
        speechRecognizer.stopRecording()
    }
}


struct AudioRecorderView: View {
    @State var isRecording: Bool = false
    
    @StateObject var model = AudioRecordViewModel()
    
    @State var transcripts = [Transcript]()
    
    var body: some View {
        
        VStack {

            ScrollView {
                VStack {
                    ForEach(transcripts) { transcript in
                        AudiogramView(transcript: transcript)
                    }
                }
            }
            .frame(maxHeight: .infinity)
            
            RecordButton(isRecording: $isRecording)
                .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .bottom)
            
        }
        .onChange(of: isRecording) { newValue in
            if newValue { model.startRecording() } else { model.endRecording() }
        }
        .onChange(of: model.transcript) { newValue in
            if let transcript = model.transcript {
                transcripts.append(transcript)
            }   
        }
    }

}

struct AudioRecorderView_Previews: PreviewProvider {
    static var previews: some View {
        AudioRecorderView()
    }
}
