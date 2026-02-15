package ref_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/ref"
)

func TestCalculateRefStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          string
		wantBook       string
		wantChapters   []string
		wantVerses     []string
		wantChCount    int
		wantVerseCount int
	}{
		{
			name:           "single verse",
			input:          "John 3:16",
			wantBook:       "John",
			wantChapters:   []string{"3"},
			wantVerses:     []string{"3:16"},
			wantChCount:    1,
			wantVerseCount: 1,
		},
		{
			name:           "verse range same chapter",
			input:          "John 3:16-21",
			wantBook:       "John",
			wantChapters:   []string{"3"},
			wantVerses:     []string{"3:16-21"},
			wantChCount:    1,
			wantVerseCount: 6,
		},
		{
			name:           "verse range across chapters",
			input:          "Genesis 1:31-2:3",
			wantBook:       "Genesis",
			wantChapters:   []string{"1-2"},
			wantVerses:     []string{"1:31-2:3"},
			wantChCount:    2,
			wantVerseCount: 4,
		},
		{
			name:           "single chapter",
			input:          "Psalm 23",
			wantBook:       "Psalms",
			wantChapters:   []string{"23"},
			wantVerses:     []string{"23:1-6"},
			wantChCount:    1,
			wantVerseCount: 6,
		},
		{
			name:           "chapter range",
			input:          "Genesis 1-2",
			wantBook:       "Genesis",
			wantChapters:   []string{"1-2"},
			wantVerses:     []string{"1:1-2:25"},
			wantChCount:    2,
			wantVerseCount: 56,
		},
		{
			name:           "single chapter book",
			input:          "Philemon 5",
			wantBook:       "Philemon",
			wantChapters:   nil,
			wantVerses:     []string{"5"},
			wantChCount:    0,
			wantVerseCount: 1,
		},
		{
			name:           "single chapter book range",
			input:          "Philemon 5-10",
			wantBook:       "Philemon",
			wantChapters:   nil,
			wantVerses:     []string{"5-10"},
			wantChCount:    0,
			wantVerseCount: 6,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parsed, err := ref.ParseProper(tt.input)
			require.NoError(t, err)

			resolved, err := ref.Canonical.Resolve(parsed)
			require.NoError(t, err)

			resolvedPtrs := make([]*ref.Resolved, len(resolved))
			for i := range resolved {
				resolvedPtrs[i] = &resolved[i]
			}

			stats, err := ref.CalculateRefStats(resolvedPtrs)
			require.NoError(t, err)

			assert.Equal(t, tt.wantBook, stats.Book)
			assert.Equal(t, tt.wantChapters, stats.ChapterRanges)
			assert.Equal(t, tt.wantVerses, stats.VerseRanges)
			assert.Equal(t, tt.wantChCount, stats.ChapterCount)
			assert.Equal(t, tt.wantVerseCount, stats.VerseCount)
		})
	}
}

func TestCalculateRefStats_MultipleReferences(t *testing.T) {
	t.Parallel()

	parsed, err := ref.ParseMultiple("John 3:16; John 3:18-20")
	require.NoError(t, err)

	resolved, err := ref.Canonical.Resolve(parsed)
	require.NoError(t, err)

	resolvedPtrs := make([]*ref.Resolved, len(resolved))
	for i := range resolved {
		resolvedPtrs[i] = &resolved[i]
	}

	stats, err := ref.CalculateRefStats(resolvedPtrs)
	require.NoError(t, err)

	assert.Equal(t, "John", stats.Book)
	assert.Equal(t, []string{"3"}, stats.ChapterRanges)
	assert.Equal(t, []string{"3:16", "3:18-20"}, stats.VerseRanges)
	assert.Equal(t, 1, stats.ChapterCount)
	assert.Equal(t, 4, stats.VerseCount)
}

func TestCalculateRefStats_EmptyInput(t *testing.T) {
	t.Parallel()

	stats, err := ref.CalculateRefStats([]*ref.Resolved{})
	require.NoError(t, err)

	assert.Equal(t, "", stats.Book)
	assert.Equal(t, []string(nil), stats.ChapterRanges)
	assert.Equal(t, []string(nil), stats.VerseRanges)
	assert.Equal(t, 0, stats.ChapterCount)
	assert.Equal(t, 0, stats.VerseCount)
}

func TestCalculateRefStats_MultiBookError(t *testing.T) {
	t.Parallel()

	// Parse a reference with multiple books
	parsed, err := ref.ParseMultiple("John 3:16; Romans 8:28")
	require.NoError(t, err)

	resolved, err := ref.Canonical.Resolve(parsed)
	require.NoError(t, err)

	resolvedPtrs := make([]*ref.Resolved, len(resolved))
	for i := range resolved {
		resolvedPtrs[i] = &resolved[i]
	}

	// Should return an error because references are from different books
	stats, err := ref.CalculateRefStats(resolvedPtrs)
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "all references must be from the same book")
}

