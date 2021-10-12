//
//  TooltipManager.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

/// ROOT-VIEW MODIFIER

public extension View {
    func addTooltip() -> some View {
        modifier(
            Tooltip()
        )
    }
}

struct Tooltip: ViewModifier {
    @EnvironmentObject private var manager: TooltipManager

    private var options: TooltipOptions { manager.options }
    @State private var tooltipSize: CGSize = .zero

    @State var floatingOffset: CGSize = .zero

    func computePosition(frame: CGRect) -> (TooltipShape.Alignment, CGPoint) {
        let center = CGPoint(x: manager.baseFrame.midX, y: manager.baseFrame.midY)

        let y: CGFloat
        if manager.options.orientation == .top {
            y = center.y + manager.baseFrame.height / 2 + tooltipSize.height / 2 + options.offset
        } else {
            y = center.y - (manager.baseFrame.height / 2 + tooltipSize.height / 2 + options.offset)
        }

        if center.x < tooltipSize.width / 2 {
            return (.start, CGPoint(x: center.x + tooltipSize.width / 2 - 24, y: y))
        } else if center.x > frame.maxX - tooltipSize.width / 2 {
            return (.end, CGPoint(x: center.x - tooltipSize.width / 2 + 24, y: y))
        }
        return (.middle, CGPoint(x: center.x, y: y))
    }

    func computeFloatingOffset() -> CGSize {
        switch options.orientation {
        case .top:
            return CGSize(width: 0, height: options.floatingOffset)
        case .bottom:
            return CGSize(width: 0, height: -options.floatingOffset)
        case .left:
            return CGSize(width: options.floatingOffset, height: 0)
        case .right:
            return CGSize(width: -options.floatingOffset, height: 0)
        }
    }

    func elongate(rect: CGRect, frame: CGRect) -> CGRect {
        return CGRect(origin: CGPoint(x: frame.origin.x, y: rect.origin.y - options.touchOffset),
                      size: CGSize(width: frame.width, height: rect.height + 2 * options.touchOffset))
    }

    func contentShape(frame: CGRect) -> CGRect {
        if options.elongated {
            return elongate(rect: manager.baseFrame, frame: frame)
        }
        return manager.baseFrame
    }

    func body(content: Content) -> some View {
        GeometryReader { geom in
            ZStack {
                content

                Group {
                    if manager.isPresented {
                        let (alignment, position) = computePosition(frame: geom.frame(in: .global))

                        Group {
                            if manager.options.touchableHole {
                                Color.black.opacity(0.4)
                                    .contentShape(
                                        HoledShape(rect: contentShape(frame: geom.frame(in: .global)), shape: manager.options.shape),
                                        eoFill: true
                                    )
                            } else {
                                Color.black.opacity(0.4)
                                    .maskWithHole(rect: contentShape(frame: geom.frame(in: .global)), shape: manager.options.shape)
                            }

                            TooltipView(content: manager.content, options: manager.options, alignment: alignment)
                                .background(sizeReader($tooltipSize))
                                .position(position)
                                .offset(floatingOffset)
                        }
                        .onAppear {
                            if options.floating {
                                floatingOffset = .zero
                                withAnimation(Animation.easeInOut(duration: options.floatingSpeed).repeatForever()) {
                                    floatingOffset = computeFloatingOffset()
                                }
                            }
                        }
                    }
                }
                .ignoresSafeArea(.all)
                .simultaneousGesture(TapGesture().onEnded {
                    manager.isPresented = false
                    manager.onDismiss?()
                })
            }
        }
    }

    struct TooltipView<Content: View>: View {
        var content: Content
        var options: TooltipOptions
        var alignment: TooltipShape.Alignment

        var body: some View {
            content
                .foregroundColor(.primary)
                .padding()
                .background(
                    TooltipShape(orientation: options.orientation, alignment: alignment)
                        .stroke(style: StrokeStyle(lineWidth: 8, lineCap: .round, lineJoin: .round)).fill(options.color)
                        .overlay(TooltipShape(orientation: options.orientation, alignment: alignment).fill(options.color))
                )
        }
    }
}

/// VIEW MODIFIER

extension View {
    func tooltip<Content: View>(selectionState: Bool, options: TooltipOptions = TooltipOptions(), @ViewBuilder content: @escaping () -> Content, onDismiss: (() -> Void)? = nil) -> some View {
        TooltipAddView(selectionState: selectionState, options: options, content: content, onDismiss: onDismiss, base: self)
    }
}

