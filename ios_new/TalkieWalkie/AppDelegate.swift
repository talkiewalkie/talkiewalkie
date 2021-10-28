//
//  AppDelegate.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 06.10.21.
//

import Firebase
import FirebaseMessaging
import OSLog
import UIKit
import GiphyUISDK

@main
class AppDelegate: UIResponder, UIApplicationDelegate {
    let auth = AuthState()

    func application(_ application: UIApplication, didFinishLaunchingWithOptions _: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        FirebaseConfiguration.shared.setLoggerLevel(.min)
        FirebaseApp.configure()

        Auth.auth().addStateDidChangeListener { _, newUser in
            if let newUser = newUser {
                Messaging.messaging().subscribe(toTopic: newUser.uid)
                os_log(.debug, "subscribed to '\(newUser.uid)'")
                self.auth.connect(with: newUser)
            } else {
                // TODO: unsubscribe, but we need to get the current user before the new one is nil from that thread
                // TODO: to do so, which I don't know how to do just yet.
                // Messaging.messaging().unsubscribe(fromTopic: existingUser.uid)
                self.auth.logout()
            }
        }

        Messaging.messaging().delegate = self
        Messaging.messaging().subscribe(toTopic: "all")
        if let fbu = Auth.auth().currentUser {
            Messaging.messaging().subscribe(toTopic: fbu.uid)
            os_log(.debug, "subscribed to '\(fbu.uid)'")
        }

        // Push Notifcations
        UNUserNotificationCenter.current().delegate = self

        application.registerForRemoteNotifications()

        if let path = FileManager.default.urls(for: .libraryDirectory, in: .userDomainMask).last {
            os_log(.debug, "CoreData sqlite files: '\(path)Application\\ Support'")
        }
        
        // Giphy
        Giphy.configure(apiKey: "9eyAVdK8MCwzgBwTN1vTi0cNoIHNQ3oM")

        // HACK: display LaunchScreen 1s longer
        sleep(1)

        return true
    }

    // MARK: - UISceneSession Lifecycle

    func application(_: UIApplication, configurationForConnecting connectingSceneSession: UISceneSession, options _: UIScene.ConnectionOptions) -> UISceneConfiguration {
        // Called when a new scene session is being created.
        // Use this method to select a configuration to create the new scene with.
        return UISceneConfiguration(name: "Default Configuration", sessionRole: connectingSceneSession.role)
    }

    func application(_: UIApplication, didDiscardSceneSessions _: Set<UISceneSession>) {
        // Called when the user discards a scene session.
        // If any sessions were discarded while the application was not running, this will be called shortly after application:didFinishLaunchingWithOptions.
        // Use this method to release any resources that were specific to the discarded scenes, as they will not return.
    }

    func application(_: UIApplication,
                     didReceiveRemoteNotification notification: [AnyHashable: Any],
                     fetchCompletionHandler completionHandler: @escaping (UIBackgroundFetchResult)
                         -> Void)
    {
        #if DEBUG
            os_log("background notification received with \(notification)")
        #endif

        completionHandler(UIBackgroundFetchResult.newData)
    }
    
    func applicationWillTerminate(_ application: UIApplication) {
        if case .Connected(let api, _) = self.auth.state {
            // TODO: send last connected at
        }
    }
}

extension AppDelegate: UNUserNotificationCenterDelegate {
    func userNotificationCenter(_: UNUserNotificationCenter,
                                willPresent notification: UNNotification,
                                withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions)
                                    -> Void)
    {
        let userInfo = notification.request.content.userInfo

        if case .Connected(let api, _) = self.auth.state {
            if let uuidString = (userInfo["uuid"] as? String), let uuid = UUID(uuidString: uuidString) {
                let (user, _) = api.userByUuid(uuid)
                if let user = user {
                    User.upsert(user, context: self.auth.moc)
                }
            }
        }

        #if DEBUG
            os_log("notification center received notif with \(userInfo)")
        #endif

        // Change this to your preferred presentation option
        completionHandler([[.list, .banner, .sound]])
    }
}

extension AppDelegate: MessagingDelegate {
    func messaging(_: Messaging, didReceiveRegistrationToken _: String?) {}
}
