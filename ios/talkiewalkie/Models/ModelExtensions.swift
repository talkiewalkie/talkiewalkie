//
//  Extensions.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 05/10/2021.
//

import CoreData
import Foundation
import OSLog

extension NSManagedObject {
    static func getByUuidOrCreate(_ uuid: UUID, context: NSManagedObjectContext) -> Self {
        guard let ename = Self.entity().name else {
            os_log("why am i here?")
            return Self(context: context)
        }
        let localUsersRq = NSFetchRequest<Self>(entityName: ename)
        os_log("local query for \(ename) with uuid: [\(uuid)]")
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
