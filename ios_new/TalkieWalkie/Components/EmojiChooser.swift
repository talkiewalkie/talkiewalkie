//
//  EmojiChooser.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import SwiftUI
import SDWebImageSwiftUI

class EmojiStore {
    let smileysPeople = "😀😬😁😂😃😄😅😆😇😉😊🙂🙃☺😋😌😍😘😗😙😚😜😝😛🤑🤓😎🤗😏😶😐😑😒🙄🤔😳😞😟😠😡😔😕🙁☹😣😖😫😩😤😮😱😨😰😯😦😧😢😥😪😓😭😵😲🤐😷🤒🤕😴💤💩😈👿👹👺💀👻👽🤖😺😸😹😻😼😽🙀😿😾🙌👏👋👍👎👊✊✌👌✋👐💪🙏☝👆👇👈👉🖕🖐🤘🖖✍💅👄👅👂👃👁👀👤👥🗣👶👦👧👨👩👱👴👵👲👳👮👷💂🕵🎅👼👸👰🚶🏃💃👯👫👬👭🙇💁🙅🙆🙋🙎🙍💇💆💑👩‍❤️‍👩👨‍❤️‍👨💏👩‍❤️‍💋‍👩👨‍❤️‍💋‍👨👪👨‍👩‍👧👨‍👩‍👧‍👦👨‍👩‍👦‍👦👨‍👩‍👧‍👧👩‍👩‍👦👩‍👩‍👧👩‍👩‍👧‍👦👩‍👩‍👦‍👦👩‍👩‍👧‍👧👨‍👨‍👦👨‍👨‍👧👨‍👨‍👧‍👦👨‍👨‍👦‍👦👨‍👨‍👧‍👧👚👕👖👔👗👙👘💄💋👣👠👡👢👞👟👒🎩🎓👑⛑🎒👝👛👜💼👓🕶💍🌂"
    
    let animalsNature = "🐶🐱🐭🐹🐰🐻🐼🐨🐯🦁🐮🐷🐽🐸🐙🐵🙈🙉🙊🐒🐔🐧🐦🐤🐣🐥🐺🐗🐴🦄🐝🐛🐌🐞🐜🕷🦂🦀🐍🐢🐠🐟🐡🐬🐳🐋🐊🐆🐅🐃🐂🐄🐪🐫🐘🐐🐏🐑🐎🐖🐀🐁🐓🦃🕊🐕🐩🐈🐇🐿🐾🐉🐲🌵🎄🌲🌳🌴🌱🌿☘🍀🎍🎋🍃🍂🍁🌾🌺🌻🌹🌷🌼🌸💐🍄🌰🎃🐚🕸🌎🌍🌏🌕🌖🌗🌘🌑🌒🌓🌔🌚🌝🌛🌜🌞🌙⭐🌟💫✨☄☀🌤⛅🌥🌦☁🌧⛈🌩⚡🔥💥❄🌨☃⛄🌬💨🌪🌫☂☔💧💦🌊"
    
    let foodDrink = "🍏🍎🍐🍊🍋🍌🍉🍇🍓🍈🍒🍑🍍🍅🍆🌶🌽🍠🍯🍞🧀🍗🍖🍤🍳🍔🍟🌭🍕🍝🌮🌯🍜🍲🍥🍣🍱🍛🍙🍚🍘🍢🍡🍧🍨🍦🍰🎂🍮🍬🍭🍫🍿🍩🍪🍺🍻🍷🍸🍹🍾🍶🍵☕🍼🍴🍽"
    
    let activity = "⚽🏀🏈⚾🎾🏐🏉🎱⛳🏌🏓🏸🏒🏑🏏🎿⛷🏂⛸🏹🎣🚣🏊🏄🛀⛹🏋🚴🚵🏇🕴🏆🎽🏅🎖🎗🏵🎫🎟🎭🎨🎪🎤🎧🎼🎹🎷🎺🎸🎻🎬🎮👾🎯🎲🎰🎳"
    
    let travelPlaces = "🚗🚕🚙🚌🚎🏎🚓🚑🚒🚐🚚🚛🚜🏍🚲🚨🚔🚍🚘🚖🚡🚠🚟🚃🚋🚝🚄🚅🚈🚞🚂🚆🚇🚊🚉🚁🛩✈🛫🛬⛵🛥🚤⛴🛳🚀🛰💺⚓🚧⛽🚏🚦🚥🏁🚢🎡🎢🎠🏗🌁🗼🏭⛲🎑⛰🏔🗻🌋🗾🏕⛺🏞🛣🛤🌅🌄🏜🏖🏝🌇🌆🏙🌃🌉🌌🌠🎇🎆🌈🏘🏰🏯🏟🗽🏠🏡🏚🏢🏬🏣🏤🏥🏦🏨🏪🏫🏩💒🏛⛪🕌🕍🕋⛩"
    
