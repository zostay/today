package ref

// TODO I'd really like these categories instead:
// 1. God (Father, Son, Holy Spirit)
// 2. Creation (Origins, Fall, Nature)
// 3. History (Narrative and Poetic, including the historical parts Gospels and Acts)
// 4. Law (All Commands, NT and OT)
// 5. Wisdom and Song (All the wisdom and singing in the Bible)
// 6. Gospel (All presentations of Jesus and Salvation)
// 7. Apocalypse and Prophecy (Looking forward to the end of all things)
// These categories would overlap heavily. I think I will need to do my own
// survey of the Bible to get this to approximate what I'm looking for.

func bk(name string) *Pericope {
	p, err := Lookup(Canonical, name+"1:1ffb", "")
	if err != nil {
		panic(err)
	}
	return p
}

func vv(ref string) *Pericope {
	p, err := Lookup(Canonical, ref, "")
	if err != nil {
		panic(err)
	}
	return p
}

var Categories = map[string][]*Pericope{
	"Law": {
		bk("Genesis"),
		bk("Exodus"),
		bk("Leviticus"),
		bk("Numbers"),
		bk("Deuteronomy"),
	},
	"History": {
		bk("Joshua"),
		bk("Judges"),
		bk("Ruth"),
		bk("1 Samuel"),
		bk("2 Samuel"),
		bk("1 Kings"),
		bk("2 Kings"),
		bk("1 Chronicles"),
		bk("2 Chronicles"),
		bk("Ezra"),
		bk("Nehemiah"),
		bk("Esther"),
	},
	"Wisdom": {
		bk("Job"),
		bk("Psalms"),
		bk("Proverbs"),
		bk("Ecclesiastes"),
		bk("Song of Solomon"),
	},
	"Prophets": {
		bk("Isaiah"),
		bk("Jeremiah"),
		bk("Lamentations"),
		bk("Ezekiel"),
		bk("Daniel"),
		bk("Hosea"),
		bk("Joel"),
		bk("Amos"),
		bk("Obadiah"),
		bk("Jonah"),
		bk("Micah"),
		bk("Nahum"),
		bk("Habakkuk"),
		bk("Zephaniah"),
		bk("Haggai"),
		bk("Zechariah"),
		bk("Malachi"),
	},
	"Gospels": {
		bk("Matthew"),
		bk("Mark"),
		bk("Luke"),
		bk("John"),

		// And Acts because it's a sequel to Luke
		bk("Acts"),
	},
	"Epistles": {
		bk("Romans"),
		bk("1 Corinthians"),
		bk("2 Corinthians"),
		bk("Galatians"),
		bk("Ephesians"),
		bk("Philippians"),
		bk("Colossians"),
		bk("1 Thessalonians"),
		bk("2 Thessalonians"),
		bk("1 Timothy"),
		bk("2 Timothy"),
		bk("Titus"),
		bk("Philemon"),
		bk("Hebrews"),
		bk("James"),
		bk("1 Peter"),
		bk("2 Peter"),
		bk("1 John"),
		bk("2 John"),
		bk("3 John"),
		bk("Jude"),
	},
	"Apocalyptic": {
		// Apocalypse passages
		vv("Daniel 7:1ffb"),
		bk("Revelation"),

		// Proto-Apocalyptic Passages
		vv("Amos 7:1-9"),
		vv("Amos 8:1-13"),
		vv("Isaiah 24:1-27:13"),
		vv("Isaiah 33:1ff"),
		vv("Isaiah 55:1-56:12"),
		vv("Jeremiah 1:11-16"),
		vv("Ezekiel 38:1-39:29"),
		vv("Zechariah 9:1ffb"),
		bk("Joel"),
	},
}
