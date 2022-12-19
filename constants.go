package main

// put all the stuff that should remain in memory in one place

const BOT_VERSION = "3.0.0"

const BOT_ABOUT_INFO = `
'''prolog
Text Response Automation System
''''''md
<Version 3.0.0>
<Created_by Nathan Hedge>
''''''py
#################'''['''md
[Website](https://tras.almostd.one/)
'''](https://tras.almostd.one/)['''md
[Add Link](https://bit.ly/gotras)
'''](https://bit.ly/gotras)['''md
[Top.gg Page](https://top.gg/bot/494273862427738113)
'''](https://top.gg/bot/494273862427738113)['''md
[Git Repo](https://github.com/10Nates/tras)
'''](https://github.com/10Nates/tras)'''py
#################'''['''md
[Legal]()
[ ](TRAS operates under the MIT license)
[ ](https://github.com/10Nates/tras/LICENSE)
'''](https://github.com/10Nates/tras/LICENSE)
`

var superScriptReplacements = map[string]string{
	"0": "⁰", "1": "¹", "2": "²", "3": "³", "4": "⁴", "5": "⁵", "6": "⁶", "7": "⁷", // math
	"8": "⁸", "9": "⁹", "+": "⁺", "-": "⁻", "=": "⁼", "(": "⁽", ")": "⁾",
	"a": "ᵃ", "b": "ᵇ", "c": "ᶜ", "d": "ᵈ", "e": "ᵉ", "f": "ᶠ", "g": "ᵍ", "h": "ʰ", // lowercase
	"i": "ⁱ", "j": "ʲ", "k": "ᵏ", "l": "ˡ", "m": "ᵐ", "n": "ⁿ", "o": "ᵒ", "p": "ᵖ",
	"q": "ᵠ", "r": "ʳ", "s": "ˢ", "t": "ᵗ", "u": "ᵘ", "v": "ᵛ", "w": "ʷ", "x": "ˣ",
	"y": "ʸ", "z": "ᶻ",
	"A": "ᴬ", "B": "ᴮ", "C": "ᶜ", "D": "ᴰ", "E": "ᴱ", "F": "ᶠ", "G": "ᴳ", "H": "ᴴ", // uppercase
	"I": "ᴵ", "J": "ᴶ", "K": "ᴷ", "L": "ᴸ", "M": "ᴹ", "N": "ᴺ", "O": "ᴼ", "P": "ᴾ",
	"Q": "ᵠ", "R": "ᴿ", "S": "ˢ", "T": "ᵀ", "U": "ᵁ", "V": "ⱽ", "W": "ᵂ", "X": "ˣ",
	"Y": "ˠ", "Z": "ᶻ",
	"!": "ꜝ", ".": "ᐧ", "/": "ᐟ", "\\": "ᐠ", // extra characters
}

var HELP_COMMAND_RESPONSES = []string{ // various quips to make it friendlier
	"Here's your help hotline, hot and ready!",
	"Looking for lessons? You're in luck, here's a list!",
	"Confused and unsure? These commands will be your cure!",
	"Introducing this informative index!",
	"This list will lend a hand, just take a look and understand!",
	"Need some guidance? This directory's the key!",
	"This register has the answer, just take a look and you'll be a master!",
	"Looking for conclusions? This catalog has them all!",
	"Lost in a fog? These functions will clear the smog!",
}

var flagifyReplacements = map[string]string{
	"a": "🇦", "b": "🇧", "c": "🇨", "d": "🇩", "e": "🇪", "f": "🇫", "g": "🇬", "h": "🇭", "i": "🇮", "j": "🇯",
	"k": "🇰", "l": "🇱", "m": "🇲", "n": "🇳", "o": "🇴", "p": "🇵", "q": "🇶", "r": "🇷", "s": "🇸", "t": "🇹",
	"u": "🇺", "v": "🇻", "w": "🇼", "x": "🇽", "y": "🇾", "z": "🇿", // regional indicators, they combine into flags in fonts
}

