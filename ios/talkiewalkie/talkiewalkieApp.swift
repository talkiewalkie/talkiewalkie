//
//  talkiewalkieApp.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import SwiftUI

import CoreData
import OSLog
import Firebase

@main
struct talkiewalkieApp: App {
    @UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate

    var body: some Scene {
        WindowGroup {
            let auth = RootViewModel(appDelegate.persistentContainer.viewContext)

            RootView(vm: auth)
                .environment(\.managedObjectContext, appDelegate.persistentContainer.viewContext)
        }
    }
}

class AppDelegate: NSObject, UIApplicationDelegate {
    func application(_: UIApplication, didFinishLaunchingWithOptions _: [UIApplication.LaunchOptionsKey: Any]? = nil) -> Bool {
        FirebaseApp.configure()
        showCoreDataFilePath()
        return true
    }

    lazy var persistentContainer: NSPersistentContainer = {
        let container = NSPersistentContainer(name: "User")

        container.loadPersistentStores { _, error in
            container.viewContext.automaticallyMergesChangesFromParent = true

            if let error = error {
                fatalError("Unable to load persistent stores: \(error)")
            }
        }

        return container
    }()
}

// Helper func when debugging to inspect coredata content
// Run e.g. `sqlite3 {path}/User.sqlite
// Sourced from https://stackoverflow.com/a/49044302
private func showCoreDataFilePath() {
    guard let path = FileManager.default.urls(for: .libraryDirectory, in: .userDomainMask).last else { return }
    os_log("CoreData sqlite files: '\(path)Application\\ Support'")
}
