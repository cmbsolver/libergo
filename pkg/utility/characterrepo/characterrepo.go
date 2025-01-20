package runelib

import (
	"strings"
)

type CharacterRepo struct {
	asciiAnsiItems map[string]map[int]string
}

func NewCharacterRepo() *CharacterRepo {
	repo := &CharacterRepo{
		asciiAnsiItems: make(map[string]map[int]string),
	}

	repo.asciiAnsiItems["ASCII"] = map[int]string{
		0: "<NUL>", 1: "<SOH>", 2: "<STX>", 3: "<ETX>", 4: "<EOT>", 5: "<ENQ>", 6: "<ACK>", 7: "<BEL>",
		8: "<BS>", 9: "<HT>", 10: "\n", 11: "<VT>", 12: "<FF>", 13: "\r", 14: "<SO>", 15: "<SI>",
		16: "<DLE>", 17: "<DC1>", 18: "<DC2>", 19: "<DC3>", 20: "<DC4>", 21: "<NAK>", 22: "<SYN>", 23: "<ETB>",
		24: "<CAN>", 25: "<EM>", 26: "<SUB>", 27: "<ESC>", 28: "<FS>", 29: "<GS>", 30: "<RS>", 31: "<US>",
		32: " ", 33: "!", 34: "\"", 35: "#", 36: "$", 37: "%", 38: "&", 39: "'",
		40: "(", 41: ")", 42: "*", 43: "+", 44: ",", 45: "-", 46: ".", 47: "/",
		48: "0", 49: "1", 50: "2", 51: "3", 52: "4", 53: "5", 54: "6", 55: "7",
		56: "8", 57: "9", 58: ":", 59: ";", 60: "<", 61: "=", 62: ">", 63: "?",
		64: "@", 65: "A", 66: "B", 67: "C", 68: "D", 69: "E", 70: "F", 71: "G",
		72: "H", 73: "I", 74: "J", 75: "K", 76: "L", 77: "M", 78: "N", 79: "O",
		80: "P", 81: "Q", 82: "R", 83: "S", 84: "T", 85: "U", 86: "V", 87: "W",
		88: "X", 89: "Y", 90: "Z", 91: "[", 92: "\\", 93: "]", 94: "^", 95: "_",
		96: "`", 97: "a", 98: "b", 99: "c", 100: "d", 101: "e", 102: "f", 103: "g",
		104: "h", 105: "i", 106: "j", 107: "k", 108: "l", 109: "m", 110: "n", 111: "o",
		112: "p", 113: "q", 114: "r", 115: "s", 116: "t", 117: "u", 118: "v", 119: "w",
		120: "x", 121: "y", 122: "z", 123: "{", 124: "|", 125: "}", 126: "~", 127: "<DEL>",
	}

	repo.asciiAnsiItems["ANSI"] = map[int]string{
		0: "<NUL>", 1: "<SOH>", 2: "<STX>", 3: "<ETX>", 4: "<EOT>", 5: "<ENQ>", 6: "<ACK>", 7: "<BEL>",
		8: "<BS>", 9: "<HT>", 10: "\n", 11: "<VT>", 12: "<FF>", 13: "\r", 14: "<SO>", 15: "<SI>",
		16: "<DLE>", 17: "<DC1>", 18: "<DC2>", 19: "<DC3>", 20: "<DC4>", 21: "<NAK>", 22: "<SYN>", 23: "<ETB>",
		24: "<CAN>", 25: "<EM>", 26: "<SUB>", 27: "<ESC>", 28: "<FS>", 29: "<GS>", 30: "<RS>", 31: "<US>",
		32: " ", 33: "!", 34: "\"", 35: "#", 36: "$", 37: "%", 38: "&", 39: "'",
		40: "(", 41: ")", 42: "*", 43: "+", 44: ",", 45: "-", 46: ".", 47: "/",
		48: "0", 49: "1", 50: "2", 51: "3", 52: "4", 53: "5", 54: "6", 55: "7",
		56: "8", 57: "9", 58: ":", 59: ";", 60: "<", 61: "=", 62: ">", 63: "?",
		64: "@", 65: "A", 66: "B", 67: "C", 68: "D", 69: "E", 70: "F", 71: "G",
		72: "H", 73: "I", 74: "J", 75: "K", 76: "L", 77: "M", 78: "N", 79: "O",
		80: "P", 81: "Q", 82: "R", 83: "S", 84: "T", 85: "U", 86: "V", 87: "W",
		88: "X", 89: "Y", 90: "Z", 91: "[", 92: "\\", 93: "]", 94: "^", 95: "_",
		96: "`", 97: "a", 98: "b", 99: "c", 100: "d", 101: "e", 102: "f", 103: "g",
		104: "h", 105: "i", 106: "j", 107: "k", 108: "l", 109: "m", 110: "n", 111: "o",
		112: "p", 113: "q", 114: "r", 115: "s", 116: "t", 117: "u", 118: "v", 119: "w",
		120: "x", 121: "y", 122: "z", 123: "{", 124: "|", 125: "}", 126: "~", 127: "<DEL>",
		128: "�", 129: "", 130: "‚", 131: "ƒ", 132: "„", 133: "…", 134: "†", 135: "‡",
		136: "ˆ", 137: "‰", 138: "Š", 139: "‹", 140: "Œ", 141: "", 142: "Ž", 143: "",
		144: "", 145: "‘", 146: "’", 147: "“", 148: "”", 149: "•", 150: "–", 151: "—",
		152: "˜", 153: "™", 154: "š", 155: "›", 156: "œ", 157: "", 158: "ž", 159: "Ÿ",
		160: "", 161: "¡", 162: "¢", 163: "£", 164: "¤", 165: "¥", 166: "¦", 167: "§",
		168: "¨", 169: "©", 170: "ª", 171: "«", 172: "¬", 173: "", 174: "®", 175: "¯",
		176: "°", 177: "±", 178: "²", 179: "³", 180: "´", 181: "µ", 182: "¶", 183: "·",
		184: "¸", 185: "¹", 186: "º", 187: "»", 188: "¼", 189: "½", 190: "¾", 191: "¿",
		192: "À", 193: "Á", 194: "Â", 195: "Ã", 196: "Ä", 197: "Å", 198: "Æ", 199: "Ç",
		200: "È", 201: "É", 202: "Ê", 203: "Ë", 204: "Ì", 205: "Í", 206: "Î", 207: "Ï",
		208: "Ð", 209: "Ñ", 210: "Ò", 211: "Ó", 212: "Ô", 213: "Õ", 214: "Ö", 215: "×",
		216: "Ø", 217: "Ù", 218: "Ú", 219: "Û", 220: "Ü", 221: "Ý", 222: "Þ", 223: "ß",
		224: "à", 225: "á", 226: "â", 227: "ã", 228: "ä", 229: "å", 230: "æ", 231: "ç",
		232: "è", 233: "é", 234: "ê", 235: "ë", 236: "ì", 237: "í", 238: "î", 239: "ï",
		240: "ð", 241: "ñ", 242: "ò", 243: "ó", 244: "ô", 245: "õ", 246: "ö", 247: "÷",
		248: "ø", 249: "ù", 250: "ú", 251: "û", 252: "ü", 253: "ý", 254: "þ", 255: "ÿ",
	}

	return repo
}

