import SwiftUI

/// A SwiftUI TextView implementation that supports both scrolling and auto-sizing layouts
public struct AutoTextField: View {
    @Environment(\.layoutDirection) private var layoutDirection

    @Binding private var text: NSAttributedString
    @Binding private var isEmpty: Bool
    @Binding private var isFocused: Bool

    @State private var fixedLinesHeight: CGFloat = 22
    @State private var flexibleHeight: CGFloat = 44

    @State var isScrollingEnabled: Bool = true

    private var onEditingChanged: (() -> Void)?
    private var shouldEditInRange: ((Range<String.Index>, String) -> Bool)?
    private var onCommit: (() -> Void)?

    var placeholderView: AnyView?
    var foregroundColor: UIColor = .label
    var autocapitalization: UITextAutocapitalizationType = .sentences
    var multilineTextAlignment: TextAlignment = .leading
    var font: UIFont = .preferredFont(forTextStyle: .body)
    var returnKeyType: UIReturnKeyType?
    var clearsOnInsertion: Bool = false
    var autocorrection: UITextAutocorrectionType = .default
    var truncationMode: NSLineBreakMode = .byTruncatingTail
    var maximumNumberOfLines: Int = 0
    var isEditable: Bool = true
    var isSelectable: Bool = true
    var enablesReturnKeyAutomatically: Bool?
    var autoDetectionTypes: UIDataDetectorTypes = []
    var allowRichText: Bool
    var maxHeight: CGFloat = .infinity

    /// Makes a new TextView with the specified configuration
    /// - Parameters:
    ///   - text: A binding to the text
    ///   - shouldEditInRange: A closure that's called before an edit it applied, allowing the consumer to prevent the change
    ///   - onEditingChanged: A closure that's called after an edit has been applied
    ///   - onCommit: If this is provided, the field will automatically lose focus when the return key is pressed
    public init(_ text: Binding<String>,
                isFocused: Binding<Bool>,
                shouldEditInRange: ((Range<String.Index>, String) -> Bool)? = nil,
                onEditingChanged: (() -> Void)? = nil,
                onCommit: (() -> Void)? = nil)
    {
        _text = Binding(
            get: { NSAttributedString(string: text.wrappedValue) },
            set: { text.wrappedValue = $0.string }
        )

        _isEmpty = Binding(
            get: { text.wrappedValue.isEmpty },
            set: { _ in }
        )

        _isFocused = isFocused
        self.onCommit = onCommit
        self.shouldEditInRange = shouldEditInRange
        self.onEditingChanged = onEditingChanged

        allowRichText = false
    }

    /// Makes a new TextView that supports `NSAttributedString`
    /// - Parameters:
    ///   - text: A binding to the attributed text
    ///   - onEditingChanged: A closure that's called after an edit has been applied
    ///   - onCommit: If this is provided, the field will automatically lose focus when the return key is pressed
    public init(_ text: Binding<NSAttributedString>,
                isFocused: Binding<Bool>,
                onEditingChanged: (() -> Void)? = nil,
                onCommit: (() -> Void)? = nil)
    {
        _text = text
        _isEmpty = Binding(
            get: { text.wrappedValue.string.isEmpty },
            set: { _ in }
        )

        _isFocused = isFocused
        self.onCommit = onCommit
        self.onEditingChanged = onEditingChanged

        allowRichText = true
    }

    public var body: some View {
        let height = min(maximumNumberOfLines > 0 ? fixedLinesHeight : flexibleHeight, maxHeight)

        Representable(
            text: $text,
            isFocused: $isFocused,
            flexibleHeight: $flexibleHeight,
            fixedLinesHeight: $fixedLinesHeight,
            isScrollingEnabled: $isScrollingEnabled,
            foregroundColor: foregroundColor,
            autocapitalization: autocapitalization,
            multilineTextAlignment: multilineTextAlignment,
            font: font,
            returnKeyType: returnKeyType,
            clearsOnInsertion: clearsOnInsertion,
            autocorrection: autocorrection,
            truncationMode: truncationMode,
            maximumNumberOfLines: maximumNumberOfLines,
            isEditable: isEditable,
            isSelectable: isSelectable,
            enablesReturnKeyAutomatically: enablesReturnKeyAutomatically,
            autoDetectionTypes: autoDetectionTypes,
            allowsRichText: allowRichText,
            maxHeight: maxHeight,
            onEditingChanged: onEditingChanged,
            shouldEditInRange: shouldEditInRange,
            onCommit: onCommit
        )
        .frame(minHeight: height, maxHeight: height)
        .background(
            placeholderView?
                .foregroundColor(Color(.placeholderText))
                .multilineTextAlignment(multilineTextAlignment)
                .font(Font(font))
                .padding(.horizontal, isScrollingEnabled ? 5 : 0)
                .padding(.vertical, isScrollingEnabled ? 8 : 0)
                .opacity(isEmpty ? 1 : 0),
            alignment: .topLeading
        )
    }
}

final class UIKitTextView: UITextView {
    override var keyCommands: [UIKeyCommand]? {
        return (super.keyCommands ?? []) + [
            UIKeyCommand(input: UIKeyCommand.inputEscape, modifierFlags: [], action: #selector(escape(_:))),
        ]
    }

    @objc private func escape(_: Any) {
        resignFirstResponder()
    }
}
