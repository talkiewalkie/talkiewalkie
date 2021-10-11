//
//  UserStore.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import Foundation
import SwiftUI


class UserStore : ObservableObject {
    static var languageCode: String {
            return Locale.current.languageCode ?? "en"
        }
    
    static func openSettings() {
            if let appSettings = URL(string: UIApplication.openSettingsURLString) {
                UIApplication.shared.open(appSettings)
            }
        }
    
}
