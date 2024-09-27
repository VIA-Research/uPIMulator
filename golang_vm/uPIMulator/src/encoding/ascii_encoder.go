package encoding

type AsciiEncoder struct {
	table          map[string]uint8
	inverted_table map[uint8]string
	unknown        string
}

func (this *AsciiEncoder) Init() {
	this.table = make(map[string]uint8)
	this.inverted_table = make(map[uint8]string)

	this.table["\t"] = 9
	this.table[" "] = 32
	this.table["!"] = 33
	this.table["\""] = 34
	this.table["#"] = 35
	this.table["$"] = 36
	this.table["%"] = 37
	this.table["&"] = 38
	this.table["'"] = 39
	this.table["("] = 40
	this.table[")"] = 41
	this.table["*"] = 42
	this.table["+"] = 43
	this.table[","] = 44
	this.table["-"] = 45
	this.table["."] = 46
	this.table["/"] = 47
	this.table["0"] = 48
	this.table["1"] = 49
	this.table["2"] = 50
	this.table["3"] = 51
	this.table["4"] = 52
	this.table["5"] = 53
	this.table["6"] = 54
	this.table["7"] = 55
	this.table["8"] = 56
	this.table["9"] = 57
	this.table[":"] = 58
	this.table[""] = 59
	this.table["<"] = 60
	this.table["="] = 61
	this.table[">"] = 62
	this.table["?"] = 63
	this.table["@"] = 64
	this.table["A"] = 65
	this.table["B"] = 66
	this.table["C"] = 67
	this.table["D"] = 68
	this.table["E"] = 69
	this.table["F"] = 70
	this.table["G"] = 71
	this.table["H"] = 72
	this.table["I"] = 73
	this.table["J"] = 74
	this.table["K"] = 75
	this.table["L"] = 76
	this.table["M"] = 77
	this.table["N"] = 78
	this.table["O"] = 79
	this.table["P"] = 80
	this.table["Q"] = 81
	this.table["R"] = 82
	this.table["S"] = 83
	this.table["T"] = 84
	this.table["U"] = 85
	this.table["V"] = 86
	this.table["W"] = 87
	this.table["X"] = 88
	this.table["Y"] = 89
	this.table["Z"] = 90
	this.table["["] = 91
	this.table["\\"] = 92
	this.table["]"] = 93
	this.table["^"] = 94
	this.table["_"] = 95
	this.table["`"] = 96
	this.table["a"] = 97
	this.table["b"] = 98
	this.table["c"] = 99
	this.table["d"] = 100
	this.table["e"] = 101
	this.table["f"] = 102
	this.table["g"] = 103
	this.table["h"] = 104
	this.table["i"] = 105
	this.table["j"] = 106
	this.table["k"] = 107
	this.table["l"] = 108
	this.table["m"] = 109
	this.table["n"] = 110
	this.table["o"] = 111
	this.table["p"] = 112
	this.table["q"] = 113
	this.table["r"] = 114
	this.table["s"] = 115
	this.table["t"] = 116
	this.table["u"] = 117
	this.table["v"] = 118
	this.table["w"] = 119
	this.table["x"] = 120
	this.table["y"] = 121
	this.table["z"] = 122
	this.table["{"] = 123
	this.table["|"] = 124
	this.table["}"] = 125
	this.table["~"] = 126
	this.table["Ç"] = 128
	this.table["ü"] = 129
	this.table["é"] = 130
	this.table["â"] = 131
	this.table["ä"] = 132
	this.table["à"] = 133
	this.table["å"] = 134
	this.table["ç"] = 135
	this.table["ê"] = 136
	this.table["ë"] = 137
	this.table["è"] = 138
	this.table["ï"] = 139
	this.table["î"] = 140
	this.table["ì"] = 141
	this.table["Ä"] = 142
	this.table["Å"] = 143
	this.table["É"] = 144
	this.table["æ"] = 145
	this.table["Æ"] = 146
	this.table["ô"] = 147
	this.table["ö"] = 148
	this.table["ò"] = 149
	this.table["û"] = 150
	this.table["ù"] = 151
	this.table["ÿ"] = 152
	this.table["Ö"] = 153
	this.table["Ü"] = 154
	this.table["ø"] = 155
	this.table["£"] = 156
	this.table["Ø"] = 157
	this.table["×"] = 158
	this.table["ƒ"] = 159
	this.table["á"] = 160
	this.table["í"] = 161
	this.table["ó"] = 162
	this.table["ú"] = 163
	this.table["ñ"] = 164
	this.table["Ñ"] = 165
	this.table["ª"] = 166
	this.table["º"] = 167
	this.table["¿"] = 168
	this.table["®"] = 169
	this.table["¬"] = 170
	this.table["½"] = 171
	this.table["¼"] = 172
	this.table["¡"] = 173
	this.table["«"] = 174
	this.table["»"] = 175
	this.table["░"] = 176
	this.table["▒"] = 177
	this.table["▓"] = 178
	this.table["│"] = 179
	this.table["┤"] = 180
	this.table["Á"] = 181
	this.table["Â"] = 182
	this.table["À"] = 183
	this.table["©"] = 184
	this.table["╣"] = 185
	this.table["║"] = 186
	this.table["╗"] = 187
	this.table["╝"] = 188
	this.table["¢"] = 189
	this.table["¥"] = 190
	this.table["┐"] = 191
	this.table["└"] = 192
	this.table["┴"] = 193
	this.table["┬"] = 194
	this.table["├"] = 195
	this.table["─"] = 196
	this.table["┼"] = 197
	this.table["ã"] = 198
	this.table["Ã"] = 199
	this.table["╚"] = 200
	this.table["╔"] = 201
	this.table["╩"] = 202
	this.table["╦"] = 203
	this.table["╠"] = 204
	this.table["═"] = 205
	this.table["╬"] = 206
	this.table["¤"] = 207
	this.table["ð"] = 208
	this.table["Ð"] = 209
	this.table["Ê"] = 210
	this.table["Ë"] = 211
	this.table["È"] = 212
	this.table["ı"] = 213
	this.table["Í"] = 214
	this.table["Î"] = 215
	this.table["Ï"] = 216
	this.table["┘"] = 217
	this.table["┌"] = 218
	this.table["█"] = 219
	this.table["▄"] = 220
	this.table["¦"] = 221
	this.table["Ì"] = 222
	this.table["▀"] = 223
	this.table["Ó"] = 224
	this.table["ß"] = 225
	this.table["Ô"] = 226
	this.table["Ò"] = 227
	this.table["õ"] = 228
	this.table["Õ"] = 229
	this.table["µ"] = 230
	this.table["þ"] = 231
	this.table["Þ"] = 232
	this.table["Ú"] = 233
	this.table["Û"] = 234
	this.table["Ù"] = 235
	this.table["ý"] = 236
	this.table["Ý"] = 237
	this.table["¯"] = 238
	this.table["´"] = 239
	this.table["≡"] = 240
	this.table["±"] = 241
	this.table["‗"] = 242
	this.table["¾"] = 243
	this.table["¶"] = 244
	this.table["§"] = 245
	this.table["÷"] = 246
	this.table["¸"] = 247
	this.table["°"] = 248
	this.table["¨"] = 249
	this.table["·"] = 250
	this.table["¹"] = 251
	this.table["³"] = 252
	this.table["²"] = 253
	this.table["■"] = 254

	for character, value := range this.table {
		this.inverted_table[value] = character
	}

	this.unknown = "■"
}

func (this *AsciiEncoder) Encode(characters string) *ByteStream {
	byte_stream := new(ByteStream)
	byte_stream.Init()

	for _, character := range characters {
		value := this.table[string(character)]
		byte_stream.Append(value)
	}

	return byte_stream
}

func (this *AsciiEncoder) Decode(byte_stream *ByteStream) string {
	characters := ""

	for i := int64(0); i < byte_stream.Size(); i++ {
		value := byte_stream.Get(int(i))
		characters += this.inverted_table[value]
	}

	return characters
}

func (this *AsciiEncoder) Unknown() string {
	return this.unknown
}
