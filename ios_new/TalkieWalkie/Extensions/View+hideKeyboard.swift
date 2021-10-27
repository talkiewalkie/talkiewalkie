//
//  View+hideKeyboard.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 26.10.21.
//

import SwiftUI

extension View {
    func hideKeyboard() {
        UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
    }
}
