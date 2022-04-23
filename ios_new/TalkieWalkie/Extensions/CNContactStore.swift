//
//  CNContactStore.swift
//  TalkieWalkie
//
//  Created by Théo Matussière on 17/10/2021.
//

import Contacts
import Foundation
import OSLog

struct ContactItem: Identifiable {
    var id: String { return phone }

    var displayName: String
    var label: String
    var phone: String
}

extension CNContactStore {
    // https://stackoverflow.com/a/49129459
    func allLocalPhoneNumbers() -> [ContactItem] {
        var contacts = [ContactItem]()
        let keys = [
            CNContactFormatter.descriptorForRequiredKeys(for: .fullName),
            CNContactPhoneNumbersKey,
            CNContactEmailAddressesKey,
        ] as [Any]
        let request = CNContactFetchRequest(keysToFetch: keys as! [CNKeyDescriptor])
        do {
            try enumerateContacts(with: request) { contact, _ in
                for phoneNumber in contact.phoneNumbers {
                    if let number = phoneNumber.value as? CNPhoneNumber, let label = phoneNumber.label {
                        let localizedLabel = CNLabeledValue<CNPhoneNumber>.localizedString(forLabel: label)
                        contacts.append(ContactItem(displayName: "\(contact.givenName) \(contact.familyName)", label: localizedLabel, phone: number.stringValue))
                    }
                }
            }
            return contacts
        } catch {
            os_log(.error, "unable to fetch contacts: \(error.localizedDescription)")
            return []
        }
    }
}
