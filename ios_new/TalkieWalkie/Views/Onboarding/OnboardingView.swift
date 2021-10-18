//
//  OnboardingView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import AVFoundation
import FirebaseAuth
import Introspect
import OSLog
import PhoneNumberKit
import SwiftUI

class OnboardingViewModel: ObservableObject {
    @Published var page: Int = Auth.auth().currentUser != nil ? 5 : 0
    @Binding var name: String
    @Published var phoneCountryCode: Int
    @Published var phoneRegionID: String?
    @Published var phoneNumber: String
    @Published var verificationCode: String = ""
    @Published var verificationID: String = ""

    init(name: Binding<String>, phoneCountryCode: Int, phoneRegionID: String, phoneNumber: String) {
        _name = name
        self.phoneCountryCode = phoneCountryCode
        self.phoneRegionID = phoneRegionID.isEmpty ? nil : phoneRegionID
        self.phoneNumber = phoneNumber
    }

    var fullPhoneNumber: String {
        "+\(phoneCountryCode)\(phoneNumber)"
    }

    func next() {
        page += 1
    }
}

struct OnboardingView: View {
    var onboardingDone: () -> Void

    @AppStorage("name") var name: String = ""
    @AppStorage("phoneCountryCode") var phoneCountryCode: Int = 33
    @AppStorage("phoneRegionID") var phoneRegionID: String = "FR"
    @AppStorage("phoneNumber") var phoneNumber: String = ""

    var body: some View {
        Onboarding(
            onboardingDone: onboardingDone,
            model: .init(name: $name,
                         phoneCountryCode: phoneCountryCode,
                         phoneRegionID: phoneRegionID,
                         phoneNumber: phoneNumber)
        )
    }
}

struct Onboarding: View {
    var onboardingDone: () -> Void

    @StateObject var model: OnboardingViewModel

    var body: some View {
        ZStack {
            Color.blue.opacity(0.5).ignoresSafeArea()

            Group {
                switch model.page {
                case 0: FirstScreen()
                case 1: TypeNameView()
                case 2: TypePhoneNumberView()
                case 3: VerifyPhoneNumberView()
                case 4: PhoneNumberSuccessView()
                case 5: TurnOnNotificationsView()
                case 6: MicrophoneAuthorizationView()
                case 7: ContactListView()
                case 8: AddTalkieWalkieContactsView()
                case 9: AddOtherContactsView()
                case 10: FinalSuccessView(onboardingDone: onboardingDone)
                default: EmptyView()
                }
            }
        }
        .environmentObject(model)
        .foregroundColor(.white)
    }
}

struct OnboardingNavControls: View {
    @Binding var page: Int

    var showPrev: Bool = true
    var showNext: Bool = true

    var loading: Bool = false

    var onNext: (() -> Void)?

    var body: some View {
        VStack {
            Spacer()
            HStack {
                if showPrev && !loading {
                    TWButton(action: {
                        page -= 1
                    }, primary: false) {
                        Image(systemName: "arrow.backward")
                            .font(.system(size: 22, weight: .heavy))
                    }
                }

                Spacer()

                if showNext {
                    TWButton(action: {
                        guard let next = onNext else { return page += 1 }
                        next()
                    }) {
                        Group {
                            if loading {
                                ProgressView()
                                    .environment(\.colorScheme, .dark)
                            } else {
                                Image(systemName: "arrow.forward")
                                    .font(.system(size: 22, weight: .heavy))
                            }
                        }
                    }
                }
            }
        }
        .frame(alignment: .bottom)
        .padding()
    }
}

struct OnboardingTitle: View {
    var text: String

    var body: some View {
        Text(text.uppercased())
            .font(.title2)
            .fontWeight(.black)
            .foregroundColor(.white)
            .frame(maxWidth: .infinity, alignment: .topLeading)
    }
}

// MARK: First Screen

struct FirstScreen: View {
    @EnvironmentObject var model: OnboardingViewModel

