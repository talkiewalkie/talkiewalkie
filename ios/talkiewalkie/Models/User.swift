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
    static func byUuid(_ uuid: UUID, context: NSManagedObjectContext, api: Api) -> User {
        let localUsersRq = NSFetchRequest<User>(entityName: "User")
        localUsersRq.predicate = NSPredicate(format: "uuid = %@", uuid.uuidString)
        let localUsers = (try? context.fetch(localUsersRq)) ?? []

        if let me = localUsers.first {
            return me
        } else {
            let u = Me(context: context)
            api.me { res, _ in
                if let remoteU = res {
                    u.uuid = remoteU.uuid
                    u.handle = remoteU.handle
                }
            }

            u.objectWillChange.send()
            try? context.save()

            return u
        }
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
