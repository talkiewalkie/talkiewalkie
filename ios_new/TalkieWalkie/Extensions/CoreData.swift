//
//  CoreData.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 12/10/2021.
//

import CoreData
import Foundation
import OSLog

extension NSPredicate {
    static var all = NSPredicate(format: "TRUEPREDICATE")
    static var none = NSPredicate(format: "FALSEPREDICATE")
}

extension NSManagedObjectContext {
    func saveOrLogError() {
        do { try save() }
        catch { os_log("Failed to save coredata: \(error.localizedDescription)") }
    }

    func executeOrLogError(_ request: NSPersistentStoreRequest) {
        do { try execute(request) }
        catch { os_log("failed to execute request: \(error.localizedDescription)") }
    }
}

extension NSManagedObject {
    static func getByUuidOrCreate(_ uuid: UUID, context: NSManagedObjectContext) -> Self {
        guard let ename = Self.entity().name else {
            os_log("why am i here?")
            return Self(context: context)
        }

        let localUsersRq = NSFetchRequest<Self>(entityName: ename)
        localUsersRq.predicate = NSPredicate(format: "uuid = %@", uuid.uuidString)
        let localUsers = (try? context.fetch(localUsersRq)) ?? []

        if let me = localUsers.first {
            os_log(.debug, "[coredata:\(ename)] found item for uuid:[\(uuid)]")
            return me
        } else {
            os_log(.debug, "[coredata:\(ename)] creating item for uuid:[\(uuid)]")
            let new = Self(context: context)
            try? context.save()

            return new
        }
    }
}
