//
//  RecordSheetView.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 20.10.21.
//

import SwiftUI

struct RecordSheetView: View {
    @State var isRecording: Bool = false
    
    var body: some View {
        VStack {
            RecordButton(isRecording: $isRecording)
        }
    }
}

struct RecordSheetView_Previews: PreviewProvider {
    static var previews: some View {
        
        ZStack {
            List(1..<20) { i in
                Text("\(i)")
            }
            
            RecordSheetView()
        }
        
    }
}