// mockVerseTextResolver implements VerseTextResolver for testing.
type mockVerseTextResolver struct {
	text string
	err  error
}

func (m *mockVerseTextResolver) VerseText(_ context.Context, _ *ref.Resolved) (string, error) {
	return m.text, m.err
}

func TestCalculateESVStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          string
		mockText       string
		wantParagraphs int
		wantLines      int
		wantWords      int
		wantChars      int
	}{
		{
			name:           "single verse single paragraph",
			input:          "John 3:16",
			mockText:       "For God so loved the world, that he gave his only Son.",
			wantParagraphs: 1,
			wantLines:      1,
			wantWords:      12,
			wantChars:      54,
		},
		{
			name:           "multiple paragraphs",
			input:          "Psalm 23",
			mockText:       "The LORD is my shepherd; I shall not want.\n\nHe makes me lie down in green pastures.",
			wantParagraphs: 2,
			wantLines:      3,
			wantWords:      17,
			wantChars:      83,
		},
		{
			name:           "multiple lines single paragraph",
			input:          "Genesis 1:1",
			mockText:       "In the beginning, God created\nthe heavens and the earth.",
			wantParagraphs: 1,
			wantLines:      2,
			wantWords:      10,
			wantChars:      56,
		},
		{
			name:           "empty text",
			input:          "John 3:16",
			mockText:       "",
			wantParagraphs: 0,
			wantLines:      0,
			wantWords:      0,
			wantChars:      0,
		},
		{
			name:           "whitespace only",
			input:          "John 3:16",
			mockText:       "   \n\n   ",
			wantParagraphs: 2,
			wantLines:      3,
			wantWords:      0,
			wantChars:      8,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parsed, err := ref.ParseProper(tt.input)
			require.NoError(t, err)

			resolved, err := ref.Canonical.Resolve(parsed)
			require.NoError(t, err)

			resolvedPtrs := make([]*ref.Resolved, len(resolved))
			for i := range resolved {
				resolvedPtrs[i] = &resolved[i]
			}

			mockResolver := &mockVerseTextResolver{text: tt.mockText}

			stats, err := ref.CalculateESVStats(context.Background(), resolvedPtrs, mockResolver)
			require.NoError(t, err)

			assert.Equal(t, tt.wantParagraphs, stats.Paragraphs)
			assert.Equal(t, tt.wantLines, stats.Lines)
			assert.Equal(t, tt.wantWords, stats.Words)
			assert.Equal(t, tt.wantChars, stats.Characters)

			// Verify RefStats fields are also populated
			assert.NotEmpty(t, stats.Book)
			assert.NotNil(t, stats.VerseRanges)
		})
	}
}

func TestCalculateESVStats_MultipleReferences(t *testing.T) {
	t.Parallel()

	// Use references from the same book
	parsed, err := ref.ParseMultiple("John 3:16; John 3:18")
	require.NoError(t, err)

	resolved, err := ref.Canonical.Resolve(parsed)
	require.NoError(t, err)

	resolvedPtrs := make([]*ref.Resolved, len(resolved))
	for i := range resolved {
		resolvedPtrs[i] = &resolved[i]
	}

	// Mock returns the same text for both references
	// They will be concatenated with \n\n between them
	simpleResolver := &mockVerseTextResolver{text: "Sample verse text."}

	stats, err := ref.CalculateESVStats(context.Background(), resolvedPtrs, simpleResolver)
	require.NoError(t, err)

	// Two verses concatenated with \n\n between them
	// "Sample verse text.\n\nSample verse text."
	assert.Equal(t, 2, stats.Paragraphs) // Two paragraphs separated by \n\n
	assert.Equal(t, 3, stats.Lines)      // 3 lines (text, blank, text)
	assert.Greater(t, stats.Words, 0)
	assert.Greater(t, stats.Characters, 0)
}

func TestCalculateESVStats_Error(t *testing.T) {
	t.Parallel()

	parsed, err := ref.ParseProper("John 3:16")
	require.NoError(t, err)

	resolved, err := ref.Canonical.Resolve(parsed)
	require.NoError(t, err)

	resolvedPtrs := make([]*ref.Resolved, len(resolved))
	for i := range resolved {
		resolvedPtrs[i] = &resolved[i]
	}

	mockResolver := &mockVerseTextResolver{err: assert.AnError}

	stats, err := ref.CalculateESVStats(context.Background(), resolvedPtrs, mockResolver)
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "failed to fetch verse text")
}
