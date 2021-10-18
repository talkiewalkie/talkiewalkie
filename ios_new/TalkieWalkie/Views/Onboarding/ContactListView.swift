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
    @EnvironmentObject var authState: AuthState

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
                    if case .Connected(let api, _) = authState.state {
                        loading = true
                        let (twCL, _) = api.syncContactList(phones: contactList.map {
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
                            authState.moc.saveOrLogError()
                        }
                    }
                } else {
                    hasRefusedSharingContactList = true
                }
            }
        }
    }
}

struct ContactListView_Previews: PreviewProvider {
    static var previews: some View {
        ContactListView()
    }
}