    var body: some View {
        ZStack {
            VStack {
                Spacer()

                Image("curves")
                    .resizable()
                    .aspectRatio(contentMode: .fit)
                    .scaleEffect(2.5)
                    .opacity(0.5)
            }
            .ignoresSafeArea()

            ZStack {
                VStack {
                    LottieView(name: "speech-therapy")
                        .frame(maxHeight: 250)

                    Spacer()
                }

                VStack {
                    Image("logo_text")
                        .resizable()
                        .aspectRatio(contentMode: .fit)
                        .scaleEffect(0.75)

                    Text("Voice messaging, done right.")
                        .font(.title2)
                        .fontWeight(.medium)
                        .foregroundColor(.white)
                }

                VStack {
                    Spacer()

                    VStack {
                        Text("By tapping on \"Sign up & Accept\", you agree to the Privacy Policy and Terms of Services.")
                            .font(.callout)
                            .fontWeight(.medium)
                            .foregroundColor(.init(.tertiaryLabel))
                            .multilineTextAlignment(.center)
                            .scaleEffect(0.9)

                        TWButton(action: model.next) {
                            Text("Sign up & Accept".uppercased())
                        }
                    }
                }
            }
            .padding()
        }
    }
}

// MARK: Type Name

struct TypeNameView: View {
    @EnvironmentObject var model: OnboardingViewModel

    @State var showInvalidNameAlert: Bool = false

    var body: some View {
        ZStack {
            VStack {
                OnboardingTitle(text: "Hi! What's your name?")

                Spacer()
            }

            TextField("Your Name", text: $model.name, onCommit: validate)
                .disableAutocorrection(true)
                .textContentType(UITextContentType.givenName)
                .multilineTextAlignment(.center)
                .accentColor(Color("Cyan"))
                .font(.title2.weight(.heavy))
                .introspectTextField { textField in
                    textField.becomeFirstResponder()
                }

            OnboardingNavControls(page: $model.page, showPrev: false, onNext: validate)
        }
        .alert(isPresented: $showInvalidNameAlert) {
            Alert(title: Text("This doesn't look like a real name! 😛"))
        }
        .padding()
    }

    func validate() {
        guard model.name.count > 0 else { return showInvalidNameAlert = true }

        model.next()
    }
}

// MARK: Type Phone

struct PhoneCountryCodeButtonView: View {
    @Binding var countryCode: Int
    @Binding var regionID: String?

    var flag: String {
        guard let regionID = regionID else { return "" }
        return getFlag(from: regionID)
    }

    var body: some View {
        TWButton(action: {}, primary: false, padding: 10) {
            Text("\(flag)+\(countryCode)")
        }
    }

    func getFlag(from regionID: String) -> String {
        regionID
            .unicodeScalars
            .map { 127_397 + $0.value }
            .compactMap(UnicodeScalar.init)
            .map(String.init)
            .joined()
    }
}

struct TypePhoneNumberView: View {
    @EnvironmentObject var model: OnboardingViewModel

    @State var showInvalidPhoneNumberAlert: Bool = false
    @State var showSendingSMSAlert: Bool = false
    @State var phoneNumberKit = PhoneNumberKit()

    @AppStorage("phoneCountryCode") var storedPhoneCountryCode: Int = 33
    @AppStorage("phoneRegionID") var storedPhoneRegionID: String = "FR"
    @AppStorage("phoneNumber") var storedPhoneNumber: String = ""

    var body: some View {
        let phoneNumber = Binding<String>(
            get: {
                model.phoneNumber
            },
            set: {
                if $0.count > model.phoneNumber.count + 1 {
                    self.parseNumber($0) // AutoFill
                } else {
                    model.phoneNumber = $0
                }
            }
        )

        ZStack {
            VStack {
                OnboardingTitle(text: "Hey \(model.name), I need your phone number to identify you.")

                Spacer()
            }

            HStack {
                PhoneCountryCodeButtonView(countryCode: $model.phoneCountryCode, regionID: $model.phoneRegionID)

                TextField("Phone number", text: phoneNumber, onCommit: validate)
                    .disableAutocorrection(true)
                    .keyboardType(.numberPad)
                    .textContentType(UITextContentType.telephoneNumber)
                    .accentColor(Color("Cyan"))
                    .font(.title2.weight(.heavy))
                    .introspectTextField { textField in
                        textField.becomeFirstResponder()
                    }
            }
            .padding(.horizontal)

            OnboardingNavControls(page: $model.page, showPrev: false, onNext: validate)
        }
        .padding()
        .alert(isPresented: $showSendingSMSAlert) {
            Alert(title: Text("Awesome! 🙌 I'm sending a code to: \(model.fullPhoneNumber)"),
                  dismissButton: .default(Text("OK"), action: model.next))
        }
        .overlay(
            EmptyView()
                .alert(isPresented: $showInvalidPhoneNumberAlert) {
                    Alert(title: Text("Please enter a valid phone number 😙"))
                }
        )
    }