func (repo *CharacterRepo) GetANSICharFromBin(bin string, includeControlCharacters bool) string {
	for _, char := range repo.asciiAnsiItems["ANSI"] {
		if strings.Contains(char, bin) {
			if !includeControlCharacters && strings.HasPrefix(char, "<") && strings.HasSuffix(char, ">") {
				return ""
			}
			return char
		}
	}
	return ""
}

func (repo *CharacterRepo) GetANSICharFromDec(dec int, includeControlCharacters bool) string {
	char, exists := repo.asciiAnsiItems["ANSI"][dec]
	if exists {
		if !includeControlCharacters && strings.HasPrefix(char, "<") && strings.HasSuffix(char, ">") {
			return ""
		}
		return char
	}
	return ""
}

func (repo *CharacterRepo) GetASCIICharFromBin(bin string, includeControlCharacters bool) string {
	for _, char := range repo.asciiAnsiItems["ASCII"] {
		if strings.Contains(char, bin) {
			if !includeControlCharacters && strings.HasPrefix(char, "<") && strings.HasSuffix(char, ">") {
				return ""
			}
			return char
		}
	}
	return ""
}

func (repo *CharacterRepo) GetASCIICharFromDec(dec int, includeControlCharacters bool) string {
	char, exists := repo.asciiAnsiItems["ASCII"][dec]
	if exists {
		if !includeControlCharacters && strings.HasPrefix(char, "<") && strings.HasSuffix(char, ">") {
			return ""
		}
		return char
	}
	return ""
}

