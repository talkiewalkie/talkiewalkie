//
//  DummyStore.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 27.10.21.
//

import SwiftUI

struct DummyStore {
    static let shared = DummyStore()

    @AppStorage("isDarkMode") var isDarkMode: Bool = false

    var auth = AuthState()
    var messageViewModel = MessageViewModel()
    var onboardingViewModel = OnboardingViewModel(name: .constant(""), handle: .constant(""), phoneCountryCode: 33, phoneRegionID: "FR", phoneNumber: "")

    var tooltipManager = TooltipManager()
    var partialSheetManager = PartialSheetManager()

    init() {}
}

extension View {
    func withDummmyEnvironments() -> some View {
        environmentObject(DummyStore.shared.auth)
            .environmentObject(DummyStore.shared.messageViewModel)
            .environmentObject(DummyStore.shared.onboardingViewModel)

//            .environment(\.managedObjectContext, DummyStore.shared.auth.moc)

            .preferredColorScheme(DummyStore.shared.isDarkMode ? .dark : .light)
            .environment(\.colorScheme, DummyStore.shared.isDarkMode ? .dark : .light)

            .addTooltip()
            .environmentObject(DummyStore.shared.tooltipManager)
            .addPartialSheet()
            .environmentObject(DummyStore.shared.partialSheetManager)
    }
}
