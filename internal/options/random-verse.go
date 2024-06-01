package options

type VerseRandomReferenceOption interface {
	applyVRR(*Verse, *RandomReference)
}

func MakeVerseRandomReferenceOptions(opts []VerseRandomReferenceOption) (*Verse, *RandomReference) {
	v := defaultVerse()
	rr := defaultRandomReference()
	for _, o := range opts {
		o.applyVRR(v, rr)
	}
	return v, rr
}
