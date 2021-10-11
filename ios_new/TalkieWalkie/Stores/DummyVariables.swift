//
//  DummyVariables.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import Foundation
import SwiftUI


struct DummyVariables {
    static let shared = DummyVariables()
    
    @AppStorage("isDarkMode") var isDarkMode: Bool = false
    
    var dummyUserStore: UserStore

    var dummyTooltipManager = TooltipManager()

    var dummyPartialSheetManager = PartialSheetManager()
    
    var dummyMessageViewModel: MessageViewModel
    
    init() {
        dummyUserStore = UserStore()
        
        dummyMessageViewModel = MessageViewModel()
        dummyMessageViewModel.message = dummyChatMessages[0]
    }
}


extension View {
    public func withDummyVariables() -> some View {
        return environmentObject(DummyVariables.shared.dummyUserStore)

            .preferredColorScheme(DummyVariables.shared.isDarkMode ? .dark : .light)
            .environment(\.colorScheme, DummyVariables.shared.isDarkMode ? .dark : .light)
            
            .addTooltip()
            .environmentObject(DummyVariables.shared.dummyTooltipManager)
        
            .addPartialSheet()
            .environmentObject(DummyVariables.shared.dummyPartialSheetManager)
        
            .environmentObject(DummyVariables.shared.dummyMessageViewModel)
    }
}
