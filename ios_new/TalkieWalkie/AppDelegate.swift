//
//  AppDelegate.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 06.10.21.
//

import UIKit
import Firebase
import FirebaseMessaging

@main
class AppDelegate: UIResponder, UIApplicationDelegate {

    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        // Firebase
        FirebaseConfiguration.shared.setLoggerLevel(.min)
        FirebaseApp.configure()
        
        Messaging.messaging().delegate = self
        
        // Push Notifcations
        UNUserNotificationCenter.current().delegate = self
        
        application.registerForRemoteNotifications()
        
        
        if let path = FileManager.default.urls(for: .libraryDirectory, in: .userDomainMask).last {
            print("CoreData sqlite files: '\(path)Application\\ Support'")
        }
        
        // HACK: display LaunchScreen 1s longer
        sleep(1)
        
        return true
    }

    // MARK: UISceneSession Lifecycle
    func application(_ application: UIApplication, configurationForConnecting connectingSceneSession: UISceneSession, options: UIScene.ConnectionOptions) -> UISceneConfiguration {
        // Called when a new scene session is being created.
        // Use this method to select a configuration to create the new scene with.
        return UISceneConfiguration(name: "Default Configuration", sessionRole: connectingSceneSession.role)
    }

    func application(_ application: UIApplication, didDiscardSceneSessions sceneSessions: Set<UISceneSession>) {
        // Called when the user discards a scene session.
        // If any sessions were discarded while the application was not running, this will be called shortly after application:didFinishLaunchingWithOptions.
        // Use this method to release any resources that were specific to the discarded scenes, as they will not return.
    }
    
    func application(_ application: UIApplication,
                       didReceiveRemoteNotification userInfo: [AnyHashable: Any],
                       fetchCompletionHandler completionHandler: @escaping (UIBackgroundFetchResult)
                         -> Void) {
        
        #if DEBUG
        print(userInfo)
        #endif

        completionHandler(UIBackgroundFetchResult.newData)
    }

}


extension AppDelegate : UNUserNotificationCenterDelegate {
    func userNotificationCenter(_ center: UNUserNotificationCenter,
                              willPresent notification: UNNotification,
                              withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions)
                                -> Void) {
        let userInfo = notification.request.content.userInfo

        #if DEBUG
        print(userInfo)
        #endif

        // Change this to your preferred presentation option
        completionHandler([[.list, .banner, .sound]])
    }

}


extension AppDelegate : MessagingDelegate {

  func messaging(_ messaging: Messaging, didReceiveRegistrationToken fcmToken: String?) {
    guard let token = fcmToken else { return }
    
    #if DEBUG
    print("Firebase registration token: \(token)")
    #endif
  }

}
