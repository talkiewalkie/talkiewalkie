//
//  User.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 27/09/2021.
//

import Foundation

struct User: Codable, Identifiable {
    let uuid: String
    var id: String { uuid }
    
    let handle: String
    let createdAt: Date
}
