//
//  User.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 04/10/2021.
//

import CoreData
import Foundation
import OSLog

import FirebaseAuth

extension User {
    @discardableResult
    static func fromProto(_ u: App_User, context: NSManagedObjectContext, block: (_ user: User) -> Void = { _ in }) -> User {
        let localU = User.getByUuidOrCreate(u.uuid.uuidOrThrow(), context: context)
        localU.uuid = u.uuid.uuidOrThrow()
        localU.displayName = u.displayName
        localU.phone = u.phone
    
        block(localU)
        
        return localU
    }
}

extension Me {
    static func fromCache(context: NSManagedObjectContext) -> Me? {
        let req = NSFetchRequest<Me>(entityName: "Me")
        req.predicate = .all

        let res = (try? context.fetch(req)) ?? []

        if res.count > 1 {
            os_log(.error, "getting \(res.count) instances of [Me] object in core data in ctx(\(context.description))")
            // I hate Core Data. Why would we have multiple instance of the Me object???
            context.perform {
                res.filter { $0.uuid == nil }.forEach { context.delete($0) }
                context.saveOrLogError()
            }
        } else if res.isEmpty {
            os_log(.debug, "no [Me] object found in core data in ctx(\(context.description))")
        }
        

        return res.first { $0.uuid != nil }
    }
}
