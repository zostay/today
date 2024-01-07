package ref

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type expectedRefType int

const (
	expectEither expectedRefType = iota
	expectJustVerse
	expectChapterAndVerse
)

type parseVerseOpts struct {
	expectedRefType
	allowWildcard bool
}

// ParseVerseOption is an option for ParseVerseRef.
type ParseVerseOption func(*parseVerseOpts)

// ExpectJustVerse changes ParseVerseRef so that it expects a verse number only.
// If it encounters chapter and verse (e.g. 3:16), it will return an error. This
// cannot be combined with ExpectChapterAndVerse.
func ExpectJustVerse(opts *parseVerseOpts) {
	opts.expectedRefType = expectJustVerse
}

// ExpectChapterAndVerse changes ParseVerseRef so that it expects a chapter and
// verse number. If it encounters a verse number only (e.g. 16), it will return
// an error. This cannot be combined with ExpectJustVerse.
func ExpectChapterAndVerse(opts *parseVerseOpts) {
	opts.expectedRefType = expectChapterAndVerse
}

// AllowWildcard changes ParseVerseRef so that it allows a wildcard character
// in the verse reference. The wildcard character is "*". If the wildcard
// character is used, the returned VerseRef will have a Wildcard() value of
// WildcardChapter if the chapter is a wildcard, WildcardVerse if the verse is
// a wildcard, or WildcardNone if neither is a wildcard. If the wildcard
// character is used in a chapter reference, the verse must also be wildcarded.
func AllowWildcard(opts *parseVerseOpts) {
	opts.allowWildcard = true
}

// ParseVerseRef parses a verse reference into a VerseRef. The verse reference
// can be a verse number only (e.g. 16), or a chapter and verse number (e.g.
// 3:16). If the verse reference is a chapter and verse number, it must be
// separated by a colon. If the verse reference is a verse number only, it must
// not contain a colon. If the verse reference is a wildcard, it must be "*".
// If the verse reference is a chapter and verse number, the chapter and verse
// must be valid integers. If the verse reference is a verse number only, the
// verse must be a valid integer. If the verse reference is a wildcard, the
// chapter and verse must be "*".
func ParseVerseRef(ref string, opt ...ParseVerseOption) (VerseRef, error) {
	opts := &parseVerseOpts{
		expectedRefType: expectEither,
		allowWildcard:   false,
	}
	for _, o := range opt {
		o(opts)
	}

	return parseVerseRef(ref, opts)
}

func parseVerseRef(ref string, opt *parseVerseOpts) (VerseRef, error) {
	parts := strings.Split(ref, ":")
	if len(parts) == 1 {
		if opt.expectedRefType == expectChapterAndVerse {
			return nil, errors.New("invalid verse reference: expected chapter and verse number")
		}

		if opt.allowWildcard && parts[0] == "*" {
			return &JustVerse{verse: Final}, nil
		}

		verse, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid verse reference: expected a verse number: %w", err)
		}

		return &JustVerse{verse: verse}, nil
	} else if len(parts) == 2 {
		if opt.expectedRefType == expectJustVerse {
			return nil, errors.New("invalid verse reference: expected verse number only")
		}

		if opt.allowWildcard && parts[0] == "*" {
			if parts[1] != "*" {
				return nil, errors.New("invalid verse reference: chapter is wildcard, but verse is not")
			}
			return &ChapterVerse{chapter: Final, verse: Final}, nil
		}

		chapter, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid verse reference: expected a chapter number: %w", err)
		}

		if opt.allowWildcard && parts[1] == "*" {
			return &ChapterVerse{chapter: chapter, verse: Final}, nil
		}

		verse, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid verse reference: expected a verse number: %w", err)
		}

		return &ChapterVerse{chapter: chapter, verse: verse}, nil
	} else {
		return nil, errors.New("invalid verse reference: too many colons")
	}
}
