//
//  HoledMask.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

private func holedMask<ShapeType: Shape>(in frame: CGRect, rect: CGRect, shape: ShapeType) -> Path {
    var path = Rectangle().path(in: frame)
    path.addPath(shape.path(in: rect))
    return path
}

struct HoledShape<ShapeType: Shape>: Shape {
    let rect: CGRect
    let shape: ShapeType

    func path(in frame: CGRect) -> Path {
        holedMask(in: frame, rect: rect, shape: shape)
    }
}

struct HoledShapeRatio<ShapeType: Shape>: Shape {
    let shape: ShapeType
    let scale: CGSize

    func path(in frame: CGRect) -> Path {
        let size = CGSize(width: frame.width * scale.width, height: frame.height * scale.height)
        let frameRect = CGRect(origin: frame.origin, size: size)

        return holedMask(in: frame, rect: frameRect, shape: shape)
    }
}

struct HoledMask<ShapeType: Shape>: ViewModifier {
    let rect: CGRect
    let shape: ShapeType

    func body(content: Content) -> some View {
        GeometryReader { geom in
            content
                .mask(
                    holedMask(in: geom.frame(in: .global), rect: rect, shape: shape)
                        .fill(style: FillStyle(eoFill: true))
                )
        }
    }
}

public extension View {
    func maskWithHole<ShapeType: Shape>(rect: CGRect, shape: ShapeType) -> some View {
        modifier(
            HoledMask(rect: rect, shape: shape)
        )
    }
}

struct HoledMask_Previews: PreviewProvider {
    static let rect = CGRect(origin: CGPoint(x: 100, y: 100), size: CGSize(width: 200, height: 100))

    static var previews: some View {
        Color.gray
            .maskWithHole(rect: rect, shape: Circle())
            .ignoresSafeArea(.all)
    }
}
