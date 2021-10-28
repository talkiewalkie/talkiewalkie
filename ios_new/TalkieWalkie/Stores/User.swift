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
    @discardableResult
    static func upsert(_ u: App_User, context: NSManagedObjectContext) -> User {
        let localU = User(context: context)
        localU.uuid = u.uuid.uuidOrThrow()
        localU.displayName = u.displayName
        localU.phone = u.phone
        return localU
    }
}

extension Me {
    static func fromCache(context: NSManagedObjectContext) -> Me? {
        let req = NSFetchRequest<Me>(entityName: "Me")
        req.predicate = .all

        let res = (try? context.fetch(req)) ?? []
        
        // I hate Core Data. Why would we have multiple instance of the Me object???
        context.perform {
            res.filter { $0.uuid == nil }.forEach { context.delete($0) }
            context.saveOrLogError()
        }
        
        return res.first { $0.uuid != nil }
    }
}