    func parseNumber(_ numberString: String, withRegion: String? = nil) -> Bool {
        guard let phoneNumber = try? phoneNumberKit.parse(
            numberString,
            withRegion: withRegion ?? PhoneNumberKit.defaultRegionCode()
        ) else { return false }

        model.phoneNumber = "\(phoneNumber.nationalNumber)"; storedPhoneNumber = model.phoneNumber
        model.phoneCountryCode = Int(phoneNumber.countryCode); storedPhoneCountryCode = model.phoneCountryCode
        model.phoneRegionID = phoneNumber.regionID ?? ""; storedPhoneRegionID = model.phoneRegionID ?? ""

        return true
    }

    func validate() {
        let valid = parseNumber(model.phoneNumber, withRegion: model.phoneRegionID)

        guard valid else { return showInvalidPhoneNumberAlert = true }

        showSendingSMSAlert = true

        os_log("verif id: '\(model.verificationID)'")
        os_log("verif code: '\(model.verificationCode)'")

        PhoneAuthProvider.provider()
            .verifyPhoneNumber(model.fullPhoneNumber, uiDelegate: nil) { verificationID, error in
                if let error = error {
                    os_log(.error, "failed to verify phone number '\(model.fullPhoneNumber)':  \(error.localizedDescription)") // TODO:
                    return
                }

                guard let verificationID = verificationID, verificationID != "" else { return }
                model.verificationID = verificationID
            }
    }
}

struct VerifyPhoneNumberView: View {
    @EnvironmentObject var model: OnboardingViewModel

    @State var showInvalidVerificationCodeAlert: Bool = false
    @State var loading: Bool = false

    var body: some View {
        ZStack {
            VStack {
                VStack(alignment: .leading) {
                    OnboardingTitle(text: "Please enter the code I sent you.")
                    Text("Sent to \(model.fullPhoneNumber)".uppercased())
                        .font(.caption.weight(.semibold))
                        .foregroundColor(Color("DarkBlue"))
                }

                Spacer()
            }

            TextField("••••••", text: $model.verificationCode, onCommit: validate)
                .disableAutocorrection(true)
                .multilineTextAlignment(.center)
                .keyboardType(.numberPad)
                .textContentType(UITextContentType.oneTimeCode)
                .accentColor(Color("Cyan"))
                .font(.system(size: 50, weight: .heavy))
                .introspectTextField { textField in
                    textField.becomeFirstResponder()
                }

            OnboardingNavControls(page: $model.page, loading: loading, onNext: validate)
        }
        .alert(isPresented: $showInvalidVerificationCodeAlert) {
            Alert(title: Text("Please enter a valid code 😙"))
        }
        .padding()
    }

    func validate() {
        guard model.verificationCode.count == 6 else { return showInvalidVerificationCodeAlert = true }

        loading = true

        let credential = PhoneAuthProvider.provider().credential(
            withVerificationID: model.verificationID,
            verificationCode: model.verificationCode
        )

        Auth.auth().signIn(with: credential) { _, error in
            loading = false

            if let error = error {
                let authError = error as NSError
                os_log(.error, "failed to auth: \(authError.description)")
                return
            }

            // User has signed in successfully and currentUser object is valid
            let currentUserInstance = Auth.auth().currentUser
            os_log(.debug, "Success! \(currentUserInstance!)")
            model.next()
        }
    }
}

struct PhoneNumberSuccessView: View {
    @EnvironmentObject var model: OnboardingViewModel

    var body: some View {
        ZStack {
            LottieView(name: "confetti2")
        }
        .padding()
        .onAppear {
            DispatchQueue.main.asyncAfter(deadline: .now() + 2) {
                model.next()
            }
        }
    }
}

// MARK: Allow Notifications

struct TurnOnNotificationsView: View {
    @EnvironmentObject var model: OnboardingViewModel

    @State var notificationAuthorizationStatus: AVCaptureDevice.AuthorizationStatus?

