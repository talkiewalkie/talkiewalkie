//
//  AuthenticationState.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 18/10/2021.
//

import CoreData
import FirebaseAuth
import Foundation
import OSLog

#if DEBUG
private let env = "dev"
#else
private let env = "prod"
#endif

enum AuthenticationState {
    case Disconnected
    case Connecting
    case Connected(AuthedGrpcApi, Me)
}

class AuthState: ObservableObject {
    let moc: NSManagedObjectContext
    var me: Me? { Me.fromCache(context: moc) }

    private var logger = Logger.withLabel("AuthState")
    private let config = Config.load(version: env)

    @Published private(set) var state = AuthenticationState.Disconnected

    init() {
        logger.debug("loading Core Data store")

        let container = NSPersistentContainer(name: "LocalModels")

        container.loadPersistentStores { _, error in
            container.viewContext.automaticallyMergesChangesFromParent = true

            if let error = error {
                fatalError("Unable to load persistent stores: \(error)")
            }
        }
        self.moc = container.viewContext
    }

    // MARK: - Utils

    var firebaseUser: FirebaseAuth.User? { Auth.auth().currentUser }

    func connect(with user: FirebaseAuth.User) {
        setConnecting()
        user.getIDTokenResult { res, err in
            guard let res = res else {
                guard let err = err else { return }
                self.logger.error("failed to get firebase token: \(err.localizedDescription)")
                return
            }

            let api = AuthedGrpcApi(url: self.config.apiUrl, token: res.token, context: self.moc)

            if let me = Me.fromCache(context: self.moc) {
                self.logger.debug("loaded user info from cache: \(me)")
                self.state = AuthenticationState.Connected(api, me)
            }

            DispatchQueue.main.async {
                let (res, _) = api.me()
                if let res = res {
                    let uuid = UUID(uuidString: res.user.uuid)!
                    let me = Me.getByUuidOrCreate(uuid, context: self.moc)
                    me.uuid = uuid
                    me.displayName = res.user.displayName
                    me.firebaseUid = self.firebaseUser?.uid

                    me.objectWillChange.send()
                    self.moc.saveOrLogError()

                    self.state = AuthenticationState.Connected(api, me)
                }
            }
        }
    }

    func logout() {
        state = AuthenticationState.Disconnected
        
        do { try Auth.auth().signOut() }
        catch {
            os_log(.error, "failed to signout!")
            return
        }
        
        // Clear coredata on logout
        // No strong candidate to do this better: https://stackoverflow.com/questions/1077810
        moc.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: User.entity().name!)))
        moc.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: Me.entity().name!)))
        moc.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: Conversation.entity().name!)))
        moc.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: Message.entity().name!)))
        moc.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: MessageContent.entity().name!)))
        moc.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: TextMessage.entity().name!)))
        moc.executeOrLogError(NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: VoiceMessage.entity().name!)))
        moc.saveOrLogError()
    }

    func setConnecting() {
        state = AuthenticationState.Connecting
    }
}

private class Config: Decodable, ObservableObject {
    private enum CodingKeys: String, CodingKey {
        case apiHost, apiPort
    }

    let apiHost: String
    let apiPort: Int

    static func load(version: String) -> Config {
        let url = Bundle.main.url(forResource: "Config.\(version)", withExtension: "plist")!
        let data = try! Data(contentsOf: url)
        let decoder = PropertyListDecoder()
        return try! decoder.decode(Config.self, from: data)
    }

    var apiUrl: URL {
        #if DEBUG
        let transport = "http://"
        #else
        let transport = "https://"
        #endif

        return URL(string: "\(transport)\(apiHost):\(apiPort)")!
    }
}
