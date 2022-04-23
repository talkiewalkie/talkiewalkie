//
//  LocalNotificationManager.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 28.10.21.
//

import SwiftUI

struct LocalNotification: Identifiable {
    let id = UUID()
    var title: String
    var message: String
    var image: UIImage?
    var date: Date?
}

public extension View {
    func addLocalNotification() -> some View {
        modifier(
            LocalNotificationModifier()
        )
    }
}

struct LocalNotificationModifier: ViewModifier {
    @EnvironmentObject private var manager: LocalNotificationManager

    func body(content: Content) -> some View {
        GeometryReader { _ in
            ZStack {
                content

                ZStack {
                    ForEach(manager.notifications) {
                        LocalNotificationView(notification: $0)
                            .transition(
                                .asymmetric(insertion: .move(edge: .top),
                                            removal: .scale.combined(with: .opacity))
                            )
                    }
                }
                .padding(.horizontal, 10)
                .frame(maxHeight: .infinity, alignment: .top)
            }
        }
    }
}

struct LocalNotificationView: View {
    var notification: LocalNotification

    var formattedDate: String {
        if let date = notification.date {
            return Self.dateFormatter.string(from: date)
        }
        return "now"
    }

    var body: some View {
        HStack {
            Image(uiImage: notification.image ?? UIImage(named: "AppIcon")!)
                .resizable()
                .aspectRatio(1, contentMode: .fill)
                .frame(width: DrawingConstraints.imageSize,
                       height: DrawingConstraints.imageSize)
                .cornerRadius(10)
                .clipped()

            HStack(alignment: .top) {
                VStack(alignment: .leading) {
                    Text(notification.title)
                        .bold()
                    Text(notification.message)
                        .lineLimit(1)
                }

                Spacer()

                Text(formattedDate)
                    .font(.subheadline)
                    .foregroundColor(.secondary)
            }
        }
        .padding(DrawingConstraints.padding)
        .background(
            VisualEffectView(effect: UIBlurEffect(style: .systemThinMaterial))
                .cornerRadius(20)
        )
        .shadow(color: .black.opacity(0.15), radius: 3)
    }

    static let dateFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "HH:mm"
        return formatter
    }()

    enum DrawingConstraints {
        static let imageSize: CGFloat = 50
        static let padding: CGFloat = 12
    }
}

class LocalNotificationManager: ObservableObject {
    @Published fileprivate var notifications: [LocalNotification] = []

    public func showNotification(_ notification: LocalNotification, displayDuration: Double = 2) {
        withAnimation(.easeInOut) {
            self.notifications.append(notification)
        }
        withAnimation(.easeInOut.delay(0.1)) {
            self.notifications.removeFirst(self.notifications.count - 1)
        }

        DispatchQueue.main.asyncAfter(deadline: .now() + displayDuration) {
            if let index = self.notifications.firstIndex(where: { $0.id == notification.id }) {
                _ = withAnimation(.easeInOut) {
                    self.notifications.remove(at: index)
                }
            }
        }
    }
}

struct LocalNotificationManager_Previews: PreviewProvider {
    static var previews: some View {
        VStack {
            TestView()
                .addLocalNotification()
                .environmentObject(LocalNotificationManager())
        }
    }

    struct TestView: View {
        @EnvironmentObject var localNotificationManager: LocalNotificationManager

        var body: some View {
            VStack {
                Button(action: {
                    localNotificationManager.showNotification(
                        LocalNotification(title: "Author", message: "Message")
                    )
                }) {
                    Text("Click")
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
    }
}
