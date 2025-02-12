package font

type Font4x4 struct {
	Pixels [4][4]uint8
}

var One = Font4x4{
	Pixels: [4][4]uint8{
		{0, 0, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
	},
}

var Two = Font4x4{
	Pixels: [4][4]uint8{
		{1, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{1, 1, 1, 1},
	},
}

var Three = Font4x4{
	Pixels: [4][4]uint8{
		{1, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 1, 0},
		{1, 1, 1, 0},
	},
}

var Four = Font4x4{
	Pixels: [4][4]uint8{
		{1, 0, 1, 0},
		{1, 0, 1, 0},
		{1, 1, 1, 0},
		{0, 0, 1, 0},
	},
}

var Five = Font4x4{
	Pixels: [4][4]uint8{
		{1, 1, 1, 0},
		{1, 0, 0, 0},
		{0, 1, 1, 0},
		{1, 1, 0, 0},
	},
}

var Six = Font4x4{
	Pixels: [4][4]uint8{
		{1, 1, 1, 0},
		{1, 0, 0, 0},
		{1, 1, 1, 0},
		{1, 1, 1, 0},
	},
}

var Seven = Font4x4{
	Pixels: [4][4]uint8{
		{1, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
	},
}

var Eight = Font4x4{
	Pixels: [4][4]uint8{
		{1, 1, 1, 0},
		{1, 0, 1, 0},
		{1, 1, 1, 0},
		{1, 1, 1, 0},
	},
}

var Nine = Font4x4{
	Pixels: [4][4]uint8{
		{1, 1, 1, 0},
		{1, 0, 1, 0},
		{1, 1, 1, 0},
		{0, 0, 1, 0},
	},
}

var Zero = Font4x4{
	Pixels: [4][4]uint8{
		{1, 1, 1, 0},
		{1, 0, 1, 0},
		{1, 0, 1, 0},
		{1, 1, 1, 0},
	},
}

var A = Font4x4{
	Pixels: [4][4]uint8{
		{0, 1, 0, 0},
		{1, 0, 1, 0},
		{1, 1, 1, 0},
		{1, 0, 1, 0},
	},
}

var B = Font4x4{
	Pixels: [4][4]uint8{
		{1, 0, 0, 0},
		{1, 1, 0, 0},
		{1, 0, 1, 0},
		{1, 1, 0, 0},
	},
}

var C = Font4x4{
	Pixels: [4][4]uint8{
		{0, 1, 1, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{0, 1, 1, 0},
	},
}

var D = Font4x4{
	Pixels: [4][4]uint8{
		{0, 0, 1, 0},
		{1, 1, 1, 0},
		{1, 0, 1, 0},
		{1, 1, 1, 0},
	},
}

var E = Font4x4{
	Pixels: [4][4]uint8{
		{1, 1, 1, 0},
		{1, 1, 0, 0},
		{1, 0, 0, 0},
		{1, 1, 1, 0},
	},
}

var F = Font4x4{
	Pixels: [4][4]uint8{
		{1, 1, 1, 0},
		{1, 1, 0, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 0},
	},
}

var CharMap = map[rune]Font4x4{
	'1': One,
	'2': Two,
	'3': Three,
	'4': Four,
	'5': Five,
	'6': Six,
	'7': Seven,
	'8': Eight,
	'9': Nine,
	'0': Zero,
	'A': A,
	'B': B,
	'C': C,
	'D': D,
	'E': E,
	'F': F,
}