func (repo *CharacterRepo) GetGematriaRunes() []string {
	return []string{
		"ᛝ", // ING
		"ᛟ", // OE
		"ᛇ", // EO
		"ᛡ", // IO
		"ᛠ", // EA
		"ᚫ", // AE
		"ᚦ", // TH
		"ᚠ", // F
		"ᚢ", // U
		"ᚩ", // O
		"ᚱ", // R
		"ᚳ", // C/K
		"ᚷ", // G
		"ᚹ", // W
		"ᚻ", // H
		"ᚾ", // N
		"ᛁ", // I
		"ᛄ", // J
		"ᛈ", // P
		"ᛉ", // X
		"ᛋ", // S/Z
		"ᛏ", // T
		"ᛒ", // B
		"ᛖ", // E
		"ᛗ", // M
		"ᛚ", // L
		"ᛞ", // D
		"ᚪ", // A
		"ᚣ", // Y
	}
}

func (repo *CharacterRepo) GetCharFromRune(value string) string {
	var retval string
	switch value {
	case "ᛝ":
		retval = "ING"
	case "ᛟ":
		retval = "OE"
	case "ᛇ":
		retval = "EO"
	case "ᛡ":
		retval = "IO"
	case "ᛠ":
		retval = "EA"
	case "ᚫ":
		retval = "AE"
	case "ᚦ":
		retval = "TH"
	case "ᚠ":
		retval = "F"
	case "ᚢ":
		retval = "U"
	case "ᚩ":
		retval = "O"
	case "ᚱ":
		retval = "R"
	case "ᚳ":
		retval = "C"
	case "ᚷ":
		retval = "G"
	case "ᚹ":
		retval = "W"
	case "ᚻ":
		retval = "H"
	case "ᚾ":
		retval = "N"
	case "ᛁ":
		retval = "I"
	case "ᛄ":
		retval = "J"
	case "ᛈ":
		retval = "P"
	case "ᛉ":
		retval = "X"
	case "ᛋ":
		retval = "S"
	case "ᛏ":
		retval = "T"
	case "ᛒ":
		retval = "B"
	case "ᛖ":
		retval = "E"
	case "ᛗ":
		retval = "M"
	case "ᛚ":
		retval = "L"
	case "ᛞ":
		retval = "D"
	case "ᚪ":
		retval = "A"
	case "ᚣ":
		retval = "Y"
	case "•":
		retval = " "
	case "⊹":
		retval = "."
	default:
		retval = value
	}
	return retval
}

func (repo *CharacterRepo) GetRuneFromChar(value string) string {
	var retval string
	switch value {
	case "ING", "NG":
		retval = "ᛝ"
	case "OE":
		retval = "ᛟ"
	case "EO":
		retval = "ᛇ"
	case "IO", "IA":
		retval = "ᛡ"
	case "EA":
		retval = "ᛠ"
	case "AE":
		retval = "ᚫ"
	case "TH":
		retval = "ᚦ"
	case "F":
		retval = "ᚠ"
	case "V", "U":
		retval = "ᚢ"
	case "O":
		retval = "ᚩ"
	case "R":
		retval = "ᚱ"
	case "Q", "K", "C":
		retval = "ᚳ"
	case "G":
		retval = "ᚷ"
	case "W":
		retval = "ᚹ"
	case "H":
		retval = "ᚻ"
	case "N":
		retval = "ᚾ"
	case "I":
		retval = "ᛁ"
	case "J":
		retval = "ᛄ"
	case "P":
		retval = "ᛈ"
	case "X":
		retval = "ᛉ"
	case "Z", "S":
		retval = "ᛋ"
	case "T":
		retval = "ᛏ"
	case "B":
		retval = "ᛒ"
	case "E":
		retval = "ᛖ"
	case "M":
		retval = "ᛗ"
	case "L":
		retval = "ᛚ"
	case "D":
		retval = "ᛞ"
	case "A":
		retval = "ᚪ"
	case "Y":
		retval = "ᚣ"
	case " ":
		retval = "•"
	case ".":
		retval = "⊹"
	default:
		retval = value
	}
	return retval
}

