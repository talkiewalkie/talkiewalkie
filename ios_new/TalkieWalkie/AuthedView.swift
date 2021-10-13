//
//  AuthedView.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 12/10/2021.
//

import SwiftUI

struct AuthedView: View {
    @EnvironmentObject var tooltipManager: TooltipManager
    @Namespace var namespace
    @AppStorage("onboardGuideShown") var onboardGuideShown: Bool = false

    @State var guideState = false
    @State var isRecording = false

    func showGuide() {
        if onboardGuideShown { return }

        guideState.toggle()

        DispatchQueue.main.asyncAfter(deadline: .now() + 0.2) {
            withAnimation(.easeIn) {
                tooltipManager.isPresented = true
            }
        }
    }

    var body: some View {
        ZStack {
            DiscussionListView(namespace: namespace)

            VStack {
                Spacer()

                RecordButton(isRecording: $isRecording)
                    .tooltip(selectionState: guideState,
                             options: .init(orientation: .bottom,
                                            padding: 0,
                                            floating: true), content: {
                                 Text("Record a first voice message!")
                             }, onDismiss: {
                                 onboardGuideShown = true
                             })
                    .padding()
            }
        }
        .onAppear {
            showGuide()
        }
    }
}

// struct AuthedView_Previews: PreviewProvider {
//    static var previews: some View {
//        AuthedView()
//    }
// }
