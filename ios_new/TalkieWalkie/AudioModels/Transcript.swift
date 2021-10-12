//
//  Transcript.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 07.10.21.
//

import Foundation
import Speech

struct Transcript: Identifiable, Equatable {
    static func == (lhs: Transcript, rhs: Transcript) -> Bool {
        return lhs.string == rhs.string && lhs.segments == rhs.segments
    }

    let id = UUID()

    let string: String
    let segments: [Segment]

    var totalDuration: Double {
        return segments.map { Double($0.timestamp + $0.timestamp) }.max() ?? .zero
    }

    struct Segment: Identifiable, Equatable {
        let id = UUID()

        let string: String
        let timestamp: TimeInterval
        let duration: TimeInterval
        let substringRange: NSRange
    }

    static func fromTranscription(transcript: SFTranscription) -> Transcript {
        let segments = transcript.segments.map { segment in
            Segment(string: segment.substring,
                    timestamp: segment.timestamp,
                    duration: segment.duration,
                    substringRange: segment.substringRange)
        }

        return Transcript(string: transcript.formattedString, segments: segments)
    }
}

let dummyTranscript = Transcript(string: "Hello comment ça va",
                                 segments: [
                                     Transcript.Segment(string: "Hello", timestamp: 0.3, duration: 0.5, substringRange: .init(location: 0, length: 5)),
                                     Transcript.Segment(string: "comment", timestamp: 1.0, duration: 0.5, substringRange: .init(location: 6, length: 7)),
                                     Transcript.Segment(string: "ça", timestamp: 1.6, duration: 0.2, substringRange: .init(location: 14, length: 2)),
                                     Transcript.Segment(string: "va", timestamp: 1.8, duration: 0.3, substringRange: .init(location: 17, length: 2)),
                                 ])
