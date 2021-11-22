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
    let persistentContainer: NSPersistentContainer
    var readContext: NSManagedObjectContext { persistentContainer.viewContext }
    private let writeContext: NSManagedObjectContext

    private lazy var backgroundQueue: OperationQueue = {
        let persistentContainerQueue = OperationQueue()
        persistentContainerQueue.maxConcurrentOperationCount = 1

        return persistentContainerQueue
    }()

    var me: Me? { Me.fromCache(context: readContext) }
    var meOrThrow: Me { me! }
    private lazy var backgroundMe: Me? = { Me.fromCache(context: writeContext) }()

    private var logger = Logger.withLabel("AuthState")
    private let config = Config.load(version: env)

    @Published private(set) var state = AuthenticationState.Disconnected

    init() {
        logger.debug("loading Core Data store")

        let persistentContainer = NSPersistentContainer(name: "LocalModels")

        persistentContainer.loadPersistentStores { _, error in
            persistentContainer.viewContext.mergePolicy = NSMergePolicy.error
            persistentContainer.viewContext.automaticallyMergesChangesFromParent = true
            if let error = error {
                fatalError("Unable to load persistent stores: \(error)")
            }
        }
        self.persistentContainer = persistentContainer

        self.writeContext = persistentContainer.newBackgroundContext()
        writeContext.mergePolicy = NSMergePolicy.error
        writeContext.automaticallyMergesChangesFromParent = true
    }

    // MARK: - Utils

    var firebaseUser: FirebaseAuth.User? { Auth.auth().currentUser }

    func connect(with fbUser: FirebaseAuth.User) {
        setConnecting()
        fbUser.getIDTokenResult { res, err in
            guard let res = res else {
                guard let err = err else { return }
                self.logger.error("failed to get firebase token: \(err.localizedDescription)")
                return
            }

            let api = AuthedGrpcApi.withUrlAndToken(url: self.config.apiUrl, token: res.token, writer: {fn in self.withWriteContext { ctx, me in fn(ctx, me)}})

            if let me = Me.fromCache(context: self.readContext) {
                self.logger.debug("loaded my user info from cache")
                self.state = AuthenticationState.Connected(api, me)
            }

            DispatchQueue.global(qos: .background).async {
                let (res, _) = api.me()
                if let res = res {
                    let uuid = UUID(uuidString: res.user.uuid)!

                    // This is a very particuliar moment in the app startup and the only place we'll ever need blocking save
                    self.readContext.performAndWait {
                        let me = Me.fromCache(context: self.readContext) ?? Me(context: self.readContext)
                        me.uuid = uuid
                        me.displayName = res.user.displayName
                        me.firebaseUid = fbUser.uid
                        try! self.readContext.save()
                        DispatchQueue.main.async { self.state = AuthenticationState.Connected(api, me) }
                    }
                }
            }
        }
    }

    func logout() {
        state = AuthenticationState.Disconnected

        DispatchQueue.global(qos: .background).async {
            do { try Auth.auth().signOut() }
            catch {
                os_log(.error, "failed to signout!")
                return
            }
        }

        withWriteContext { ctx, _ in self.cleanCoreData(context: ctx) }
    }

    func cleanCoreData(context: NSManagedObjectContext) {
        // Clear coredata on logout
        // No strong candidate to do this better: https://stackoverflow.com/questions/1077810
        context.deleteAndMergeChanges(using: NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: User.entity().name!)))
        context.deleteAndMergeChanges(using: NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: Me.entity().name!)))
        context.deleteAndMergeChanges(using: NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: Conversation.entity().name!)))
        context.deleteAndMergeChanges(using: NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: Message.entity().name!)))
        context.deleteAndMergeChanges(using: NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: MessageContent.entity().name!)))
        context.deleteAndMergeChanges(using: NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: TextMessage.entity().name!)))
        context.deleteAndMergeChanges(using: NSBatchDeleteRequest(fetchRequest: NSFetchRequest(entityName: VoiceMessage.entity().name!)))
    }

    func setConnecting() {
        state = AuthenticationState.Connecting
    }

    func withWriteContext(block: @escaping (_ context: NSManagedObjectContext, _ me: Me?) -> Void) {
        backgroundQueue.addOperation {
            self.writeContext.performAndWait {
                block(self.writeContext, self.backgroundMe)
                try! self.writeContext.save()
            }
        }
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
