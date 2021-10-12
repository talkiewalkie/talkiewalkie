//
//  Date+Format.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import Foundation

extension Date {
    func getFormattedDate(format: String) -> String {
        let dateformat = DateFormatter()
        dateformat.dateFormat = format
        return dateformat.string(from: self)
    }
}
