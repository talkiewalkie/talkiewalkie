//
//  Dictionary+Extension.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import Foundation

func +<Key, Value> (lhs: [Key: Value], rhs: [Key: Value]) -> [Key: Value] {
    var result = lhs
    rhs.forEach{ result[$0] = $1 }
    return result
}
