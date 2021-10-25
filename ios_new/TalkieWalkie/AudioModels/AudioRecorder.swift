//
//  AudioRecorder.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 25.10.21.
//

import Foundation
import AVFoundation
import CoreGraphics


enum AudioRecorderState {
    case stopped
    case playing
    case recording
}


class AudioRecorder: NSObject {
    
    var player: AVAudioPlayer?
    var recorder: AVAudioRecorder!
    
    private var audioURL: URL?
    
    private(set) var state: AudioRecorderState = .stopped
    
    override init() {
        super.init()
        
        setupRecorder()
        setupAudioSession()
        enableBuiltInMic()
    }
    
    // MARK: - Audio Session Configuration
    func setupAudioSession() {
        do {
            let session = AVAudioSession.sharedInstance()
            try session.setCategory(.playAndRecord, options: [.defaultToSpeaker, .allowBluetooth])
            try session.setActive(true)
        } catch {
            fatalError("Failed to configure and activate session.")
        }
    }
    
    private func enableBuiltInMic() {
        // Get the shared audio session.
        let session = AVAudioSession.sharedInstance()
        
        // Find the built-in microphone input.
        guard let availableInputs = session.availableInputs,
              let builtInMicInput = availableInputs.first(where: { $0.portType == .builtInMic }) else {
            print("The device must have a built-in microphone.")
            return
        }
        
        // Make the built-in microphone input the preferred input.
        do {
            try session.setPreferredInput(builtInMicInput)
        } catch {
            print("Unable to set the built-in mic as the preferred input.")
        }
    }
    
    // MARK: - Audio Recording and Playback
    func setupRecorder() {
        let tempDir = URL(fileURLWithPath: NSTemporaryDirectory())
        let fileURL = tempDir.appendingPathComponent("recording.wav")

        do {
//            let settings: [String: Any] = [
//                AVFormatIDKey: kAudioFormatOpus,
//                //AVEncoderAudioQualityKey: AVAudioQuality.high.rawValue,
//                AVNumberOfChannelsKey: 1,
//                AVSampleRateKey: 24000.0,
//            ]
            let settings: [String: Any] = [
                AVFormatIDKey: Int(kAudioFormatLinearPCM),
                AVLinearPCMIsNonInterleaved: false,
                AVSampleRateKey: 44_100.0,
                AVNumberOfChannelsKey: 1,
                AVLinearPCMBitDepthKey: 16
            ]
            recorder = try AVAudioRecorder(url: fileURL, settings: settings)
        } catch {
            fatalError("Unable to create audio recorder: \(error.localizedDescription)")
        }
        
        recorder.delegate = self
        recorder.isMeteringEnabled = true
        recorder.prepareToRecord()
    }
    
    @discardableResult
    func record() -> Bool {
        let started = recorder.record()
        state = .recording
        return started
    }
    
    // Stops recording and calls the completion callback when the recording finishes.
    func stopRecording() {
        recorder.stop()
        
        state = .stopped
    }
    
    func play() {
        guard let url = audioURL else { print("No recording to play"); return }
        player = try? AVAudioPlayer(contentsOf: url)
        player?.isMeteringEnabled = true
        player?.delegate = self
        player?.play()
        state = .playing
    }
    
    func stopPlayback() {
        player?.stop()
        state = .stopped
    }
    
    var recordingURL: URL? {
        let directory = FileManager.default.urls(for: .applicationSupportDirectory, in: .userDomainMask).first!
        let url = directory.appendingPathComponent("recording.wav")
        return FileManager.default.fileExists(atPath: url.path) ? url : nil
    }
    
}

extension AudioRecorder: AVAudioRecorderDelegate {
    
    // AVAudioRecorderDelegate method.
    func audioRecorderDidFinishRecording(_ recorder: AVAudioRecorder, successfully flag: Bool) {
        
        let directory = FileManager.default.urls(for: .applicationSupportDirectory, in: .userDomainMask).first!
        let destURL = directory.appendingPathComponent("recording.wav")
        try? FileManager.default.removeItem(at: destURL)
        try? FileManager.default.copyItem(at: recorder.url, to: destURL)
        recorder.prepareToRecord()
        
        audioURL = destURL
    }
    
    
}


extension AudioRecorder: AVAudioPlayerDelegate {
    func audioPlayerDidFinishPlaying(_ player: AVAudioPlayer, successfully flag: Bool) {
        
    }
}