var unicodifyReplacements = map[string]string{
	"0": "𝟢", "1": "𝟣", "2": "𝟤", "3": "𝟥", "4": "𝟦", "5": "𝟧", // numbers
	"6": "𝟨", "7": "𝟩", "8": "𝟪", "9": "𝟫",
	"a": "а", "b": "𝖻", "c": "с", "d": "𝖽", "e": "е", "f": "𝖿", // lowercase
	"g": "𝗀", "h": "𝗁", "i": "𝗂", "j": "𝗃", "k": "𝗄", "l": "ⅼ",
	"m": "𝗆", "n": "𝗇", "o": "о", "p": "р", "q": "𝗊", "r": "𝗋",
	"s": "ѕ", "t": "𝗍", "u": "𝗎", "v": "𝗏", "w": "𝗐", "x": "𝗑",
	"y": "𝗒", "z": "𝗓",
	"A": "А", "B": "В", "C": "С", "D": "𝖣", "E": "Е", "F": "𝖥", // uppercase
	"G": "𝖦", "H": "𝖧", "I": "І", "J": "𝖩", "K": "𝖪", "L": "𝖫",
	"M": "𝖬", "N": "𝖭", "O": "О", "P": "Р", "Q": "𝖰", "R": "𝖱",
	"S": "Ѕ", "T": "𝖳", "U": "𝖴", "V": "𝖵", "W": "𝖶", "X": "𝖷",
	"Y": "𝖸", "Z": "𝖹",
}

var emojifyReplacements = map[string]string{
	"cl": "🆑", "ab": "🆎", "ok": "🆗", "tm": "™️", // combo words first
	"10": "🔟", "0": "0️⃣", "1": "1️⃣", "2": "2️⃣",
	"3": "3️⃣", "4": "4️⃣", "5": "5️⃣", "6": "6️⃣",
	"7": "7️⃣", "8": "8️⃣", "9": "9️⃣",
	"a": "🅰", "b": "🅱️", "c": "©️", "h": "♓", // special lookalikes
	"i": "ℹ️", "m": "♏️", "o": "🅾️", "p": "🅿️",
	"r": "®️", "s": "💲", "t": "✝️", "x": "❌",
	" ": " ⬜ ", "!": "❗️", "?": "❓", "+": "➕", // misc lookalikes
	"-": "➖",
	"d": "🇩", "e": "🇪", "f": "🇫", "g": "🇬", // regional indicators as fallback
	"j": "🇯", "k": "🇰", "l": "🇱", "n": "🇳",
	"q": "🇶", "u": "🇺", "v": "🇻", "w": "🇼",
	"y": "🇾", "z": "🇿",
}

var boldReplacements = map[string]string{
	"A": "𝗔", "B": "𝗕", "C": "𝗖", "D": "𝗗", "E": "𝗘", "F": "𝗙", "G": "𝗚",
	"H": "𝗛", "I": "𝗜", "J": "𝗝", "K": "𝗞", "L": "𝗟", "M": "𝗠", "N": "𝗡",
	"O": "𝗢", "P": "𝗣", "Q": "𝗤", "R": "𝗥", "S": "𝗦", "T": "𝗧", "U": "𝗨",
	"V": "𝗩", "W": "𝗪", "X": "𝗫", "Y": "𝗬", "Z": "𝗭", "a": "𝗮", "b": "𝗯",
	"c": "𝗰", "d": "𝗱", "e": "𝗲", "f": "𝗳", "g": "𝗴", "h": "𝗵", "i": "𝗶",
	"j": "𝗷", "k": "𝗸", "l": "𝗹", "m": "𝗺", "n": "𝗻", "o": "𝗼", "p": "𝗽",
	"q": "𝗾", "r": "𝗿", "s": "𝘀", "t": "𝘁", "u": "𝘂", "v": "𝘃", "w": "𝘄",
	"x": "𝘅", "y": "𝘆", "z": "𝘇", "0": "𝟬", "1": "𝟭", "2": "𝟮", "3": "𝟯",
	"4": "𝟰", "5": "𝟱", "6": "𝟲", "7": "𝟳", "8": "𝟴", "9": "𝟵",
}

