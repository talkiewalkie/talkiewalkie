//
//  ContactListView.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 16/10/2021.
//

import Contacts
import CoreData
import FirebaseAuth
import PhoneNumberKit
import SwiftUI

struct ContactListView: View {
    @EnvironmentObject var model: OnboardingViewModel

    @State private var nonTalkieWalkieContacts: [ContactItem] = []
    @State private var talkiewalkieContacts: [ContactItem] = []
    @State var loading = false

    @State var phoneNumberKit = PhoneNumberKit()

    @AppStorage("hasRefusedSharingContactList") var hasRefusedSharingContactList: Bool = false

    var body: some View {
        ZStack {
            if loading {
                ProgressView()
            } else {
                VStack {
                    if !talkiewalkieContacts.isEmpty {
                        Text("You already have \(talkiewalkieContacts.count) friends")
                    }
                    Text("Invite your close ones!")
                        .padding()
                        .background(Color.white)
                        .foregroundColor(.yellow)
                        .rotationEffect(Angle(degrees: 15))

                    ScrollView {
                        LazyVStack(spacing: 40) {
                            ForEach(nonTalkieWalkieContacts) { contact in
                                HStack(spacing: 20) {
                                    AutomaticAvatar(String(contact.displayName.prefix(1)), color: generateColorFor(text: contact.phone))
                                    VStack(alignment: .leading, spacing: 10) {
                                        Text(contact.displayName).foregroundColor(.black)
                                        Text(contact.phone).foregroundColor(.gray).font(.caption)
                                    }
                                    Spacer()
                                    Button("whatsapp") {
                                        // TODO: invite logic
                                    }.foregroundColor(.green)
                                }
                                .frame(minWidth: 0, maxWidth: .infinity, alignment: .topLeading)
                                .padding()
                                .background(Color.white)
                                .cornerRadius(15)
                                .shadow(radius: 10)
                            }
                        }
                        .frame(minHeight: 0, maxHeight: .infinity)
                    }
                    if !hasRefusedSharingContactList {
                        OnboardingNavControls(page: $model.page)
                    }
                }
            }

            if hasRefusedSharingContactList {
                VStack {
                    Text("Here's how to share your contact list.")
                    Spacer()
                    Text("You can't proceed without sharing your contact list.")
                }
                .padding()
                .background(Color.white)
                .shadow(radius: 10)
            }
        }
        .padding()
        .onAppear {
            CNContactStore().requestAccess(for: .contacts) { access, _ in
                if access {
                    hasRefusedSharingContactList = false
                    let contactStore = CNContactStore()
                    let contactList = contactStore.allLocalPhoneNumbers()
                    // TODO: we really shouldn't start a new connection to the server here.
                    if let fbU = Auth.auth().currentUser {
                        loading = true
                        let persistentContainer = NSPersistentContainer(name: "LocalModels")

                        persistentContainer.loadPersistentStores { _, error in
                            persistentContainer.viewContext.automaticallyMergesChangesFromParent = true

                            if let error = error {
                                fatalError("Unable to load persistent stores: \(error)")
                            }
                        }

                        AuthenticatedState.build(Config.load(version: "dev"), fbU: fbU, context: persistentContainer.viewContext) { st in
                            let (twCL, _) = st.gApi.syncContactList(phones: contactList.map {
                                guard let phoneNumber = try? phoneNumberKit.parse(
                                    $0.phone,
                                    withRegion: PhoneNumberKit.defaultRegionCode()
                                ) else { return "" }
                                return self.phoneNumberKit.format(phoneNumber, toType: .e164)
                            })
                            loading = false
                            if let twCL = twCL {
                                talkiewalkieContacts = contactList.filter { twCL.users.map { u in u.phone }.contains($0.phone) }
                                nonTalkieWalkieContacts = contactList.filter { !twCL.users.map { u in u.phone }.contains($0.phone) }
                                persistentContainer.viewContext.saveOrLogError()
                            }
                            // disabling polling since it crashes once passed onboarding, likely due to conflicting core data instances
                            // necessary for tomorrow though
                            // DispatchQueue.main.asyncAfter(deadline: .now() + 1) { pollContactList(st: st, contactList: contactList, phoneNumberKit: self.phoneNumberKit)}
                        }
                    } else {
                        fatalError("unreachable state")
                    }
                } else {
                    hasRefusedSharingContactList = true
                }
            }
        }
    }
}

private func pollContactList(st: AuthenticatedState, contactList: [ContactItem], phoneNumberKit: PhoneNumberKit) {
    st.gApi.syncContactList(phones: contactList.map {
        guard let phoneNumber = try? phoneNumberKit.parse(
            $0.phone,
            withRegion: PhoneNumberKit.defaultRegionCode()
        ) else { return "" }
        return phoneNumberKit.format(phoneNumber, toType: .e164)
    })
    DispatchQueue.main.asyncAfter(deadline: .now() + 2) { pollContactList(st: st, contactList: contactList, phoneNumberKit: phoneNumberKit) }
}

struct ContactListView_Previews: PreviewProvider {
    static var previews: some View {
        ContactListView()
    }
}
