//
//  ErrorPopupManager.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct ErrorPopupMessage: Identifiable {
    let id = UUID()
    var message: String
}

public extension View {
    func addErrorPopup() -> some View {
        modifier(
            ErrorPopupModifier()
        )
    }
}

struct ErrorPopupModifier: ViewModifier {
    @EnvironmentObject private var manager: ErrorPopupManager

    func body(content: Content) -> some View {
        GeometryReader { _ in
            ZStack {
                content

                ZStack {
                    ForEach(manager.errorMessages) {
                        ErrorPopupView(errorMessage: $0)
                            .transition(.move(edge: .top).combined(with: .opacity))
                    }
                }
                .padding(.horizontal, 10)
                .frame(maxHeight: .infinity, alignment: .top)
            }
        }
    }
}

struct ErrorPopupView: View {
    var errorMessage: ErrorPopupMessage

    var body: some View {
        Text(errorMessage.message)
            .foregroundColor(.white)
            .padding(.vertical, 10).padding(.horizontal)
            .frame(maxWidth: .infinity)
            .multilineTextAlignment(.center)
            .background(Color.red)
            .opacity(0.9)
    }
}

class ErrorPopupManager: ObservableObject {
    @Published fileprivate var errorMessages: [ErrorPopupMessage] = []

    func showError(_ message: String, displayDuration: Double = 2) {
        let errorMessage = ErrorPopupMessage(message: message)

        withAnimation(.easeInOut) {
            self.errorMessages.append(errorMessage)
        }
        withAnimation(.easeInOut.delay(0.1)) {
            self.errorMessages.removeFirst(self.errorMessages.count - 1)
        }

        DispatchQueue.main.asyncAfter(deadline: .now() + displayDuration) {
            if let index = self.errorMessages.firstIndex(where: { $0.id == errorMessage.id }) {
                _ = withAnimation(.easeInOut) {
                    self.errorMessages.remove(at: index)
                }
            }
        }
    }
}

struct ErrorPopupManager_Previews: PreviewProvider {
    static var previews: some View {
        VStack {
            TestView()
                .addErrorPopup()
                .environmentObject(ErrorPopupManager())
        }
    }

    struct TestView: View {
        @EnvironmentObject var errorPopupManager: ErrorPopupManager

        var body: some View {
            VStack {
                Button(action: {
                    errorPopupManager.showError("Error message")
                }) {
                    Text("Click")
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
    }
}
