//
//  format.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import Foundation

func distance(_ dist: Float) -> String {
    if dist < 2000 {
        return "\(Int(dist.rounded()))m"
    } else {
        return "\(Int((dist / 1000).rounded()))km"
    }
}
