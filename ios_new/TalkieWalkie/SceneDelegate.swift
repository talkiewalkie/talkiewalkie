//
//  SceneDelegate.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 06.10.21.
//

import CoreData
import SwiftUI
import UIKit


class Config: Decodable, ObservableObject {
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


class SceneDelegate: UIResponder, UIWindowSceneDelegate {
    var window: UIWindow?

    lazy var persistentContainer: NSPersistentContainer = {
        let container = NSPersistentContainer(name: "LocalModels")

        container.loadPersistentStores { _, error in
            container.viewContext.automaticallyMergesChangesFromParent = true

            if let error = error {
                fatalError("Unable to load persistent stores: \(error)")
            }
        }

        return container
    }()

    func scene(_ scene: UIScene, willConnectTo _: UISceneSession, options _: UIScene.ConnectionOptions) {
        // Use this method to optionally configure and attach the UIWindow `window` to the provided UIWindowScene `scene`.
        // If using a storyboard, the `window` property will automatically be initialized and attached to the scene.
        // This delegate does not imply the connecting scene or session are new (see `application:configurationForConnectingSceneSession` instead).

        #if DEBUG
        let config = Config.load(version: "dev")
        #else
        let config = Config.load(version: "prod")
        #endif
        let tooltipManager = TooltipManager()
        let hmv = HomeViewModel(persistentContainer.viewContext, config: config)

        let contentView = HomeView(homeViewModel: hmv)
            .environmentObject(UserStore(persistentContainer.viewContext))
            .addTooltip()
            .environmentObject(tooltipManager)
            .environmentObject(config)

        // Use a UIHostingController as window root view controller.
        if let windowScene = scene as? UIWindowScene {
            let window = UIWindow(windowScene: windowScene)
            window.rootViewController = UIHostingController(rootView: contentView)
            self.window = window
            window.makeKeyAndVisible()
        }
    }

    func sceneDidDisconnect(_: UIScene) {
        // Called as the scene is being released by the system.
        // This occurs shortly after the scene enters the background, or when its session is discarded.
        // Release any resources associated with this scene that can be re-created the next time the scene connects.
        // The scene may re-connect later, as its session was not necessarily discarded (see `application:didDiscardSceneSessions` instead).
    }

    func sceneDidBecomeActive(_: UIScene) {
        // Called when the scene has moved from an inactive state to an active state.
        // Use this method to restart any tasks that were paused (or not yet started) when the scene was inactive.
    }

    func sceneWillResignActive(_: UIScene) {
        // Called when the scene will move from an active state to an inactive state.
        // This may occur due to temporary interruptions (ex. an incoming phone call).
    }

    func sceneWillEnterForeground(_: UIScene) {
        // Called as the scene transitions from the background to the foreground.
        // Use this method to undo the changes made on entering the background.
    }

    func sceneDidEnterBackground(_: UIScene) {
        // Called as the scene transitions from the foreground to the background.
        // Use this method to save data, release shared resources, and store enough scene-specific state information
        // to restore the scene back to its current state.
    }
}
