//
//  Bundle+AppInfo.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import Foundation

extension Bundle {
    var appName: String {
        return (infoDictionary?["CFBundleName"] as? String) ?? "Litso"
    }

    var bundleId: String {
        return bundleIdentifier ?? "io.litso.Litso"
    }

    var versionNumber: String {
        return (infoDictionary?["CFBundleShortVersionString"] as? String) ?? "1.0.0"
    }

    var buildNumber: String {
        return (infoDictionary?["CFBundleVersion"] as? String) ?? "1"
    }
}
