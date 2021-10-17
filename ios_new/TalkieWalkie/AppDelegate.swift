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

@main
class AppDelegate: UIResponder, UIApplicationDelegate {
    func application(_ application: UIApplication, didFinishLaunchingWithOptions _: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        // Firebase
        FirebaseConfiguration.shared.setLoggerLevel(.min)
        FirebaseApp.configure()

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

        // HACK: display LaunchScreen 1s longer
        sleep(1)

        return true
    }

    // MARK: UISceneSession Lifecycle

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
                     didReceiveRemoteNotification userInfo: [AnyHashable: Any],
                     fetchCompletionHandler completionHandler: @escaping (UIBackgroundFetchResult)
                         -> Void)
    {
        #if DEBUG
            os_log("notification received with \(userInfo)")
        #endif

        completionHandler(UIBackgroundFetchResult.newData)
    }
}

extension AppDelegate: UNUserNotificationCenterDelegate {
    func userNotificationCenter(_: UNUserNotificationCenter,
                                willPresent notification: UNNotification,
                                withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions)
                                    -> Void)
    {
        let userInfo = notification.request.content.userInfo

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
