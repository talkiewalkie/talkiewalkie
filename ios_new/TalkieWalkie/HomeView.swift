//
//  HomeView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

class MessageViewModel: ObservableObject {
    @Published var message: ChatMessage?
    @Published var showDetailView: Bool = false
}

struct HomeView: View {
    @AppStorage("showOnboarding") var showOnboarding: Bool = true

    @State var isRecording = false

    @EnvironmentObject var tooltipManager: TooltipManager
    @State var guideState = false
    @AppStorage("onboardGuideShown") var onboardGuideShown: Bool = false

    @StateObject var messageViewModel = MessageViewModel()
    @Namespace var namespace

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
        Group {
            if showOnboarding {
                OnboardingView(onboardingDone: onboardingDone)
            } else {
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

                    if messageViewModel.showDetailView {
                        MessageDetailView(namespace: namespace)
                    }
                }
                .onAppear {
                    showGuide()
                }
            }
        }
        .environmentObject(messageViewModel)
    }

    func onboardingDone() {
        showOnboarding = false
    }
}

struct HomeView_Previews: PreviewProvider {
    static var previews: some View {
        HomeView(showOnboarding: false, onboardGuideShown: true)
            .withDummyVariables()
    }
}
