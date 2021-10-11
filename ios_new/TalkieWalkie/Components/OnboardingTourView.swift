//
//  OnboardingTourView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI


extension View {
    public func addGuidedTour() -> some View {
        self.modifier(
            GuidedTour()
        )
    }

    func guidedTour<Content: View>(isPresented: Binding<Bool>, @ViewBuilder content: @escaping () -> Content) -> some View {
        GuidedTourAddView(isPresented: isPresented, content: content, base: self)
    }
}

struct GuidedTour: ViewModifier {
    @EnvironmentObject private var manager: GuidedTourManager

    @State var toolTipSize: CGSize = .zero
    var offset: CGFloat = 20

    func tooltipPosition() -> CGPoint {
        return CGPoint(x: manager.baseFrame.midX, y: manager.baseFrame.minY - toolTipSize.height / 2 - offset)
    }

    func body(content: Content) -> some View {
        ZStack {

            content

            Group {
                if manager.isPresented {
                    Color.black.opacity(0.4)
                        .maskWithHole(rect: manager.baseFrame, shape: Capsule())

                ToolTipView(content: manager.content)
                    .background(sizeReader($toolTipSize))
                    .position(tooltipPosition())

                }
            }
            .ignoresSafeArea(.all)
            .onTapGesture {
                manager.isPresented = false
                manager.onDismiss?()
            }
        }
    }

    struct ToolTipView<Content: View>: View {
        var content: Content

        var body: some View {
            content
                .foregroundColor(.primary)
                .padding()
                .background(
                    TooltipShape(orientation: .bottom)
                        .stroke(style: StrokeStyle(lineWidth: 8, lineCap: .round, lineJoin: .round)).fill(Color(#colorLiteral(red: 0.8786171876, green: 0.8786171876, blue: 0.8786171876, alpha: 1)))
                        .overlay(TooltipShape(orientation: .bottom).fill(Color(#colorLiteral(red: 0.8786171876, green: 0.8786171876, blue: 0.8786171876, alpha: 1))))
                )
        }
    }
}

struct GuidedTourAddView<Base: View, InnerContent: View>: View {
    @EnvironmentObject var manager: GuidedTourManager

    @Binding var isPresented: Bool
    let content: () -> InnerContent
    let base: Base

    var body: some View {
        base
            .background(rectReader($manager.baseFrame))
            .onChange(of: isPresented) { _ in
                DispatchQueue.main.async(execute: updateContent)
            }
    }

    func updateContent() {
        manager.updateGuidedTour(isPresented: isPresented, content: content, onDismiss: {
            self.isPresented = false
        })
    }
}

public class GuidedTourManager: ObservableObject {
    /// Published var to present or hide the partial sheet
    @Published var isPresented: Bool = false {
        didSet {
            if !isPresented {
                DispatchQueue.main.async { [weak self] in
                    self?.content = AnyView(EmptyView())
                    self?.onDismiss = nil
                }
            }
        }
    }
    /// The content of the sheet
    @Published private(set) var content: AnyView
    /// the onDismiss code runned when the partial sheet is closed
    private(set) var onDismiss: (() -> Void)?

    @Published var baseFrame: CGRect = .zero

    public init() {
        self.content = AnyView(EmptyView())
    }

    /**
     Updates some properties of the **Partial Sheet**
    - parameter isPresented: If the partial sheet is presented
    - parameter content: The content to place inside of the Partial Sheet.
    - parameter onDismiss: This code will be runned when the sheet is dismissed.
    */
    public func updateGuidedTour<T>(isPresented: Bool? = nil, content: (() -> T)? = nil, onDismiss: (() -> Void)? = nil) where T: View {
        if let content = content {
            self.content = AnyView(content())
        }
        if let onDismiss = onDismiss {
            self.onDismiss = onDismiss
        }
        if let isPresented = isPresented {
            withAnimation {
                self.isPresented = isPresented
            }
        }

    }
}


struct OnboardingTourView_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
    }

    struct TestView: View {
        @State var showGuidedTour = false

        var body: some View {
            ZStack {
                Color.red
                    .ignoresSafeArea(.all)

                VStack {
                    
                    Spacer()

                    Button(action: {
                        showGuidedTour = true
                    }, label: {
                        Text("Lorem Ipsum")
                            .foregroundColor(.primary)
                            .padding()
                            .background(Capsule().fill(Color.white))
                    })

                    Spacer()

                    Text("Lorem Ipsum")
                        .padding()
                        .background(Capsule().fill(Color.white))
                        .guidedTour(isPresented: $showGuidedTour, content: {
                            Text("Lorem Ipsum")
                        })

                    Spacer()

                }
            }
            .addGuidedTour()
            .environmentObject(GuidedTourManager())
        }
    }

}
