//
//  AnyShape.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

public struct AnyShape: Shape {
    public var make: (CGRect, inout Path) -> Void

    public init(_ make: @escaping (CGRect, inout Path) -> Void) {
        self.make = make
    }

    public init<S: Shape>(_ shape: S) {
        make = { rect, path in
            path = shape.path(in: rect)
        }
    }

    public func path(in rect: CGRect) -> Path {
        return Path { [make] in make(rect, &$0) }
    }
}