    var body: some View {
        ZStack {
            VStack {
                OnboardingTitle(text: "TalkieWalkie works better with notifications on.")
                Spacer()
            }

            VStack {
                LottieView(name: "notification")
                    .scaleEffect(0.8)
                    .frame(maxHeight: 400)

                TWButton(action: {
                    let center = UNUserNotificationCenter.current()
                    center.requestAuthorization(options: [.alert, .sound, .badge]) { _, error in
                        if let error = error {
                            os_log(.error, "\(error.localizedDescription)")
                        }
                    }
                    model.next()
                }, primary: true) {
                    Text("Activate!")
                }
                
                Spacer()
                HStack {
                    Spacer()

                    TWButton(action: model.next, primary: false) {
                        Text("Later".uppercased())
                    }
                    .opacity(0.75)
                }
            }
        }
        .padding()
    }
}

// MARK: Allow microphone

struct MicrophoneAuthorizationView: View {
    @EnvironmentObject var model: OnboardingViewModel

    @State var audioAuthorizationStatus: AVCaptureDevice.AuthorizationStatus?

    func authorizeAudio() {
        DispatchQueue.global(qos: .userInitiated).async {
            AVCaptureDevice.authorizeAudio(completion: { audioStatus in
                print("Status: ", audioStatus, " superstatus: ", audioStatus.superStatus)
                DispatchQueue.main.async {
                    withAnimation {
                        audioAuthorizationStatus = audioStatus
                        print("Authorized?: ", audioAuthorizationStatus?.superStatus == .authorized)
                    }
                }
            })
        }
    }

    func openSettings() {
        if let appSettings = URL(string: UIApplication.openSettingsURLString) {
            UIApplication.shared.open(appSettings)
        }
    }

    var body: some View {
        ZStack {
            VStack {
                OnboardingTitle(text: "For TalkieWalkie to work, allow microphone access.")

                Spacer()
            }
            if audioAuthorizationStatus?.superStatus == .authorized {
                authorizationSuccessView
            } else if audioAuthorizationStatus?.superStatus == .denied {
                authorizationDeniedView
            } else {
                requestAuthorizationView
            }
        }
        .padding()
    }

    var requestAuthorizationView: some View {
        ZStack {
            LottieView(name: "52786-recording-bubble")
                .frame(maxHeight: 400)

            VStack {
                Spacer()

                TWButton(action: {
                    authorizeAudio()
                }, compact: false) {
                    Text("Allow".uppercased())
                }
            }
        }
    }

    var authorizationSuccessView: some View {
        ZStack {}
            .onAppear(perform: model.next)
    }

    var authorizationDeniedView: some View {
        ZStack {
            VStack {
                LottieView(name: "mic-off")
                    .frame(maxHeight: 200)

                Text("Go to settings, and toggle \"Microphone\"".uppercased())

                Image("microphone_toggle")
                    .resizable()
                    .aspectRatio(contentMode: .fit)
            }

            VStack {
                Spacer()

                TWButton(action: {
                    openSettings()
                }, compact: false) {
                    Text("Open phone settings".uppercased())
                }
            }
        }
    }
}

struct AddTalkieWalkieContactsView: View {
    @EnvironmentObject var model: OnboardingViewModel

    var body: some View {
        ZStack {
            Text("Hello")

            OnboardingNavControls(page: $model.page)
        }
        .padding()
    }
}

struct AddOtherContactsView: View {
    @EnvironmentObject var model: OnboardingViewModel

    var body: some View {
        ZStack {
            Text("Hello")

            OnboardingNavControls(page: $model.page)
        }
        .padding()
    }
}

// MARK: Final Screen

struct FinalSuccessView: View {
    @EnvironmentObject var model: OnboardingViewModel
    var onboardingDone: () -> Void

    var body: some View {
        ZStack {
            LottieView(name: "check_success", loopMode: .playOnce)

            VStack {
                Spacer()

                TWButton(action: { onboardingDone() }) {
                    Text("Let's start!".uppercased())
                }
            }
        }
        .padding()
    }
}

struct Onboarding_Previews: PreviewProvider {
    static var previews: some View {
        TestView()
    }

    struct TestView: View {
        @StateObject var model: OnboardingViewModel = {
            let model = OnboardingViewModel(name: .constant("Alex"), phoneCountryCode: 33, phoneRegionID: "FR", phoneNumber: "")
            model.page = 6
            return model
        }()

        var body: some View {
            Onboarding(onboardingDone: {}, model: model)
        }
    }
}
