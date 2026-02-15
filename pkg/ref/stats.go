package ref

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// RefStats contains basic statistics about a resolved reference.
type RefStats struct {
	Book          string
	ChapterRanges []string
	VerseRanges   []string
	ChapterCount  int
	VerseCount    int
}

// CalculateRefStats calculates basic statistics for a set of resolved references.
// All references should be from the same book.
func CalculateRefStats(resolved []*Resolved) *RefStats {
	if len(resolved) == 0 {
		return &RefStats{}
	}

	// Get book name from first reference
	bookName := resolved[0].Book.Name
	isJustVerse := resolved[0].Book.JustVerse

	// Collect all chapters and verses
	chapterSet := make(map[int]bool)
	verseCount := 0
	verseRanges := make([]string, 0, len(resolved))

	for _, r := range resolved {
		// Count verses
		verses := r.Verses()
		verseCount += len(verses)

		// Build verse range string for this resolved reference
		if isJustVerse {
			// Single-chapter books (e.g., Philemon)
			first := r.First.(N)
			last := r.Last.(N)
			if first.Number == last.Number {
				verseRanges = append(verseRanges, fmt.Sprintf("%d", first.Number))
			} else {
				verseRanges = append(verseRanges, fmt.Sprintf("%d-%d", first.Number, last.Number))
			}
		} else {
			// Multi-chapter books
			firstCV := r.First.(CV)
			lastCV := r.Last.(CV)

			// Collect chapters
			for ch := firstCV.Chapter; ch <= lastCV.Chapter; ch++ {
				chapterSet[ch] = true
			}

			// Format verse range
			if firstCV.Chapter == lastCV.Chapter {
				if firstCV.Verse == lastCV.Verse {
					verseRanges = append(verseRanges, fmt.Sprintf("%d:%d", firstCV.Chapter, firstCV.Verse))
				} else {
					verseRanges = append(verseRanges, fmt.Sprintf("%d:%d-%d", firstCV.Chapter, firstCV.Verse, lastCV.Verse))
				}
			} else {
				verseRanges = append(verseRanges, fmt.Sprintf("%d:%d-%d:%d", firstCV.Chapter, firstCV.Verse, lastCV.Chapter, lastCV.Verse))
			}
		}
	}

	// Build chapter ranges
	var chapterRanges []string
	if !isJustVerse {
		chapters := make([]int, 0, len(chapterSet))
		for ch := range chapterSet {
			chapters = append(chapters, ch)
		}
		sort.Ints(chapters)

		// Compact consecutive chapters into ranges
		if len(chapters) > 0 {
			start := chapters[0]
			prev := chapters[0]

			for i := 1; i < len(chapters); i++ {
				if chapters[i] == prev+1 {
					prev = chapters[i]
				} else {
					if start == prev {
						chapterRanges = append(chapterRanges, fmt.Sprintf("%d", start))
					} else {
						chapterRanges = append(chapterRanges, fmt.Sprintf("%d-%d", start, prev))
					}
					start = chapters[i]
					prev = chapters[i]
				}
			}

			// Add final range
			if start == prev {
				chapterRanges = append(chapterRanges, fmt.Sprintf("%d", start))
			} else {
				chapterRanges = append(chapterRanges, fmt.Sprintf("%d-%d", start, prev))
			}
		}
	}

	return &RefStats{
		Book:          bookName,
		ChapterRanges: chapterRanges,
		VerseRanges:   verseRanges,
		ChapterCount:  len(chapterSet),
		VerseCount:    verseCount,
	}
}

// ESVStats contains both reference statistics and ESV-specific text statistics.
type ESVStats struct {
	RefStats
	Paragraphs  int
	Lines       int
	Words       int
	Characters  int
}

// VerseTextResolver is an interface for resolving verse text.
// This is used to allow mocking in tests.
type VerseTextResolver interface {
	VerseText(ctx context.Context, ref *Resolved) (string, error)
}

// CalculateESVStats calculates statistics including ESV-specific metrics by
// fetching the verse text from the ESV API.
func CalculateESVStats(ctx context.Context, resolved []*Resolved, resolver VerseTextResolver) (*ESVStats, error) {
	// Calculate basic ref stats
	refStats := CalculateRefStats(resolved)

	// Fetch text for all resolved references and concatenate
	var allText strings.Builder
	for _, r := range resolved {
		text, err := resolver.VerseText(ctx, r)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch verse text: %w", err)
		}
		if allText.Len() > 0 {
			allText.WriteString("\n\n")
		}
		allText.WriteString(text)
	}

	text := allText.String()

	// Count paragraphs (sequences of text separated by blank lines)
	paragraphs := 1
	if len(text) > 0 {
		paragraphs = strings.Count(text, "\n\n") + 1
	} else {
		paragraphs = 0
	}

	// Count lines (newline characters + 1)
	lines := 1
	if len(text) > 0 {
		lines = strings.Count(text, "\n") + 1
	} else {
		lines = 0
	}

	// Count words (split on whitespace)
	words := 0
	if len(strings.TrimSpace(text)) > 0 {
		fields := strings.Fields(text)
		words = len(fields)
	}

	// Count characters (including spaces)
	characters := len(text)

	return &ESVStats{
		RefStats:   *refStats,
		Paragraphs: paragraphs,
		Lines:      lines,
		Words:      words,
		Characters: characters,
	}, nil
}
