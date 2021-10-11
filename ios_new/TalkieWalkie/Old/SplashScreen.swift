//
//  SplashScreen.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct SplashScreen: View {
    var body: some View {
        ZStack {
            LinearGradient(colors: [Color("Purple"), Color("Red"), Color("Yellow")], startPoint: .init(x: -0.2, y: -0.4), endPoint: .bottomTrailing)
            
            
            Image("logo_bubble")
                .resizable()
                .aspectRatio(contentMode: .fit)
                .scaleEffect(0.4)
                .opacity(0.975)
        }
        .ignoresSafeArea()
    }
}

struct SplashScreen_Previews: PreviewProvider {
    static var previews: some View {
        SplashScreen()
    }
}
