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

var lb = MustLookupBook
var lbe = MustLookupBookExtract
var Categories = map[string][]BookExtract{
	"Law": {
		lb("Genesis"),
		lb("Exodus"),
		lb("Leviticus"),
		lb("Numbers"),
		lb("Deuteronomy"),
	},
	"History": {
		lb("Joshua"),
		lb("Judges"),
		lb("Ruth"),
		lb("1 Samuel"),
		lb("2 Samuel"),
		lb("1 Kings"),
		lb("2 Kings"),
		lb("1 Chronicles"),
		lb("2 Chronicles"),
		lb("Ezra"),
		lb("Nehemiah"),
		lb("Esther"),
	},
	"Wisdom": {
		lb("Job"),
		lb("Psalms"),
		lb("Proverbs"),
		lb("Ecclesiastes"),
		lb("Song of Solomon"),
	},
	"Prophets": {
		lb("Isaiah"),
		lb("Jeremiah"),
		lb("Lamentations"),
		lb("Ezekiel"),
		lb("Daniel"),
		lb("Hosea"),
		lb("Joel"),
		lb("Amos"),
		lb("Obadiah"),
		lb("Jonah"),
		lb("Micah"),
		lb("Nahum"),
		lb("Habakkuk"),
		lb("Zephaniah"),
		lb("Haggai"),
		lb("Zechariah"),
		lb("Malachi"),
	},
	"Gospels": {
		lb("Matthew"),
		lb("Mark"),
		lb("Luke"),
		lb("John"),

		// And Acts because it's a sequel to Luke
		lb("Acts"),
	},
	"Epistles": {
		lb("Romans"),
		lb("1 Corinthians"),
		lb("2 Corinthians"),
		lb("Galatians"),
		lb("Ephesians"),
		lb("Philippians"),
		lb("Colossians"),
		lb("1 Thessalonians"),
		lb("2 Thessalonians"),
		lb("1 Timothy"),
		lb("2 Timothy"),
		lb("Titus"),
		lb("Philemon"),
		lb("Hebrews"),
		lb("James"),
		lb("1 Peter"),
		lb("2 Peter"),
		lb("1 John"),
		lb("2 John"),
		lb("3 John"),
		lb("Jude"),
	},
	"Apocalyptic": {
		// Apocalypse passages
		lbe("Daniel", "7:1", "*:*"),
		lb("Revelation"),

		// Proto-Apocalyptic Passages
		lbe("Amos", "7:1", "7:9"),
		lbe("Amos", "8:1", "8:13"),
		lbe("Isaiah", "24:1", "27:*"),
		lbe("Isaiah", "33:1", "33:*"),
		lbe("Isaiah", "55:1", "56:*"),
		lbe("Jeremiah", "1:11", "1:16"),
		lbe("Ezekiel", "38:1", "39:*"),
		lbe("Zechariah", "9:1", "14:*"),
		lb("Joel"),
	},
}
