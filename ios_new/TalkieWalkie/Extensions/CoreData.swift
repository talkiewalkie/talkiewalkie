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
        performAndWait {
            if hasChanges {
                do { try save() }
                catch { os_log(.error, "Failed to save coredata: \(error.localizedDescription)") }
            }
        }
    }

    func executeOrLogError(_ request: NSPersistentStoreRequest) -> NSPersistentStoreResult? {
        var res: NSPersistentStoreResult?
        do { res = try execute(request) }
        catch { os_log("failed to execute request: \(error.localizedDescription)") }
        return res
    }

    /// Executes the given `NSBatchDeleteRequest` and directly merges the changes to bring the given managed object context up to date.
    /// From https://stackoverflow.com/a/60266079, although it is not updating the contexts though, hence not meeting its purpose.
    /// Leaving it for reference.
    ///
    /// - Parameter batchDeleteRequest: The `NSBatchDeleteRequest` to execute.
    /// - Throws: An error if anything went wrong executing the batch deletion.
    public func deleteAndMergeChanges(using batchDeleteRequest: NSBatchDeleteRequest) {
        batchDeleteRequest.resultType = .resultTypeObjectIDs
        let result = executeOrLogError(batchDeleteRequest) as? NSBatchDeleteResult
        let changes: [AnyHashable: Any] = [NSDeletedObjectsKey: result?.result as? [NSManagedObjectID] ?? []]
        NSManagedObjectContext.mergeChanges(fromRemoteContextSave: changes, into: [self])
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
            return new
        }
    }
}
