import SwiftUI

public extension UIFont {
    static var caption2: UIFont = .preferredFont(forTextStyle: .caption2)
    static var caption: UIFont = .preferredFont(forTextStyle: .caption1)
    static var footnote: UIFont = .preferredFont(forTextStyle: .footnote)
    static var callout: UIFont = .preferredFont(forTextStyle: .callout)
    static var body: UIFont = .preferredFont(forTextStyle: .body)
    static var subheadline: UIFont = .preferredFont(forTextStyle: .subheadline)
    static var headline: UIFont = .preferredFont(forTextStyle: .headline)
    static var title3: UIFont = .preferredFont(forTextStyle: .title3)
    static var title2: UIFont = .preferredFont(forTextStyle: .title2)
    static var title: UIFont = .preferredFont(forTextStyle: .title1)
    static var largeTitle: UIFont = .preferredFont(forTextStyle: .largeTitle)
}

public extension UIFont {
    enum Leading {
        case loose
        case tight
    }

    private func addingAttributes(_ attributes: [UIFontDescriptor.AttributeName: Any]) -> UIFont {
        return UIFont(descriptor: fontDescriptor.addingAttributes(attributes), size: pointSize)
    }

    static func system(size: CGFloat, weight: UIFont.Weight, design: UIFontDescriptor.SystemDesign = .default) -> UIFont {
        let descriptor = UIFont.systemFont(ofSize: size).fontDescriptor
            .addingAttributes([
                UIFontDescriptor.AttributeName.traits: [
                    UIFontDescriptor.TraitKey.weight: weight.rawValue,
                ],
            ]).withDesign(design)!
        return UIFont(descriptor: descriptor, size: size)
    }

    static func system(_ style: UIFont.TextStyle, design: UIFontDescriptor.SystemDesign = .default) -> UIFont {
        let descriptor = UIFontDescriptor.preferredFontDescriptor(withTextStyle: style).withDesign(design)!
        return UIFont(descriptor: descriptor, size: 0)
    }

    func weight(_ weight: UIFont.Weight) -> UIFont {
        return addingAttributes([
            UIFontDescriptor.AttributeName.traits: [
                UIFontDescriptor.TraitKey.weight: weight.rawValue,
            ],
        ])
    }

    func italic() -> UIFont {
        let descriptor = fontDescriptor.withSymbolicTraits(.traitItalic)!
        return UIFont(descriptor: descriptor, size: 0)
    }

    func bold() -> UIFont {
        let descriptor = fontDescriptor.withSymbolicTraits(.traitBold)!
        return UIFont(descriptor: descriptor, size: 0)
    }

    func leading(_ leading: Leading) -> UIFont {
        let descriptor = fontDescriptor.withSymbolicTraits(leading == .loose ? .traitLooseLeading : .traitTightLeading)!
        return UIFont(descriptor: descriptor, size: 0)
    }

    func smallCaps() -> UIFont {
        return addingAttributes([
            .featureSettings: [
                [
                    UIFontDescriptor.FeatureKey.featureIdentifier: kLowerCaseType,
                    UIFontDescriptor.FeatureKey.typeIdentifier: kLowerCaseSmallCapsSelector,
                ],
                [
                    UIFontDescriptor.FeatureKey.featureIdentifier: kUpperCaseType,
                    UIFontDescriptor.FeatureKey.typeIdentifier: kUpperCaseSmallCapsSelector,
                ],
            ],
        ])
    }

    func lowercaseSmallCaps() -> UIFont {
        return addingAttributes([
            .featureSettings: [
                [
                    UIFontDescriptor.FeatureKey.featureIdentifier: kLowerCaseType,
                    UIFontDescriptor.FeatureKey.typeIdentifier: kLowerCaseSmallCapsSelector,
                ],
            ],
        ])
    }

    func uppercaseSmallCaps() -> UIFont {
        return addingAttributes([
            .featureSettings: [
                [
                    UIFontDescriptor.FeatureKey.featureIdentifier: kUpperCaseType,
                    UIFontDescriptor.FeatureKey.typeIdentifier: kUpperCaseSmallCapsSelector,
                ],
            ],
        ])
    }

    func monospacedDigit() -> UIFont {
        return addingAttributes([
            .featureSettings: [
                [
                    UIFontDescriptor.FeatureKey.featureIdentifier: kNumberSpacingType,
                    UIFontDescriptor.FeatureKey.typeIdentifier: kMonospacedNumbersSelector,
                ],
            ],
        ])
    }
}
