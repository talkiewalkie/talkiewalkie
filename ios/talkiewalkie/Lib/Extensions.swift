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
}
