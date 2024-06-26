package ref

var Abbreviations = &BookAbbreviations{
	Abbreviations: []BookAbbreviation{
		{
			Name:      "Genesis",
			Preferred: "Gen.",
			Accepts: []string{
				"Genesis",
				"Gn",
			},
		},
		{
			Name:      "Exodus",
			Preferred: "Ex.",
			Accepts: []string{
				"Exodus",
			},
		},
		{
			Name:      "Leviticus",
			Preferred: "Lev.",
			Accepts: []string{
				"Leviticus",
				"Lv",
			},
		},
		{
			Name:      "Numbers",
			Preferred: "Num.",
			Accepts: []string{
				"Numbers",
				"Nm",
				"Nb",
			},
		},
		{
			Name:      "Deuteronomy",
			Preferred: "Deut.",
			Accepts: []string{
				"Deuteronomy",
				"Dt",
			},
		},
		{
			Name:      "Joshua",
			Preferred: "Josh.",
			Accepts: []string{
				"Joshua",
				"Jsh",
			},
		},
		{
			Name:      "Judges",
			Preferred: "Judg.",
			Accepts: []string{
				"Judges",
				"Jg",
				"Jdgs",
			},
		},
		{
			Name:      "Ruth",
			Preferred: "Ruth",
			Accepts: []string{
				"Ruth",
				"Rth",
			},
		},
		{
			Name:      "1 Samuel",
			Preferred: "1 Sam.",
			Ordinal:   1,
			Accepts: []string{
				"1Samuel",
				"FirstSamuel",
				"1stSamuel",
				"ISamuel",
				"ⅠSamuel",
				"1Sm",
				"FirstSm",
				"1stSm",
				"ISm",
				"ⅠSm",
			},
		},
		{
			Name:      "2 Samuel",
			Preferred: "2 Sam.",
			Ordinal:   2,
			Accepts: []string{
				"2Samuel",
				"SecondSamuel",
				"2ndSamuel",
				"IISamuel",
				"ⅡSamuel",
				"2Sm",
				"SecondSm",
				"2ndSm",
				"IISm",
				"ⅡSm",
			},
		},
		{
			Name:      "1 Kings",
			Preferred: "1 Kings",
			Ordinal:   1,
			Accepts: []string{
				"1Kings",
				"FirstKings",
				"1stKings",
				"IKings",
				"ⅠKings",
				"1Kgs",
				"FirstKgs",
				"1stKgs",
				"IKgs",
				"ⅠKgs",
			},
		},
		{
			Name:      "2 Kings",
			Preferred: "2 Kings",
			Ordinal:   2,
			Accepts: []string{
				"2Kings",
				"SecondKings",
				"2ndKings",
				"IIKings",
				"ⅡKings",
				"2Kgs",
				"SecondKgs",
				"2ndKgs",
				"IIKgs",
				"ⅡKgs",
			},
		},
		{
			Name:      "1 Chronicles",
			Preferred: "1 Chron.",
			Ordinal:   1,
			Accepts: []string{
				"1Chronicles",
				"FirstChronicles",
				"1stChronicles",
				"IChronicles",
				"ⅠChronicles",
				"1Chr",
				"FirstChr",
				"1stChr",
				"IChr",
				"ⅠChr",
			},
		},
		{
			Name:      "2 Chronicles",
			Preferred: "2 Chron.",
			Ordinal:   2,
			Accepts: []string{
				"2Chronicles",
				"SecondChronicles",
				"2ndChronicles",
				"IIChronicles",
				"ⅡChronicles",
				"2Chr",
				"SecondChr",
				"2ndChr",
				"IIChr",
				"ⅡChr",
			},
		},
		{
			Name:      "Ezra",
			Preferred: "Ezra",
			Accepts: []string{
				"Ezra",
			},
		},
		{
			Name:      "Nehemiah",
			Preferred: "Neh.",
			Accepts: []string{
				"Nehemiah",
			},
		},
		{
			Name:      "Esther",
			Preferred: "Est.",
			Accepts: []string{
				"Esther",
			},
		},
		{
			Name:      "Job",
			Preferred: "Job",
			Accepts: []string{
				"Job",
				"Jb",
			},
		},
		{
			Name:      "Psalms",
			Preferred: "Ps.",
			Singular:  "Psalm",
			Accepts: []string{
				"Psalms",
				"Pslm",
				"Psm",
				"Pss",
			},
		},
		{
			Name:      "Proverbs",
			Preferred: "Prov.",
			Accepts: []string{
				"Proverbs",
				"Prv",
			},
		},
		{
			Name:      "Ecclesiastes",
			Preferred: "Eccles.",
			Accepts: []string{
				"Ecclesiastes",
				"Qoheleth",
			},
		},
		{
			Name:      "Song of Solomon",
			Preferred: "Song",
			Accepts: []string{
				"Song of Solomon",
				"Song of Songs",
				"SOS",
				"Canticle of Canticles",
				"Canticles",
			},
		},
		{
			Name:      "Isaiah",
			Preferred: "Isa.",
			Accepts: []string{
				"Isaiah",
			},
		},
		{
			Name:      "Jeremiah",
			Preferred: "Jer.",
			Accepts: []string{
				"Jeremiah",
				"Jr",
			},
		},
		{
			Name:      "Lamentations",
			Preferred: "Lam.",
			Accepts: []string{
				"Lamentations",
			},
		},
		{
			Name:      "Ezekiel",
			Preferred: "Ezek.",
			Accepts: []string{
				"Ezekiel",
				"Ezk",
			},
		},
		{
			Name:      "Daniel",
			Preferred: "Dan.",
			Accepts: []string{
				"Daniel",
				"Dn",
			},
		},
		{
			Name:      "Hosea",
			Preferred: "Hos.",
			Accepts: []string{
				"Hosea",
			},
		},
		{
			Name:      "Joel",
			Preferred: "Joel",
			Accepts: []string{
				"Joel",
				"Jl",
			},
		},
		{
			Name:      "Amos",
			Preferred: "Amos",
			Accepts: []string{
				"Amos",
			},
		},
		{
			Name:      "Obadiah",
			Preferred: "Obad.",
			Accepts: []string{
				"Obadiah",
			},
		},
		{
			Name:      "Jonah",
			Preferred: "Jonah",
			Accepts: []string{
				"Jonah",
				"Jnh",
			},
		},
		{
			Name:      "Micah",
			Preferred: "Mic.",
			Accepts: []string{
				"Micah",
				"Mc",
			},
		},
		{
			Name:      "Nahum",
			Preferred: "Nah.",
			Accepts: []string{
				"Nahum",
			},
		},
		{
			Name:      "Habakkuk",
			Preferred: "Hab.",
			Accepts: []string{
				"Habakkuk",
				"Hb",
			},
		},
		{
			Name:      "Zephaniah",
			Preferred: "Zeph.",
			Accepts: []string{
				"Zephaniah",
				"Zp",
			},
		},
		{
			Name:      "Haggai",
			Preferred: "Hag.",
			Accepts: []string{
				"Haggai",
				"Hg",
			},
		},
		{
			Name:      "Zechariah",
			Preferred: "Zech.",
			Accepts: []string{
				"Zechariah",
				"Zc",
			},
		},
		{
			Name:      "Malachi",
			Preferred: "Mal.",
			Accepts: []string{
				"Malachi",
				"Ml",
			},
		},
		{
			Name:      "Matthew",
			Preferred: "Matt.",
			Accepts: []string{
				"Matthew",
				"Mt",
			},
		},
		{
			Name:      "Mark",
			Preferred: "Mark",
			Accepts: []string{
				"Mark",
				"Mrk",
				"Mk",
			},
		},
		{
			Name:      "Luke",
			Preferred: "Luke",
			Accepts: []string{
				"Luke",
				"Lk",
			},
		},
		{
			Name:      "John",
			Preferred: "John",
			Accepts: []string{
				"John",
				"Jhn",
				"Jn",
			},
		},
		{
			Name:      "Acts",
			Preferred: "Acts",
			Accepts: []string{
				"Acts",
			},
		},
		{
			Name:      "Romans",
			Preferred: "Rom.",
			Accepts: []string{
				"Romans",
				"Rm",
			},
		},
		{
			Name:      "1 Corinthians",
			Preferred: "1 Cor.",
			Ordinal:   1,
			Accepts: []string{
				"1Corinthians",
				"FirstCorinthians",
				"1stCorinthians",
				"ICorinthians",
				"ⅠCorinthians",
			},
		},
		{
			Name:      "2 Corinthians",
			Preferred: "2 Cor.",
			Ordinal:   2,
			Accepts: []string{
				"2Corinthians",
				"SecondCorinthians",
				"2ndCorinthians",
				"IICorinthians",
				"ⅡCorinthians",
			},
		},
		{
			Name:      "Galatians",
			Preferred: "Gal.",
			Accepts: []string{
				"Galatians",
			},
		},
		{
			Name:      "Ephesians",
			Preferred: "Eph.",
			Accepts: []string{
				"Ephesians",
			},
		},
		{
			Name:      "Philippians",
			Preferred: "Phil.",
			Accepts: []string{
				"Philippians",
				"Php",
				"Pp",
			},
		},
		{
			Name:      "Colossians",
			Preferred: "Col.",
			Accepts: []string{
				"Colossians",
			},
		},
		{
			Name:      "1 Thessalonians",
			Preferred: "1 Thess.",
			Ordinal:   1,
			Accepts: []string{
				"1Thessalonians",
				"FirstThessalonians",
				"1stThessalonians",
				"IThessalonians",
				"ⅠThessalonians",
			},
		},
		{
			Name:      "2 Thessalonians",
			Preferred: "2 Thess.",
			Ordinal:   2,
			Accepts: []string{
				"2Thessalonians",
				"SecondThessalonians",
				"2ndThessalonians",
				"IIThessalonians",
				"ⅡThessalonians",
			},
		},
		{
			Name:      "1 Timothy",
			Preferred: "1 Tim.",
			Ordinal:   1,
			Accepts: []string{
				"1Timothy",
				"FirstTimothy",
				"1stTimothy",
				"ITimothy",
				"ⅠTimothy",
			},
		},
		{
			Name:      "2 Timothy",
			Preferred: "2 Tim.",
			Ordinal:   2,
			Accepts: []string{
				"2Timothy",
				"SecondTimothy",
				"2ndTimothy",
				"IITimothy",
				"ⅡTimothy",
			},
		},
		{
			Name:      "Titus",
			Preferred: "Titus",
			Accepts: []string{
				"Titus",
			},
		},
		{
			Name:      "Philemon",
			Preferred: "Philem.",
			Accepts: []string{
				"Philemon",
				"Phm",
				"Pm",
			},
		},
		{
			Name:      "Hebrews",
			Preferred: "Heb.",
			Accepts: []string{
				"Hebrews",
			},
		},
		{
			Name:      "James",
			Preferred: "James",
			Accepts: []string{
				"James",
				"Jas",
				"Jm",
			},
		},
		{
			Name:      "1 Peter",
			Preferred: "1 Pet.",
			Ordinal:   1,
			Accepts: []string{
				"1Peter",
				"FirstPeter",
				"1stPeter",
				"IPeter",
				"ⅠPeter",
				"1Pt",
				"FirstPt",
				"1stPt",
				"IPt",
				"ⅠPt",
			},
		},
		{
			Name:      "2 Peter",
			Preferred: "2 Pet.",
			Ordinal:   2,
			Accepts: []string{
				"2Peter",
				"SecondPeter",
				"2ndPeter",
				"IIPeter",
				"ⅡPeter",
				"2Pt",
				"SecondPt",
				"2ndPt",
				"IIPt",
				"ⅡPt",
			},
		},
		{
			Name:      "1 John",
			Preferred: "1 John",
			Ordinal:   1,
			Accepts: []string{
				"1John",
				"FirstJohn",
				"1stJohn",
				"IJohn",
				"ⅠJohn",
				"1Jhn",
				"FirstJhn",
				"1stJhn",
				"IJhn",
				"ⅠJhn",
				"1Jn",
				"FirstJn",
				"1stJn",
				"IJn",
				"ⅠJn",
			},
		},
		{
			Name:      "2 John",
			Preferred: "2 John",
			Ordinal:   2,
			Accepts: []string{
				"2John",
				"SecondJohn",
				"2ndJohn",
				"IIJohn",
				"ⅡJohn",
				"2Jhn",
				"SecondJhn",
				"2ndJhn",
				"IIJhn",
				"ⅡJhn",
				"2Jn",
				"SecondJn",
				"2ndJn",
				"IIJn",
				"ⅡJn",
			},
		},
		{
			Name:      "3 John",
			Preferred: "3 John",
			Ordinal:   3,
			Accepts: []string{
				"3John",
				"ThirdJohn",
				"3rdJohn",
				"IIIJohn",
				"ⅢJohn",
				"3Jhn",
				"ThirdJhn",
				"3rdJhn",
				"IIIJhn",
				"ⅢJhn",
				"3Jn",
				"ThirdJn",
				"3rdJn",
				"IIIJn",
				"ⅢJn",
			},
		},
		{
			Name:      "Jude",
			Preferred: "Jude",
			Accepts: []string{
				"Jude",
				"Jd",
			},
		},
		{
			Name:      "Revelation",
			Preferred: "Rev.",
			Accepts: []string{
				"Revelation",
				"Rv",
				"The Revelation",
			},
		},
	},
}