func (repo *CharacterRepo) IsRune(value string, includeDunkus bool) bool {
	if includeDunkus {
		if value == "•" || value == "⊹" {
			return true
		}
	}

	switch value {
	case "ᛝ", "ᛟ", "ᛇ", "ᛡ", "ᛠ", "ᚫ", "ᚦ", "ᚠ", "ᚢ", "ᚩ", "ᚱ", "ᚳ", "ᚷ", "ᚹ", "ᚻ", "ᚾ", "ᛁ", "ᛄ", "ᛈ", "ᛉ", "ᛋ", "ᛏ", "ᛒ", "ᛖ", "ᛗ", "ᛚ", "ᛞ", "ᚪ", "ᚣ":
		return true
	default:
		return false
	}
}

func (repo *CharacterRepo) GetValueFromRune(rune string) int {
	var retval int
	switch rune {
	case "ᛝ":
		retval = 79
	case "ᛟ":
		retval = 83
	case "ᛇ":
		retval = 41
	case "ᛡ":
		retval = 107
	case "ᛠ":
		retval = 109
	case "ᚫ":
		retval = 101
	case "ᚦ":
		retval = 5
	case "ᚠ":
		retval = 2
	case "ᚢ":
		retval = 3
	case "ᚩ":
		retval = 7
	case "ᚱ":
		retval = 11
	case "ᚳ":
		retval = 13
	case "ᚷ":
		retval = 17
	case "ᚹ":
		retval = 19
	case "ᚻ":
		retval = 23
	case "ᚾ":
		retval = 29
	case "ᛁ":
		retval = 31
	case "ᛄ":
		retval = 37
	case "ᛈ":
		retval = 43
	case "ᛉ":
		retval = 47
	case "ᛋ":
		retval = 53
	case "ᛏ":
		retval = 59
	case "ᛒ":
		retval = 61
	case "ᛖ":
		retval = 67
	case "ᛗ":
		retval = 71
	case "ᛚ":
		retval = 73
	case "ᛞ":
		retval = 89
	case "ᚪ":
		retval = 97
	case "ᚣ":
		retval = 103
	default:
		retval = 0
	}
	return retval
}

func (repo *CharacterRepo) GetRuneFromValue(value int) string {
	var retval string
	switch value {
	case 2:
		retval = "ᚠ"
	case 3:
		retval = "ᚢ"
	case 5:
		retval = "ᚦ"
	case 7:
		retval = "ᚩ"
	case 11:
		retval = "ᚱ"
	case 13:
		retval = "ᚳ"
	case 17:
		retval = "ᚷ"
	case 19:
		retval = "ᚹ"
	case 23:
		retval = "ᚻ"
	case 29:
		retval = "ᚾ"
	case 31:
		retval = "ᛁ"
	case 37:
		retval = "ᛄ"
	case 41:
		retval = "ᛇ"
	case 43:
		retval = "ᛈ"
	case 47:
		retval = "ᛉ"
	case 53:
		retval = "ᛋ"
	case 59:
		retval = "ᛏ"
	case 61:
		retval = "ᛒ"
	case 67:
		retval = "ᛖ"
	case 71:
		retval = "ᛗ"
	case 73:
		retval = "ᛚ"
	case 79:
		retval = "ᛝ"
	case 83:
		retval = "ᛟ"
	case 89:
		retval = "ᛞ"
	case 97:
		retval = "ᚪ"
	case 101:
		retval = "ᚫ"
	case 103:
		retval = "ᚣ"
	case 107:
		retval = "ᛡ"
	case 109:
		retval = "ᛠ"
	default:
		retval = ""
	}
	return retval
}
