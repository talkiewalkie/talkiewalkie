import SwiftUI
import OSLog

extension AutoTextField.Representable {
    final class Coordinator: NSObject, UITextViewDelegate {

        internal let textView: UIKitTextView
        
        private var originalText: NSAttributedString = .init()
        private var text: Binding<NSAttributedString>
        private var isFocused: Binding<Bool>
        private var flexibleHeight: Binding<CGFloat>
        private var fixedLinesHeight: Binding<CGFloat>
        private var isScrollingEnabled: Binding<Bool>
        private var maxHeight: CGFloat

        var onCommit: (() -> Void)?
        var onEditingChanged: (() -> Void)?
        var shouldEditInRange: ((Range<String.Index>, String) -> Bool)?

        init(text: Binding<NSAttributedString>,
             isFocused: Binding<Bool>,
             flexibleHeight: Binding<CGFloat>,
             fixedLinesHeight: Binding<CGFloat>,
             isScrollingEnabled: Binding<Bool>,
             maxHeight: CGFloat,
             shouldEditInRange: ((Range<String.Index>, String) -> Bool)?,
             onEditingChanged: (() -> Void)?,
             onCommit: (() -> Void)?
        ) {
            textView = UIKitTextView()
            textView.backgroundColor = .clear
            textView.setContentCompressionResistancePriority(.defaultLow, for: .horizontal)
            
            self.text = text
            self.isFocused = isFocused
            self.flexibleHeight = flexibleHeight
            self.fixedLinesHeight = fixedLinesHeight
            self.isScrollingEnabled = isScrollingEnabled
            self.maxHeight = maxHeight
            self.shouldEditInRange = shouldEditInRange
            self.onEditingChanged = onEditingChanged
            self.onCommit = onCommit

            super.init()
            textView.delegate = self
        }

        func textViewDidBeginEditing(_ textView: UITextView) {
            DispatchQueue.main.async {
                self.isFocused.wrappedValue = true
            }
            
            originalText = text.wrappedValue
        }

        func textViewDidChange(_ textView: UITextView) {
            text.wrappedValue = NSAttributedString(attributedString: textView.attributedText)
            recalculateHeight(newMaximumNumberOfLines: textView.textContainer.maximumNumberOfLines)
            onEditingChanged?()
        }

        func textView(_ textView: UITextView, shouldChangeTextIn range: NSRange, replacementText text: String) -> Bool {
            if onCommit != nil, text == "\n" {
                onCommit?()
                originalText = NSAttributedString(attributedString: textView.attributedText)
                textView.resignFirstResponder()
                return false
            }

            return true
        }

        func textViewDidEndEditing(_ textView: UITextView) {
            DispatchQueue.main.async {
                self.isFocused.wrappedValue = false
            }
            
            // this check is to ensure we always commit text when we're not using a closure
            if onCommit != nil {
                text.wrappedValue = originalText
            }
        }

    }

}

extension AutoTextField.Representable.Coordinator {

    func update(representable: AutoTextField.Representable) {
        textView.attributedText = representable.text
        textView.font = representable.font
        textView.adjustsFontForContentSizeCategory = true
        textView.textColor = representable.foregroundColor
        textView.autocapitalizationType = representable.autocapitalization
        textView.autocorrectionType = representable.autocorrection
        textView.isEditable = representable.isEditable
        textView.isSelectable = representable.isSelectable
        textView.isScrollEnabled = representable.isScrollingEnabled
        textView.dataDetectorTypes = representable.autoDetectionTypes
        textView.allowsEditingTextAttributes = representable.allowsRichText

        switch representable.multilineTextAlignment {
        case .leading:
            textView.textAlignment = textView.traitCollection.layoutDirection ~= .leftToRight ? .left : .right
        case .trailing:
            textView.textAlignment = textView.traitCollection.layoutDirection ~= .leftToRight ? .right : .left
        case .center:
            textView.textAlignment = .center
        }

        if let value = representable.enablesReturnKeyAutomatically {
            textView.enablesReturnKeyAutomatically = value
        } else {
            textView.enablesReturnKeyAutomatically = onCommit == nil ? false : true
        }

        if let returnKeyType = representable.returnKeyType {
            textView.returnKeyType = returnKeyType
        } else {
            textView.returnKeyType = onCommit == nil ? .default : .done
        }

        if !representable.isScrollingEnabled {
            textView.textContainer.lineFragmentPadding = 0
            textView.textContainerInset = .zero
        }
        
        recalculateHeight(newMaximumNumberOfLines: representable.maximumNumberOfLines)
        
        textView.textContainer.maximumNumberOfLines = representable.maximumNumberOfLines
        if representable.maximumNumberOfLines > 0 {
            textView.textContainer.lineBreakMode = representable.truncationMode
        } else {
            textView.textContainer.lineBreakMode = .byWordWrapping
        }
        
        if representable.isFocused {
            textView.becomeFirstResponder()
        }
        
        textView.setNeedsDisplay()
    }

    private func recalculateHeight(newMaximumNumberOfLines: Int) {
        let currentMaximumNumberOfLines = textView.textContainer.maximumNumberOfLines
        let currentHeight = newMaximumNumberOfLines > 0 ? fixedLinesHeight : flexibleHeight
        
        let newSize = textView.sizeThatFits(CGSize(width: textView.frame.width, height: .greatestFiniteMagnitude))
        let enableScrolling: Bool
        
        if currentMaximumNumberOfLines == newMaximumNumberOfLines || textView.attributedText.string.isEmpty {
            enableScrolling = newSize.height > self.maxHeight
            
            if currentHeight.wrappedValue != newSize.height {
                DispatchQueue.main.async {
                    currentHeight.wrappedValue = newSize.height
                }
            }
        } else {
            if newMaximumNumberOfLines > 0 {
                enableScrolling = false
            } else {
                enableScrolling = currentHeight.wrappedValue > self.maxHeight
            }
        }
        
        if isScrollingEnabled.wrappedValue != enableScrolling {
            DispatchQueue.main.async {
                self.isScrollingEnabled.wrappedValue = enableScrolling
            }
        }
    }
}


extension UITextView {
    func numberOfLines() -> Int {
        let layoutManager = self.layoutManager
        let numberOfGlyphs = layoutManager.numberOfGlyphs
        var lineRange: NSRange = NSMakeRange(0, 1)
        var index = 0
        var numberOfLines = 0

        while index < numberOfGlyphs {
            layoutManager.lineFragmentRect(
                forGlyphAt: index, effectiveRange: &lineRange
            )
            index = NSMaxRange(lineRange)
            numberOfLines += 1
        }
        return numberOfLines
    }
}
