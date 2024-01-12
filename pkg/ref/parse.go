package ref

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type parseState struct {
	input []rune
	pos   int
}

func makeParseState(input string) *parseState {
	return &parseState{
		input: []rune(input),
	}
}

func (p *parseState) peek() rune {
	pos := p.pos
	for pos < len(p.input) && unicode.IsSpace(p.input[pos]) {
		pos++
	}

	if pos >= len(p.input) {
		return -1
	}

	return p.input[pos]
}

func (p *parseState) next() rune {
	pos := p.pos
	for pos < len(p.input) && unicode.IsSpace(p.input[pos]) {
		pos++
	}

	if pos >= len(p.input) {
		return -1
	}

	p.pos = pos + 1
	return p.input[pos]
}

func (p *parseState) expectRune(r rune) bool {
	if p.peek() != r {
		return false
	}
	p.next()
	return true
}

func (p *parseState) expect(pred func(rune) bool) (rune, bool) {
	r := p.peek()
	if !pred(r) {
		return -1, false
	}
	p.next()
	return r, true
}

func (p *parseState) expectWhile(pred func(rune) bool) []rune {
	rs := make([]rune, 0, 10)
	for {
		if r, ok := p.expect(pred); ok {
			rs = append(rs, r)
			continue
		}
		break
	}
	return rs
}

func (p *parseState) expectWhilePreserveWS(pred func(rune) bool) []rune {
	rs := make([]rune, 0, 10)
	for {
		was := p.pos
		if r, ok := p.expect(pred); ok {
			if p.pos-was > 1 {
				rs = append(rs, ' ')
			}
			rs = append(rs, r)
			continue
		}
		break
	}
	return rs
}

func (p *parseState) endOfInput() bool {
	return p.pos >= len(p.input)
}

func (p *parseState) remainderError() error {
	return &MoreInputError{Remaining: string(p.input[p.pos:])}
}

// ParseV will parse a single verse number. If there's trailing input after the
// verse number, the ref.V will be returned along with a MoreInputError.
func ParseV(ref string) (*V, error) {
	v, ps, err := expectV(*makeParseState(ref))
	if err != nil {
		return nil, err
	}

	if !ps.endOfInput() {
		return v, ps.remainderError()
	}

	return v, nil
}

// ParseCV will parse a chapter and verse number. If there's trailing input
// after the verse number, the ref.CV will be returned along with a
// MoreInputError.
func ParseCV(ref string) (*CV, error) {
	cv, ps, err := expectCV(*makeParseState(ref))
	if err != nil {
		return nil, err
	}

	if !ps.endOfInput() {
		return cv, ps.remainderError()
	}

	return cv, nil
}

// ParseSingle will parse a single verse reference, which may be a verse number
// or a chapter-and-verse reference. If there's trailing input after the verse
// number, the ref.Single will be returned along with a MoreInputError.
func ParseSingle(ref string) (*Single, error) {
	s, ps, err := expectSingle(*makeParseState(ref))
	if err != nil {
		return nil, err
	}

	if !ps.endOfInput() {
		return s, ps.remainderError()
	}

	return s, nil
}

// ParseAndFollowing will parse a verse reference followed by "ff" and then
// (optionally) either "b" or "c". If there's trailing input after the verse
// number, the ref.AndFollowing will be returned along with a MoreInputError.
func ParseAndFollowing(ref string) (*AndFollowing, error) {
	af, ps, err := expectAndFollowing(*makeParseState(ref))
	if err != nil {
		return nil, err
	}

	if !ps.endOfInput() {
		return af, ps.remainderError()
	}

	return af, nil
}

// ParseRange will parse a range of verses. If there's trailing input after the
// verse number, the ref.Range will be returned along with a MoreInputError.
func ParseRange(ref string) (*Range, error) {
	r, ps, err := expectRange(*makeParseState(ref))
	if err != nil {
		return nil, err
	}

	if !ps.endOfInput() {
		return r, ps.remainderError()
	}

	return r, nil
}

