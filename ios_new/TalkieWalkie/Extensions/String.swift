//
//  String.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 12/10/2021.
//

import Foundation

extension String {
    func uuidOrThrow() -> UUID {
        return UUID(uuidString: self)!
    }
}
