import OSLog
import SwiftUI

extension AutoTextField {
    struct Representable: UIViewRepresentable {
        @Binding var text: NSAttributedString
        @Binding var isFocused: Bool
        @Binding var flexibleHeight: CGFloat
        @Binding var fixedLinesHeight: CGFloat
        @Binding var isScrollingEnabled: Bool

        let foregroundColor: UIColor
        let autocapitalization: UITextAutocapitalizationType
        var multilineTextAlignment: TextAlignment
        let font: UIFont
        let returnKeyType: UIReturnKeyType?
        let clearsOnInsertion: Bool
        let autocorrection: UITextAutocorrectionType
        let truncationMode: NSLineBreakMode
        let maximumNumberOfLines: Int
        let isEditable: Bool
        let isSelectable: Bool
        let enablesReturnKeyAutomatically: Bool?
        var autoDetectionTypes: UIDataDetectorTypes = []
        var allowsRichText: Bool
        var maxHeight: CGFloat

        var onEditingChanged: (() -> Void)?
        var shouldEditInRange: ((Range<String.Index>, String) -> Bool)?
        var onCommit: (() -> Void)?

        func makeUIView(context: Context) -> UIKitTextView {
            context.coordinator.textView
        }

        func updateUIView(_: UIKitTextView, context: Context) {
            context.coordinator.update(representable: self)
        }

        @discardableResult func makeCoordinator() -> Coordinator {
            Coordinator(
                text: $text,
                isFocused: $isFocused,
                flexibleHeight: $flexibleHeight,
                fixedLinesHeight: $fixedLinesHeight,
                isScrollingEnabled: $isScrollingEnabled,
                maxHeight: maxHeight,
                shouldEditInRange: shouldEditInRange,
                onEditingChanged: onEditingChanged,
                onCommit: onCommit
            )
        }
    }
}