// I made this and then realized I didn't need it. Saving it because it was a lot of effort.
/*"ac": "🇦🇨", "ad": "🇦🇩", "ae": "🇦🇪", "af": "🇦🇫", "ag": "🇦🇬", "ai": "🇦🇮", "al": "🇦🇱", "am": "🇦🇲", "ao": "🇦🇴",
"aq": "🇦🇶", "ar": "🇦🇷", "as": "🇦🇸", "at": "🇦🇹", "au": "🇦🇺", "aw": "🇦🇼", "ax": "🇦🇽", "az": "🇦🇿", "ba": "🇧🇦",
"bb": "🇧🇧", "bd": "🇧🇩", "be": "🇧🇪", "bf": "🇧🇫", "bg": "🇧🇬", "bh": "🇧🇭", "bi": "🇧🇮", "bj": "🇧🇯", "bl": "🇧🇱",
"bm": "🇧🇲", "bn": "🇧🇳", "bo": "🇧🇴", "bq": "🇧🇶", "br": "🇧🇷", "bs": "🇧🇸", "bt": "🇧🇹", "bv": "🇧🇻", "bw": "🇧🇼",
"by": "🇧🇾", "bz": "🇧🇿", "ca": "🇨🇦", "cc": "🇨🇨", "cd": "🇨🇩", "cf": "🇨🇫", "cg": "🇨🇬", "ch": "🇨🇭", "ci": "🇨🇮",
"ck": "🇨🇰", "cl": "🇨🇱", "cm": "🇨🇲", "cn": "🇨🇳", "co": "🇨🇴", "cp": "🇨🇵", "cr": "🇨🇷", "cu": "🇨🇺", "cv": "🇨🇻",
"cw": "🇨🇼", "cx": "🇨🇽", "cy": "🇨🇾", "cz": "🇨🇿", "de": "🇩🇪", "dg": "🇩🇬", "dj": "🇩🇯", "dk": "🇩🇰", "dm": "🇩🇲",
"do": "🇩🇴", "dz": "🇩🇿", "ea": "🇪🇦", "ec": "🇪🇨", "ee": "🇪🇪", "eg": "🇪🇬", "eh": "🇪🇭", "er": "🇪🇷", "es": "🇪🇸",
"et": "🇪🇹", "eu": "🇪🇺", "fi": "🇫🇮", "fj": "🇫🇯", "fk": "🇫🇰", "fm": "🇫🇲", "fo": "🇫🇴", "fr": "🇫🇷", "ga": "🇬🇦",
"gb": "🇬🇧", "gd": "🇬🇩", "ge": "🇬🇪", "gf": "🇬🇫", "gg": "🇬🇬", "gh": "🇬🇭", "gi": "🇬🇮", "gl": "🇬🇱", "gm": "🇬🇲",
"gn": "🇬🇳", "gp": "🇬🇵", "gq": "🇬🇶", "gr": "🇬🇷", "gs": "🇬🇸", "gt": "🇬🇹", "gu": "🇬🇺", "gw": "🇬🇼", "gy": "🇬🇾",
"hk": "🇭🇰", "hm": "🇭🇲", "hn": "🇭🇳", "hr": "🇭🇷", "ht": "🇭🇹", "hu": "🇭🇺", "ic": "🇮🇨", "id": "🇮🇩", "ie": "🇮🇪",
"il": "🇮🇱", "im": "🇮🇲", "in": "🇮🇳", "io": "🇮🇴", "iq": "🇮🇶", "ir": "🇮🇷", "is": "🇮🇸", "it": "🇮🇹", "je": "🇯🇪",
"jm": "🇯🇲", "jo": "🇯🇴", "jp": "🇯🇵", "ke": "🇰🇪", "kg": "🇰🇬", "kh": "🇰🇭", "ki": "🇰🇮", "km": "🇰🇲", "kn": "🇰🇳",
"kp": "🇰🇵", "kr": "🇰🇷", "kw": "🇰🇼", "ky": "🇰🇾", "kz": "🇰🇿", "la": "🇱🇦", "lb": "🇱🇧", "lc": "🇱🇨", "li": "🇱🇮",
"lk": "🇱🇰", "lr": "🇱🇷", "ls": "🇱🇸", "lt": "🇱🇹", "lu": "🇱🇺", "lv": "🇱🇻", "ly": "🇱🇾", "ma": "🇲🇦", "mc": "🇲🇨",
"md": "🇲🇩", "me": "🇲🇪", "mf": "🇲🇫", "mg": "🇲🇬", "mh": "🇲🇭", "mk": "🇲🇰", "ml": "🇲🇱", "mm": "🇲🇲", "mn": "🇲🇳",
"mo": "🇲🇴", "mp": "🇲🇵", "mq": "🇲🇶", "mr": "🇲🇷", "ms": "🇲🇸", "mt": "🇲🇹", "mu": "🇲🇺", "mv": "🇲🇻", "mw": "🇲🇼",
"mx": "🇲🇽", "my": "🇲🇾", "mz": "🇲🇿", "na": "🇳🇦", "nc": "🇳🇨", "ne": "🇳🇪", "nf": "🇳🇫", "ng": "🇳🇬", "ni": "🇳🇮",
"nl": "🇳🇱", "no": "🇳🇴", "np": "🇳🇵", "nr": "🇳🇷", "nu": "🇳🇺", "nz": "🇳🇿", "om": "🇴🇲", "pa": "🇵🇦", "pe": "🇵🇪",
"pf": "🇵🇫", "pg": "🇵🇬", "ph": "🇵🇭", "pk": "🇵🇰", "pl": "🇵🇱", "pm": "🇵🇲", "pn": "🇵🇳", "pr": "🇵🇷", "ps": "🇵🇸",
"pt": "🇵🇹", "pw": "🇵🇼", "py": "🇵🇾", "qa": "🇶🇦", "re": "🇷🇪", "ro": "🇷🇴", "rs": "🇷🇸", "ru": "🇷🇺", "rw": "🇷🇼",
"sa": "🇸🇦", "sb": "🇸🇧", "sc": "🇸🇨", "sd": "🇸🇩", "se": "🇸🇪", "sg": "🇸🇬", "sh": "🇸🇭", "si": "🇸🇮", "sj": "🇸🇯",
"sk": "🇸🇰", "sl": "🇸🇱", "sm": "🇸🇲", "sn": "🇸🇳", "so": "🇸🇴", "sr": "🇸🇷", "ss": "🇸🇸", "st": "🇸🇹", "sv": "🇸🇻",
"sx": "🇸🇽", "sy": "🇸🇾", "sz": "🇸🇿", "ta": "🇹🇦", "tc": "🇹🇨", "td": "🇹🇩", "tf": "🇹🇫", "tg": "🇹🇬", "th": "🇹🇭",
"tj": "🇹🇯", "tk": "🇹🇰", "tl": "🇹🇱", "tm": "🇹🇲", "tn": "🇹🇳", "to": "🇹🇴", "tr": "🇹🇷", "tt": "🇹🇹", "tv": "🇹🇻",
"tw": "🇹🇼", "tz": "🇹🇿", "ua": "🇺🇦", "ug": "🇺🇬", "um": "🇺🇲", "un": "🇺🇳", "us": "🇺🇸", "uy": "🇺🇾", "uz": "🇺🇿",
"va": "🇻🇦", "vc": "🇻🇨", "ve": "🇻🇪", "vg": "🇻🇬", "vi": "🇻🇮", "vn": "🇻🇳", "vu": "🇻🇺", "wf": "🇼🇫", "ws": "🇼🇸",
"xk": "🇽🇰", "ye": "🇾🇪", "yt": "🇾🇹", "za": "🇿🇦", "zm": "🇿🇲", "zw": "🇿🇼",*/
