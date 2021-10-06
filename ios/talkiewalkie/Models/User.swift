//
//  User.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 04/10/2021.
//

import CoreData
import Foundation

import FirebaseAuth

extension User {
    static func upsert(_ u: App_User, context: NSManagedObjectContext) -> User {
        let localU = User.getByUuidOrCreate(u.uuid.uuidOrThrow(), context: context)
        localU.uuid = u.uuid.uuidOrThrow()
        localU.handle = u.handle
        return localU
    }
}

extension Me {
    static func fromCache(context: NSManagedObjectContext) -> Me? {
        let req = NSFetchRequest<Me>(entityName: "Me")
        req.predicate = .all

        let res = (try? context.fetch(req)) ?? []
        return res.first
    }
}