    let objects = "⌚📱📲💻⌨🖥🖨🖱🖲🕹🗜💽💾💿📀📼📷📸📹🎥📽🎞📞☎📟📠📺📻🎙🎚🎛⏱⏲⏰🕰⏳⌛📡🔋🔌💡🔦🕯🗑🛢💸💵💴💶💷💰💳💎⚖🔧🔨⚒🛠⛏🔩⚙⛓🔫💣🔪🗡⚔🛡🚬☠⚰⚱🏺🔮📿💈⚗🔭🔬🕳💊💉🌡🏷🔖🚽🚿🛁🔑🗝🛋🛌🛏🚪🛎🖼🗺⛱🗿🛍🎈🎏🎀🎁🎊🎉🎎🎐🎌🏮✉📩📨📧💌📮📪📫📬📭📦📯📥📤📜📃📑📊📈📉📄📅📆🗓📇🗃🗳🗄📋🗒📁📂🗂🗞📰📓📕📗📘📙📔📒📚📖🔗📎🖇✂📐📏📌📍🚩🏳🏴🔐🔒🔓🔏🖊🖋✒📝✏🖍🖌🔍🔎"
    
    let symbols = "❤💛💚💙💜💔❣💕💞💓💗💖💘💝💟☮✝☪🕉☸✡🔯🕎☯☦🛐⛎♈♉♊♋♌♍♎♏♐♑♒♓🆔⚛🈳🈹☢☣📴📳🈶🈚🈸🈺🈷✴🆚🉑💮🉐㊙㊗🈴🈵🈲🅰🅱🆎🆑🅾🆘⛔📛🚫❌⭕💢♨🚷🚯🚳🚱🔞📵❗❕❓❔‼⁉💯🔅🔆🔱⚜〽⚠🚸🔰♻🈯💹❇✳❎✅💠🌀➿🌐Ⓜ🏧🈂🛂🛃🛄🛅♿🚭🚾🅿🚰🚹🚺🚼🚻🚮🎦📶🈁🆖🆗🆙🆒🆕🆓0⃣1⃣2⃣3⃣4⃣5⃣6⃣7⃣8⃣9⃣🔟🔢▶⏸⏯⏹⏺⏭⏮⏩⏪🔀🔁🔂◀🔼🔽⏫⏬➡⬅⬆⬇↗↘↙↖↕↔🔄↪↩⤴⤵#⃣*⃣ℹ🔤🔡🔠🔣🎵🎶〰➰✔🔃➕➖➗✖💲💱🔚🔙🔛🔝🔜☑🔘⚪⚫🔴🔵🔸🔹🔶🔷🔺▪▫⬛⬜🔻◼◻◾◽🔲🔳🔈🔉🔊🔇📣📢🔔🔕🃏🀄♠♣♥♦🎴🗨💭🗯💬🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚🕛🕜🕝🕞🕟🕠🕡🕢🕣🕤🕥🕦🕧"
    
    let flags = "🇦🇫🇦🇽🇦🇱🇩🇿🇦🇸🇦🇩🇦🇴🇦🇮🇦🇶🇦🇬🇦🇷🇦🇲🇦🇼🇦🇺🇦🇹🇦🇿🇧🇸🇧🇭🇧🇩🇧🇧🇧🇾🇧🇪🇧🇿🇧🇯🇧🇲🇧🇹🇧🇴🇧🇶🇧🇦🇧🇼🇧🇷🇮🇴🇻🇬🇧🇳🇧🇬🇧🇫🇧🇮🇨🇻🇰🇭🇨🇲🇨🇦🇮🇨🇰🇾🇨🇫🇹🇩🇨🇱🇨🇳🇨🇽🇨🇨🇨🇴🇰🇲🇨🇬🇨🇩🇨🇰🇨🇷🇭🇷🇨🇺🇨🇼🇨🇾🇨🇿🇩🇰🇩🇯🇩🇲🇩🇴🇪🇨🇪🇬🇸🇻🇬🇶🇪🇷🇪🇪🇪🇹🇪🇺🇫🇰🇫🇴🇫🇯🇫🇮🇫🇷🇬🇫🇵🇫🇹🇫🇬🇦🇬🇲🇬🇪🇩🇪🇬🇭🇬🇮🇬🇷🇬🇱🇬🇩🇬🇵🇬🇺🇬🇹🇬🇬🇬🇳🇬🇼🇬🇾🇭🇹🇭🇳🇭🇰🇭🇺🇮🇸🇮🇳🇮🇩🇮🇷🇮🇶🇮🇪🇮🇲🇮🇱🇮🇹🇨🇮🇯🇲🇯🇵🇯🇪🇯🇴🇰🇿🇰🇪🇰🇮🇽🇰🇰🇼🇰🇬🇱🇦🇱🇻🇱🇧🇱🇸🇱🇷🇱🇾🇱🇮🇱🇹🇱🇺🇲🇴🇲🇰🇲🇬🇲🇼🇲🇾🇲🇻🇲🇱🇲🇹🇲🇭🇲🇶🇲🇷🇲🇺🇾🇹🇲🇽🇫🇲🇲🇩🇲🇨🇲🇳🇲🇪🇲🇸🇲🇦🇲🇿🇲🇲🇳🇦🇳🇷🇳🇵🇳🇱🇳🇨🇳🇿🇳🇮🇳🇪🇳🇬🇳🇺🇳🇫🇲🇵🇰🇵🇳🇴🇴🇲🇵🇰🇵🇼🇵🇸🇵🇦🇵🇬🇵🇾🇵🇪🇵🇭🇵🇳🇵🇱🇵🇹🇵🇷🇶🇦🇷🇪🇷🇴🇷🇺🇷🇼🇧🇱🇸🇭🇰🇳🇱🇨🇵🇲🇻🇨🇼🇸🇸🇲🇸🇹🇸🇦🇸🇳🇷🇸🇸🇨🇸🇱🇸🇬🇸🇽🇸🇰🇸🇮🇸🇧🇸🇴🇿🇦🇬🇸🇰🇷🇸🇸🇪🇸🇱🇰🇸🇩🇸🇷🇸🇿🇸🇪🇨🇭🇸🇾🇹🇼🇹🇯🇹🇿🇹🇭🇹🇱🇹🇬🇹🇰🇹🇴🇹🇹🇹🇳🇹🇷🇹🇲🇹🇨🇹🇻🇺🇬🇺🇦🇦🇪🇬🇧🇺🇸🇻🇮🇺🇾🇺🇿🇻🇺🇻🇦🇻🇪🇻🇳🇼🇫🇪🇭🇾🇪🇿🇲🇿🇼"
    
