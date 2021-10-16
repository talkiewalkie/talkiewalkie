//
//  UserStore.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import CoreData
import FirebaseAuth
import Foundation
import OSLog
import SwiftUI

class UserStore: ObservableObject {
    private var coredataCtx: NSManagedObjectContext
    
    init(_ ctx: NSManagedObjectContext) {
        self.coredataCtx = ctx
    }
    
    static var languageCode: String {
        return Locale.current.languageCode ?? "en"
    }
    
    var me: Me { Me.fromCache(context: coredataCtx)! }

    static func openSettings() {
        if let appSettings = URL(string: UIApplication.openSettingsURLString) {
            UIApplication.shared.open(appSettings)
        }
    }
    
    func logout() {
        do { try Auth.auth().signOut() }
        catch {
            os_log(.error, "failed to signout!")
            return
        }
        
        // Clear coredata on logout
        // No strong candidate to do this better: https://stackoverflow.com/questions/1077810
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: User.entity().name!)))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: Me.entity().name!)))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: Conversation.entity().name!)))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: Message.entity().name!)))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: MessageContent.entity().name!)))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: TextMessage.entity().name!)))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: VoiceMessage.entity().name!)))
        self.coredataCtx.saveOrLogError()
    }
}
