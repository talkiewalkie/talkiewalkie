//
//  HeaderView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct HeaderSettingsView: View {
    @State var showSettings: Bool = false

    @EnvironmentObject var userStore: UserStore
    @Environment(\.colorScheme) var colorScheme

    var body: some View {
        Button(action: { showSettings = true }) {
            Image(systemName: "gearshape")
                .renderingMode(.template)
                .font(Font.system(size: 14, weight: .medium))
                .padding(10).contentShape(Rectangle())
        }
        .buttonStyle(PlainButtonStyle())
        .sheetWithThemeEnvironment(colorScheme: colorScheme, isPresented: $showSettings) {
            SettingsView(show: $showSettings)
                .font(.body)
                .environmentObject(userStore)
        }
    }
}

struct HeaderView_Previews: PreviewProvider {
    static var previews: some View {
        HeaderSettingsView()
    }
}
