//
//  UserDefaults.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 23/11/2021.
//

import Foundation

extension UserDefaults {
    enum Keys: String, CaseIterable {
        // Onboarding vars
        case handle
        case displayName
        case phoneCountryCode
        case phoneRegionID
        case phoneNumber
        case showOnboarding
        case hasRefusedSharingContacts

        // State mgmgt
        case lastEventUuid
    }

    func reset() {
        Keys.allCases.forEach { removeObject(forKey: $0.rawValue) }
    }
}