// ParseRelated will parse a relative reference, which may be a range of verses,
// a verse followed by "ff" and then (optionally) either "b" or "c", or a single
// verse reference. If there's trailing input after the verse number, the
// ref.Relative will be returned along with a MoreInputError.
//
// A trailing comma is considered an error.
func ParseRelated(ref string) (*Related, error) {
	r, ps, err := expectRelated(*makeParseState(ref))
	if err != nil {
		return nil, err
	}

	if !ps.endOfInput() {
		return r, ps.remainderError()
	}

	return r, err
}

// ParseProper will parse a proper reference, which is a book name or
// abbreviation followed by a relative reference. If there's trailing input
// after the verse number, the ref.Proper will be returned along with a
// MoreInputError.
func ParseProper(ref string) (*Proper, error) {
	p, ps, err := expectProper(*makeParseState(ref))
	if err != nil {
		return nil, err
	}

	if !ps.endOfInput() {
		return p, ps.remainderError()
	}

	return p, nil
}

// ParseMultiple will parse a multiple reference, which is a semicolon separated
// list of proper references. There must be at least one reference. The first
// reference must be a proper reference (book and related verse reference), but
// subsequence references may be relative to the first proper reference. If
// there's trailing input after the verse number, the ref.Multiple will be
// returned along with a MoreInputError.
//
// A trailing semi-colon is allowed and will be rolled into the MoreInputError.
func ParseMultiple(ref string) (*Multiple, error) {
	m, ps, err := expectMultiple(*makeParseState(ref))
	if err != nil {
		return nil, err
	}

	if !ps.endOfInput() {
		return m, ps.remainderError()
	}

	return m, err
}

func expectNumber(ref parseState) (int, parseState, error) {
	numr := ref.expectWhile(func(r rune) bool {
		return r >= '0' && r <= '9'
	})
	if len(numr) == 0 {
		return 0, ref, fmt.Errorf("%w: expected a number", ErrParseFail)
	}

	num, err := strconv.Atoi(string(numr))
	if err != nil {
		// should be unreachable
		return 0, ref, fmt.Errorf("%w: expected a number: %w", ErrParseFail, err)
	}

	return num, ref, nil
}

var (
	ErrParseFail = errors.New("parse failed")
)

// MoreInputError indicates that parse completed partially but there's additional
// input present.
type MoreInputError struct {
	// Remaining is the trailing input.
	Remaining string
}

// Error implements error.
func (e *MoreInputError) Error() string {
	return "parse completed, but more input remains"
}

var _ error = (*MoreInputError)(nil)

func expectV(ref parseState) (*V, parseState, error) {
	num, ps, err := expectNumber(ref)
	if err != nil {
		return nil, ref, err
	}

	v := &V{Verse: num}
	if err := v.Validate(); err != nil {
		return v, ps, err
	}

	return v, ps, nil
}

func expectCV(ref parseState) (*CV, parseState, error) {
	cnum, ps, err := expectNumber(ref)
	if err != nil {
		return nil, ref, err
	}

	cvBreak := func(r rune) bool {
		return r == ':' || r == '.'
	}

	if _, ok := ps.expect(cvBreak); !ok {
		return nil, ref, fmt.Errorf("%w: expected a colon", ErrParseFail)
	}

	vnum, ps, err := expectNumber(ps)
	if err != nil {
		return nil, ref, err
	}

	cv := &CV{Chapter: cnum, Verse: vnum}
	if err := cv.Validate(); err != nil {
		return cv, ps, err
	}

	return cv, ps, nil
}

func expectSingle(ref parseState) (s *Single, ps parseState, err error) {
	var v Verse
	v, ps, err = expectCV(ref)
	if errors.Is(err, ErrParseFail) {
		v, ps, err = expectV(ref)
		if err != nil {
			return nil, ref, err
		}
	}

	s = &Single{Verse: v}
	if err := s.Validate(); err != nil {
		return s, ps, err
	}

	return s, ps, nil
}

func expectAndFollowing(ref parseState) (*AndFollowing, parseState, error) {
	s, ps, err := expectSingle(ref)
	if err != nil {
		return nil, ref, err
	}

	if !ps.expectRune('f') {
		return nil, ref, fmt.Errorf("%w: %w", ErrParseFail, errors.New(`expected "ff"`))
	}

	if !ps.expectRune('f') {
		return nil, ref, fmt.Errorf("%w: %w", ErrParseFail, errors.New(`expected "ff"`))
	}

	untilLetters := func(r rune) bool {
		return r == 'b' || r == 'c'
	}

	f := &AndFollowing{Verse: s.Verse, Following: FollowingRemainingChapter}
	if r, ok := ps.expect(untilLetters); ok && r == 'b' {
		f = &AndFollowing{Verse: s.Verse, Following: FollowingRemainingBook}
	}

	if err := f.Validate(); err != nil {
		return f, ps, err
	}

	return f, ps, nil
}

