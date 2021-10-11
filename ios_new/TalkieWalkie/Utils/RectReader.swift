//
//  RectReader.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI

func rectReader(_ binding: Binding<CGRect>) -> some View {
    return GeometryReader { (geometry) -> AnyView in
        let rect = geometry.frame(in: .global)
        DispatchQueue.main.async {
            binding.wrappedValue = rect
        }
        return AnyView(Color.clear)
    }
}

func sizeReader(_ binding: Binding<CGSize>) -> some View {
    return GeometryReader { (geometry) -> AnyView in
        let size = geometry.size
        DispatchQueue.main.async {
            binding.wrappedValue = size
        }
        return AnyView(Color.clear)
    }
}
