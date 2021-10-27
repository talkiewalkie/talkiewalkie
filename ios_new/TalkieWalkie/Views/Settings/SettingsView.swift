//
//  SettingsView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import OSLog
import SwiftUI

struct LanguageChoice: Identifiable {
    var id: String
    var name: String
}

let availableLanguages = [
    LanguageChoice(id: "en", name: "English"),
    LanguageChoice(id: "fr", name: "Français"),
    LanguageChoice(id: "de", name: "Deutsch"),
    LanguageChoice(id: "ru", name: "Русский"),
]

enum SocialMediaType {
    case Twitter, Instagram, Snapchat, TikTok, Discord

    var appURL: URL? {
        switch self {
        case .Twitter:
            return URL(string: "twitter://user?screen_name=talkiewalkieapp")!
        case .Instagram:
            return URL(string: "instagram://user?username=talkiewalkieapp")!
        default:
            return nil
        }
    }

    var webURL: URL {
        switch self {
        case .Twitter:
            return URL(string: "https://twitter.com/talkiewalkieapp")!
        case .Instagram:
            return URL(string: "https://instagram.com/talkiewalkieapp")!
        case .Snapchat:
            return URL(string: "https://www.snapchat.com/add/talkiewalkieapp")!
        case .TikTok:
            return URL(string: "https://www.tiktok.com/@talkiewalkieapp")!
        case .Discord:
            return URL(string: "https://discord.gg/HmpvaZTDmv")!
        }
    }
}

struct SettingsView: View {
    @Binding var show: Bool
    @EnvironmentObject var authState: AuthState

    var languageSelection: Int {
        guard let index = availableLanguages.firstIndex(where: { lang in
            lang.id == authState.me?.locale
        }) else { return 0 }

        return index
    }

    func openSocialMedia(socialMedia: SocialMediaType) {
        let appURL = socialMedia.appURL, webURL = socialMedia.webURL

        if let url = appURL, UIApplication.shared.canOpenURL(url) {
            UIApplication.shared.open(url)
        } else {
            UIApplication.shared.open(webURL)
        }
    }

    var body: some View {
        let languageBinding = Binding<Int>(
            get: { languageSelection }, set: { _ in }
        )

        NavigationView {
            List {
                Section {
                    NavigationLink(destination: GeneralSettingsView()) {
                        Text(LocalizedStringKey("General"))
                    }

                    #if DEBUG
                        NavigationLink(destination: DebugSettingsInfoView()) {
                            Text(LocalizedStringKey("Debug Info"))
                        }
                    #endif

                    Button(action: {
                        if let appSettings = URL(string: UIApplication.openSettingsURLString) {
                            UIApplication.shared.open(appSettings)
                        }
                    }) {
                        Picker(selection: languageBinding, label: Text(LocalizedStringKey("Language"))) {
                            ForEach(availableLanguages.indices) { index in
                                Text(availableLanguages[index].name).tag(index)
                            }
                        }
                    }.buttonStyle(PlainButtonStyle())

                    Button(action: {
                        // Intercom.presentMessenger()
                    }, label: {
                        NavigationLink(LocalizedStringKey("Contact support"), destination: EmptyView())
                    }).buttonStyle(PlainButtonStyle())

                    NavigationLink(destination: AboutSettingsView()) {
                        Text(LocalizedStringKey("About"))
                    }
                }

                Section(header: HStack {
                    Image(systemName: "bubble.right")
                    Text(LocalizedStringKey("Follow us on Social media!"))
                }) {
                    Button(action: { openSocialMedia(socialMedia: .Twitter) }) {
                        HStack {
                            Image("twitter")
                                .resizable().aspectRatio(contentMode: .fit)
                                .frame(height: 20)
                            Text("Twitter")
                        }
                    }
                    Button(action: { openSocialMedia(socialMedia: .Instagram) }) {
                        HStack {
                            Image("instagram")
                                .resizable().aspectRatio(contentMode: .fit)
                                .frame(height: 20)
                            Text("Instagram")
                        }
                    }
                    Button(action: { openSocialMedia(socialMedia: .Snapchat) }) {
                        HStack {
                            Image("snapchat")
                                .resizable().aspectRatio(contentMode: .fit)
                                .frame(height: 20)
                            Text("Snapchat")
                        }
                    }

                    Button(action: { openSocialMedia(socialMedia: .TikTok) }) {
                        HStack {
                            Image("tiktok")
                                .resizable().aspectRatio(contentMode: .fit)
                                .frame(height: 20)
                            Text("TikTok")
                        }
                    }

                    Button(action: { openSocialMedia(socialMedia: .Discord) }) {
                        HStack {
                            Image("discord")
                                .resizable().aspectRatio(contentMode: .fit)
                                .frame(height: 20)
                            Text("Discord")
                        }
                    }
                }

                Section(header: HStack {
                    Text(LocalizedStringKey("Terms"))
                }) {
                    Link(LocalizedStringKey("Privacy policy"), destination: URL(string: "https://talkiewalkie.app/privacy")!)
                    Link(LocalizedStringKey("Terms of use"), destination: URL(string: "https://talkiewalkie.app/terms")!)
                }

                #if DEBUG
                    Section {
                        Button("Log out") {
                            os_log("logging out")
                            show = false
                            DispatchQueue.global(qos: .background).async { authState.logout() }
                        }
                        Button("Wipe out local state") {
                            os_log("wiping out core data")
                            authState.backgroundMoc.perform {
                                authState.cleanCoreData(context: authState.backgroundMoc)
                            }
                        }
                    }
                #endif
            }
            .listStyle(InsetGroupedListStyle())
            .navigationBarTitle(LocalizedStringKey("Settings"), displayMode: .inline)
            .navigationBarItems(leading: CloseButton(show: $show))
        }
    }
}

