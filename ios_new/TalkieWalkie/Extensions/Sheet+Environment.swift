//
//  Sheet+Environment.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

extension View {
    func sheetWithThemeEnvironment<Content>(colorScheme: ColorScheme,
                                            isPresented: Binding<Bool>, onDismiss: (() -> Void)? = nil, content: @escaping () -> Content) -> some View where Content : View {
        return self.sheet(isPresented: isPresented, onDismiss: onDismiss, content: {
            content()
                .preferredColorScheme(colorScheme)
                .environment(\.colorScheme, colorScheme)
        })
    }
}


extension View {
    func partialSheetWithThemeEnvironment<Content>(colorScheme: ColorScheme,
                                            isPresented: Binding<Bool>, content: @escaping () -> Content) -> some View where Content : View {
        return self.partialSheet(isPresented: isPresented, content: {
            content()
                .preferredColorScheme(colorScheme)
                .environment(\.colorScheme, colorScheme)
        })
    }
}
