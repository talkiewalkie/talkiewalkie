//
//  OSLog.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 12/10/2021.
//

import Foundation
import OSLog

extension Logger {
    static func withLabel(_ label: String) -> Logger {
        return Logger(subsystem: Bundle.main.bundleIdentifier!, category: label)
    }
}