struct GeneralSettingsView: View {
    @EnvironmentObject var authState: AuthState
    @Environment(\.colorScheme) var colorScheme

    @AppStorage("addWatermark") var addWatermark: Bool = true
    @AppStorage("isDarkMode") var isDarkMode: Bool = false

    @State var showProCTA = false

    var body: some View {
        let fakeWatermarkBinding = Binding<Bool>(
            get: { true },
            set: { _ in showProCTA = true }
        )

        List {
            Section {
                Toggle(LocalizedStringKey("Watermark"), isOn: fakeWatermarkBinding)
            }

            Section(
                header: Text("Appearance")
            ) {
                Picker(selection: $isDarkMode, label: Text("Picker")) {
                    Text("Light").tag(false)
                    Text("Dark").tag(true)
                }
                .pickerStyle(SegmentedPickerStyle())
            }

            Section {}
        }
        .listStyle(InsetGroupedListStyle())
        .navigationBarTitle(LocalizedStringKey("General"), displayMode: .inline)
        .sheetWithThemeEnvironment(colorScheme: colorScheme, isPresented: $showProCTA) {
            EmptyView()
        }
    }
}

struct DebugSettingsInfoView: View {
    @EnvironmentObject var authState: AuthState

    var body: some View {
        List {
            Section {
                HStack {
                    Text("Device ID")
                    Spacer()
                }
            }
        }
    }
}

struct AboutSettingsView: View {
    var body: some View {
        List {
            Section {
                HStack {
                    Text("App version")
                    Spacer()
                    Text("v.\(Bundle.main.versionNumber) (\(Bundle.main.buildNumber))")
                        .foregroundColor(.secondary)
                }

                NavigationLink(destination: LegalSettingsView()) {
                    Text(LocalizedStringKey("Legal"))
                }
            }
        }
        .listStyle(InsetGroupedListStyle())
        .navigationBarTitle("About", displayMode: .inline)
    }
}

struct LegalSettingsView: View {
    var body: some View {
        ScrollView {
            VStack(alignment: .leading) {
                Text("Copyright information")
                    .font(.title2)

                Divider()

                Text("Tweemoji")
                    .fontWeight(.bold)

                Text("The emojis used in this app are provided by Tweemoji (tweemoji.twitter.com) and are licensed under CC BY 4.0 (creativecommons.org/licenses/by/4.0).")
            }
            .padding()
        }
        .navigationBarTitle(LocalizedStringKey("Legal"), displayMode: .inline)
    }
}

struct SettingsView_Previews: PreviewProvider {
    static var previews: some View {
        SettingsView(show: .constant(true))
    }
}
