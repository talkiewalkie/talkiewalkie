//
//  Extensions.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 05/10/2021.
//

import CoreData
import Foundation

extension NSManagedObject {
    static func getByUuidOrCreate(_ uuid: UUID, context: NSManagedObjectContext) -> Self {
        let localUsersRq = NSFetchRequest<Self>(entityName: Self.entity().name!)
        localUsersRq.predicate = NSPredicate(format: "uuid = %@", uuid.uuidString)
        let localUsers = (try? context.fetch(localUsersRq)) ?? []

        if let me = localUsers.first {
            return me
        } else {
            let new = Self(context: context)
            try? context.save()

            return new
        }
    }
}
