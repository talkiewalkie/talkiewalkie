//
//  String+NSRange.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 07.10.21.
//

import Foundation

extension String {
    func substring(with nsrange: NSRange) -> String {
        guard let range = Range(nsrange, in: self) else { return "" }
        return String(self[range])
    }
}
