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
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: User.fetchRequest()))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: Me.fetchRequest()))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: Conversation.fetchRequest()))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: Message.fetchRequest()))
        
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: MessageContent.fetchRequest()))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: TextMessage.fetchRequest()))
        self.coredataCtx.executeOrLogError(NSBatchDeleteRequest(fetchRequest: VoiceMessage.fetchRequest()))
    }
}
