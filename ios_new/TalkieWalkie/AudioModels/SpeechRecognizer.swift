/*
 See LICENSE folder for this sample’s licensing information.
 */

import AVFoundation
import Foundation
import Speech
import SwiftUI

/// A helper for transcribing speech to text using AVAudioEngine.
struct SpeechRecognizer {
    private class SpeechAssist {
        var audioEngine: AVAudioEngine?
        var recognitionRequest: SFSpeechAudioBufferRecognitionRequest?
        var recognitionTask: SFSpeechRecognitionTask?
        let speechRecognizer = SFSpeechRecognizer(locale: .init(identifier: "fr_FR"))

        deinit {
            reset()
        }

        func reset() {
            recognitionTask?.cancel()
            audioEngine?.stop()
            audioEngine = nil
            recognitionRequest = nil
            recognitionTask = nil
        }
    }

    private let assistant = SpeechAssist()

    var onTranscriptReady: ((Transcript) -> Void)?

    /**
        Begin transcribing audio.

        Creates a `SFSpeechRecognitionTask` that transcribes speech to text until you call `stopRecording()`.
        The resulting transcription is continuously written to the provided text binding.

        -  Parameters:
            - speech: A binding to a string where the transcription is written.
     */
    func record() {
        canAccess { authorized in
            guard authorized else {
                return
            }

            assistant.audioEngine = AVAudioEngine()
            guard let audioEngine = assistant.audioEngine else {
                fatalError("Unable to create audio engine")
            }
            assistant.recognitionRequest = SFSpeechAudioBufferRecognitionRequest()
            guard let recognitionRequest = assistant.recognitionRequest else {
                fatalError("Unable to create request")
            }
            recognitionRequest.shouldReportPartialResults = false

            do {
                let audioSession = AVAudioSession.sharedInstance()
                try audioSession.setCategory(.record, mode: .measurement, options: .duckOthers)
                try audioSession.setActive(true, options: .notifyOthersOnDeactivation)
                let inputNode = audioEngine.inputNode

                let recordingFormat = inputNode.outputFormat(forBus: 0)
                inputNode.installTap(onBus: 0, bufferSize: 1024, format: recordingFormat) { (buffer: AVAudioPCMBuffer, _: AVAudioTime) in
                    recognitionRequest.append(buffer)
                }

                audioEngine.prepare()
                try audioEngine.start()
                assistant.recognitionTask = assistant.speechRecognizer?.recognitionTask(with: recognitionRequest) { result, _ in
                    if let result = result {
                        let transcript = Transcript.fromTranscription(transcript: result.bestTranscription)
                        onTranscriptReady?(transcript)
                    }

                    audioEngine.stop()
                    inputNode.removeTap(onBus: 0)
                    self.assistant.recognitionRequest = nil
                }
            } catch {
                print("Error transcibing audio: " + error.localizedDescription)
                assistant.reset()
            }
        }
    }

    /// Stop transcribing audio.
    func stopRecording() {
        assistant.reset()
    }

    private func canAccess(withHandler handler: @escaping (Bool) -> Void) {
        SFSpeechRecognizer.requestAuthorization { status in
            if status == .authorized {
                AVAudioSession.sharedInstance().requestRecordPermission { authorized in
                    handler(authorized)
                }
            } else {
                handler(false)
            }
        }
    }
}
