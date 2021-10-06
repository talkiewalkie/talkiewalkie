//
//  Extensions.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 04/10/2021.
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

extension String {
    func uuidOrThrow() -> UUID {
        return UUID(uuidString: self)!
    }
}

extension Optional where Wrapped == String {
    var repr: String {
        switch self {
        case .some(let str):
            return str
        case .none:
            return "[nil]"
        }
    }
}
