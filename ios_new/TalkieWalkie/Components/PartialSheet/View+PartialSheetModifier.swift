//
//  View+PartialSheetModifier.swift
//  PartialModal
//
//  Created by Miotto Andrea on 10/11/2019.
//  Copyright © 2019 Miotto Andrea. All rights reserved.
//

import SwiftUI

public extension View {
    /**
     Add a PartialSheet to the current view. You should attach it to your Root View.
     Use the PartialSheetManager as an environment object to present it whenever you want.
     - parameter style: The style configuration for the Partial Sheet.
     */
    func addPartialSheet(
        style: PartialSheetStyle = PartialSheetStyle.defaultStyle()) -> some View
    {
        modifier(
            PartialSheet(
                style: style
            )
        )
    }
}
