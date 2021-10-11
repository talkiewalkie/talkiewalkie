//
//  TooltipShape.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

struct TooltipShape: Shape {
    var orientation: Orientation
    var alignment: Alignment = .middle
    var tipSize: CGFloat = 8
    var borderOffset: CGFloat = 8
    
    func path(in rect: CGRect) -> Path {
        
        let topLeft = CGPoint(x: rect.minX + tipSize, y: rect.minY + tipSize)
        let topRight = CGPoint(x: rect.maxX - tipSize, y: rect.minY + tipSize)
        let bottomRight = CGPoint(x: rect.maxX - tipSize, y: rect.maxY - tipSize)
        let bottomLeft = CGPoint(x: rect.minX + tipSize, y: rect.maxY - tipSize)
        
        
        let tooltipMiddle: CGPoint
        switch alignment {
        case .start:
            switch orientation.direction {
            case .horizontal:
                tooltipMiddle = .init(x: rect.minX + 2 * tipSize + borderOffset, y: orientation == .top ? rect.minY : rect.maxY)
            case .vertical:
                tooltipMiddle = .init(x: orientation == .left ? rect.minX : rect.maxX, y: rect.minY + 2 * tipSize + borderOffset)
            }
        case .middle:
            switch orientation.direction {
            case .horizontal:
                tooltipMiddle = .init(x: rect.midX, y: orientation == .top ? rect.minY : rect.maxY)
            case .vertical:
                tooltipMiddle = .init(x: orientation == .left ? rect.minX : rect.maxX, y: rect.midY)
            }
        case .end:
            switch orientation.direction {
            case .horizontal:
                tooltipMiddle = .init(x: rect.maxX - 2 * tipSize - borderOffset, y: orientation == .top ? rect.minY : rect.maxY)
            case .vertical:
                tooltipMiddle = .init(x: orientation == .left ? rect.minX : rect.maxX, y: rect.maxY - 2 * tipSize - borderOffset)
            }
        }
        
        let tooltipStart: CGPoint, tooltipEnd: CGPoint
        switch orientation {
        case .left:
            tooltipStart = .init(x: tooltipMiddle.x + tipSize, y: tooltipMiddle.y + tipSize)
            tooltipEnd = .init(x: tooltipMiddle.x + tipSize, y: tooltipMiddle.y - tipSize)
        case .top:
            tooltipStart = .init(x: tooltipMiddle.x - tipSize, y: tooltipMiddle.y + tipSize)
            tooltipEnd = .init(x: tooltipMiddle.x + tipSize, y: tooltipMiddle.y + tipSize)
        case .right:
            tooltipStart = .init(x: tooltipMiddle.x - tipSize, y: tooltipMiddle.y - tipSize)
            tooltipEnd = .init(x: tooltipMiddle.x - tipSize, y: tooltipMiddle.y + tipSize)
        case .bottom:
            tooltipStart = .init(x: tooltipMiddle.x + tipSize, y: tooltipMiddle.y - tipSize)
            tooltipEnd = .init(x: tooltipMiddle.x - tipSize, y: tooltipMiddle.y - tipSize)
        }
        
        
        return Path { path in
            path.move(to:topLeft)
            for (index, point) in [topRight, bottomRight, bottomLeft, topLeft].enumerated() {
                if index == orientation.rawValue {
                    path.addLine(to: tooltipStart)
                    path.addLine(to: tooltipMiddle)
                    path.addLine(to: tooltipEnd)
                }
                path.addLine(to: point)
            }
            path.closeSubpath()
        }
    }
    
    enum Orientation: Int, CaseIterable {
        case top = 0, right = 1, bottom = 2, left = 3
        
        var direction: Direction {
            switch self {
            case .left, .right:
                return .vertical
            case .top, .bottom:
                return .horizontal
            }
        }
        
        enum Direction {
            case horizontal, vertical
        }
    }
    
    enum Alignment: CaseIterable {
        case start, middle, end
    }
}

struct ToolTipShape_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
    }
    
    struct TestView: View {
        let size: CGFloat = 100
        
        var body: some View {
            VStack {
                LazyVGrid(columns: [GridItem(.fixed(size)), GridItem(.fixed(size)), GridItem(.fixed(size))]) {
                    
                    TooltipShape(orientation: .bottom, alignment: .start)
                        .frame(width: size, height: size)
                    
                    TooltipShape(orientation: .bottom, alignment: .middle)
                        .frame(width: size, height: size)
                    
                    TooltipShape(orientation: .bottom, alignment: .end)
                        .frame(width: size, height: size)
                    
                    TooltipShape(orientation: .top, alignment: .start)
                        .frame(width: size, height: size)
                    
                    TooltipShape(orientation: .top, alignment: .middle)
                        .frame(width: size, height: size)
                    
                    TooltipShape(orientation: .top, alignment: .end)
                        .frame(width: size, height: size)

                }
                
                LazyHGrid(rows: [GridItem(.fixed(size)), GridItem(.fixed(size)), GridItem(.fixed(size))]) {

                    TooltipShape(orientation: .right, alignment: .start)
                        .frame(width: size, height: size)
                    
                    TooltipShape(orientation: .right, alignment: .middle)
                        .frame(width: size, height: size)
                    
                    TooltipShape(orientation: .right, alignment: .end)
                        .frame(width: size, height: size)
                    
                    TooltipShape(orientation: .left, alignment: .start)
                        .frame(width: size, height: size)
                    
                    TooltipShape(orientation: .left, alignment: .middle)
                        .frame(width: size, height: size)
                    
                    TooltipShape(orientation: .left, alignment: .end)
                        .frame(width: size, height: size)
                    
                }
            }
        }
    }
}