public struct TooltipOptions {
    var elongated = false
    var touchableHole = false
    var orientation: TooltipShape.Orientation = .top
    var shape = AnyShape(Capsule())
    var padding: CGFloat = .zero

    var color = Color(#colorLiteral(red: 0.8786171876, green: 0.8786171876, blue: 0.8786171876, alpha: 1))

    var offset: CGFloat = 10
    var touchOffset: CGFloat = 15

    var floating: Bool = false
    var floatingOffset: CGFloat = 10
    var floatingSpeed: Double = 0.5
}

struct TooltipAddView<Base: View, InnerContent: View>: View {
    @EnvironmentObject var manager: TooltipManager

    var selectionState: Bool
    var options: TooltipOptions
    let content: () -> InnerContent
    let onDismiss: (() -> Void)?
    let base: Base

    @State var baseFrame: CGRect = .zero

    var body: some View {
        base
            .padding(options.padding)
            .background(rectReader($baseFrame))
            .onChange(of: selectionState) { _ in
                DispatchQueue.main.async(execute: updateContent)
            }
    }

    func updateContent() {
        manager.update(baseFrame: baseFrame, content: content, options: options, onDismiss: onDismiss)
    }
}

/// MANAGER

public class TooltipManager: ObservableObject {
    /// Published var to present or hide the tooltip
    @Published var isPresented: Bool = false {
        didSet {
//            if !isPresented {
//                DispatchQueue.main.async { [weak self] in
//                    self?.baseFrame = .zero
//                    self?.content = AnyView(EmptyView())
//                    self?.onDismiss = nil
//                }
//            }
        }
    }

    @Published var baseFrame: CGRect = .zero

    /// The content of the tooltip
    @Published private(set) var content: AnyView

    @Published var options = TooltipOptions()

    /// the onDismiss code that runs when the tooltip is closed
    private(set) var onDismiss: (() -> Void)?

    public init() {
        content = AnyView(EmptyView())
    }

    /**
      Updates some properties of the **Tooltip**
     - parameter isPresented: If the tooltip is presented
     - parameter content: The content to place inside of the Tooltip.
     - parameter onDismiss: This code will be runned when the tooltip is dismissed.
     */
    public func update<T>(isPresented: Bool? = nil, baseFrame: CGRect? = nil, content: (() -> T)? = nil, options: TooltipOptions, onDismiss: (() -> Void)? = nil) where T: View {
        if let content = content {
            self.content = AnyView(content())
        }
        self.options = options

        if let baseFrame = baseFrame {
            self.baseFrame = baseFrame
        }
        if let onDismiss = onDismiss {
            self.onDismiss = onDismiss
        }
        if let isPresented = isPresented {
            self.isPresented = isPresented
        }
    }
}

struct TooltipManager_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
            .environmentObject(TooltipManager())
    }

    struct TestView: View {
        @State var selectionState = false

        @EnvironmentObject var tooltipManager: TooltipManager

        init() {
            UINavigationBar.appearance().barTintColor = .systemBackground
            UINavigationBar.appearance().shadowImage = UIImage()
        }

        var body: some View {
            NavigationView {
                ZStack {
                    Color.white
                        .ignoresSafeArea(.all)

                    VStack {
                        Text("\(tooltipManager.baseFrame.minX)")
                        Text("\(tooltipManager.baseFrame.minY)")
                        Text("\(tooltipManager.baseFrame.width)")
                        Text("\(tooltipManager.baseFrame.height)")

                        Spacer()

                        HStack {
                            Spacer()

                            Button(action: {
                                withAnimation(.easeIn) {
                                    tooltipManager.isPresented = true
                                }
                            }, label: {
                                Text("Lorem Ipsum")
                                    .padding()
                                    .background(Capsule().fill(Color.white))
                            })
                                .tooltip(selectionState: selectionState, options: TooltipOptions(floating: true), content: {
                                    Text("Lorem Ipsum")
                                })

                            Spacer()
                        }

                        Spacer()
                    }
                }

                .navigationTitle("hello")
                .navigationBarTitleDisplayMode(.inline)
            }
            .addTooltip()
            .onAppear {
                selectionState.toggle()
            }
        }
    }
}