    var allEmojis: [Character] {
        return Array(smileysPeople + animalsNature + foodDrink + activity + travelPlaces + objects + symbols + flags)
    }
}

struct EmojiChooser: View {
    @Binding var show: Bool
    
    var onChoose: ((UIImage) -> Void)?
    
    init(show: Binding<Bool>, onChoose: ((UIImage) -> Void)? = nil) {
        self._show = show
        self.onChoose = onChoose
        
        UINavigationBar.appearance().barTintColor = .systemBackground
        UINavigationBar.appearance().shadowImage = UIImage()
    }
    
    var body: some View {
        
        NavigationView {
            EmojiSearchView(show: $show, onChoose: onChoose)
                .navigationBarTitle(LocalizedStringKey("Emojis"), displayMode: .inline)
                .navigationBarItems(trailing: CloseButton(show: $show))
        }
        
    }
}

struct EmojiSearchView: View {
    @Binding var show: Bool
    
    var onChoose: ((UIImage) -> Void)?
    
    @State private var emojiStore = EmojiStore()
    
    func emojiToId(_ emoji: Character) -> String {
        let uni = emoji.unicodeScalars
        let unicodes = uni.map { $0.value }.map { String($0, radix: 16, uppercase: false) }
        return unicodes.joined(separator: "-")
    }

    func tweemojiURL(for emoji: Character) -> URL? {
        let baseURL = URL(string: "https://twemoji.maxcdn.com/v/latest/72x72/")
        let url = baseURL?.appendingPathComponent(emojiToId(emoji)).appendingPathExtension("png")
        return url
    }
    
    func handleTap(emoji: Character) {
        SDWebImageManager.shared.loadImage(with: tweemojiURL(for: emoji), progress: nil, completed: { image, _, _, _, _, _ in
            guard let image = image else { return }
            show = false
            onChoose?(image)
        })
    }
    
    var body: some View {
        ScrollView(.vertical, showsIndicators: false) {
            VStack(alignment: .leading) {
                Text("Choose an Emoji")
                    .font(.title2)
                    .fontWeight(.medium)
                    .padding()
                
                LazyVGrid(columns: [GridItem(.adaptive(minimum: 50), spacing: 20)], spacing: 15) {
                    ForEach(emojiStore.allEmojis, id: \.self) { emoji in
                        Button(action: {
                            handleTap(emoji: emoji)
                        }) {
                            WebImage(url: tweemojiURL(for: emoji))
                                .resizable()
                                .placeholder {
                                    Circle()
                                        .stroke(style: StrokeStyle(lineCap: .round, lineJoin: .round, dash: [5]))
                                        .foregroundColor(.gray)
                                }
                                .scaledToFit()
                                .frame(width: 50, height: 50, alignment: .center)
                        }
                        .buttonStyle(PlainButtonStyle())
                        
                    }
                }
            }
            .padding()
        }
    }
}

struct EmojiChooser_Previews: PreviewProvider {
    static var previews: some View {
        EmojiChooser(show: .constant(true))
    }
}