func expectRange(ref parseState) (*Range, parseState, error) {
	sf, ps, err := expectSingle(ref)
	if err != nil {
		return nil, ref, err
	}

	if !ps.expectRune('-') {
		return nil, ref, fmt.Errorf("%w: expected a dash", ErrParseFail)
	}

	sl, ps, err := expectSingle(ps)
	if err != nil {
		return nil, ref, err
	}

	r := &Range{
		First: sf.Verse,
		Last:  sl.Verse,
	}
	if err := r.Validate(); err != nil {
		return r, ps, err
	}

	return r, ps, nil
}

func expectRelative(ref parseState) (rel Relative, ps parseState, err error) {
	rel, ps, err = expectRange(ref)
	if errors.Is(err, ErrParseFail) {
		rel, ps, err = expectAndFollowing(ps)
		if errors.Is(err, ErrParseFail) {
			return expectSingle(ps)
		}

		return rel, ps, err
	}

	return rel, ps, err
}

func expectRelated(ref parseState) (*Related, parseState, error) {
	var (
		ps  = ref
		err error
		rel Relative
	)

	rels := make([]Relative, 0, 10)
	for {
		rel, ps, err = expectRelative(ps)
		if err != nil {
			return nil, ref, err
		}

		rels = append(rels, rel)

		if !ps.expectRune(',') {
			break
		}
	}

	r := &Related{Refs: rels}
	if err := r.Validate(); err != nil {
		return r, ps, err
	}

	return r, ps, nil
}

func expectBookName(ref parseState) (string, parseState, error) {
	ps := ref
	firstLetter, ok := ps.expect(func(r rune) bool {
		return r >= 'a' && r <= 'z' ||
			r >= 'A' && r <= 'Z' ||
			r >= '1' && r <= '9'
	})

	if !ok {
		return "", ref, fmt.Errorf("%w: expected a letter or number to start book name or abbreviation", ErrParseFail)
	}

	rest := ps.expectWhilePreserveWS(func(r rune) bool {
		return r >= 'a' && r <= 'z' ||
			r >= 'A' && r <= 'Z' ||
			r == '.'
	})

	if len(rest) == 0 {
		return "", ref, fmt.Errorf("%w: expected a letter or number to start book name or abbreviation", ErrParseFail)
	}

	return string(append([]rune{firstLetter}, rest...)), ps, nil
}

func expectProper(ref parseState) (*Proper, parseState, error) {
	name, ps, err := expectBookName(ref)
	if err != nil {
		return nil, ref, err
	}

	var rel Relative
	rel, ps, err = expectRelated(ps)
	if err != nil {
		return nil, ref, err
	}

	if len(rel.(*Related).Refs) == 1 {
		rel = rel.(*Related).Refs[0]
	}

	p := &Proper{Book: name, Verse: rel}
	if err := p.Validate(); err != nil {
		return p, ps, err
	}

	return p, ps, nil
}

func expectMultiple(ref parseState) (*Multiple, parseState, error) {
	refs := make([]Ref, 0, 10)
	pref, ps, err := expectProper(ref)
	if err != nil {
		return nil, ref, err
	}

	refs = append(refs, pref)
	for {
		refPs := ps
		if !ps.expectRune(';') {
			break
		}

		pref, ps, err = expectProper(ps)
		if errors.Is(err, ErrParseFail) {
			var rel Relative
			rel, ps, err = expectRelated(ps)
			if err != nil {
				ps = refPs
				break
			}

			if len(rel.(*Related).Refs) == 1 {
				rel = rel.(*Related).Refs[0]
			}

			refs = append(refs, rel)
			continue
		}

		if err != nil {
			ps = refPs
			break
		}

		refs = append(refs, pref)
	}

	m := &Multiple{Refs: refs}
	if err := m.Validate(); err != nil {
		return m, ps, err
	}

	return m, ps, nil
}
