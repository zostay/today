package ref_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/ref"
)

func TestGetAvailableStyles(t *testing.T) {
	t.Parallel()

	styles := ref.GetAvailableStyles()
	assert.Equal(t, []string{
		"canonical",
		"abbr",
		"2letter",
		"3letter",
		"2letter.",
		"3letter.",
	}, styles)
}

func TestGetFormatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		style     string
		wantError bool
	}{
		{"canonical style", "canonical", false},
		{"abbr style", "abbr", false},
		{"2letter style", "2letter", false},
		{"3letter style", "3letter", false},
		{"2letter. style", "2letter.", false},
		{"3letter. style", "3letter.", false},
		{"invalid style", "invalid", true},
		{"empty style", "", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			formatter, err := ref.GetFormatter(tt.style)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, formatter)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, formatter)
			}
		})
	}
}

func TestCanonicalFormatter_Format(t *testing.T) {
	t.Parallel()

	formatter, err := ref.GetFormatter("canonical")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single verse", "John 3:16", "John 3:16"},
		{"verse range", "John 3:16-18", "John 3:16-18"},
		{"chapter", "Psalm 23", "Psalm 23"},
		{"chapter range", "Genesis 1-2", "Genesis 1-2"},
		{"multiple verses", "John 3:16; Romans 8:28", "John 3:16; Romans 8:28"},
		{"numbered book", "1 John 3:16", "1 John 3:16"},
		{"single chapter book", "Philemon 5", "Philemon 5"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var parsed ref.Absolute
			var err error
			parsed, err = ref.ParseProper(tt.input)
			if err != nil {
				parsed, err = ref.ParseMultiple(tt.input)
				require.NoError(t, err)
			}

			resolved, err := ref.Canonical.Resolve(parsed)
			require.NoError(t, err)

			resolvedPtrs := make([]*ref.Resolved, len(resolved))
			for i := range resolved {
				resolvedPtrs[i] = &resolved[i]
			}

			result, err := formatter.Format(resolvedPtrs)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAbbrFormatter_Format(t *testing.T) {
	t.Parallel()

	formatter, err := ref.GetFormatter("abbr")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single verse", "Genesis 1:1", "Gen. 1:1"},
		{"verse range", "Romans 8:28-39", "Rom. 8:28-39"},
		{"chapter", "Psalm 23", "Ps. 23"},
		{"multiple verses", "Genesis 1:1; Romans 8:28", "Gen. 1:1; Rom. 8:28"},
		{"numbered book", "1 John 3:16", "1 John 3:16"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var parsed ref.Absolute
			var err error
			parsed, err = ref.ParseProper(tt.input)
			if err != nil {
				parsed, err = ref.ParseMultiple(tt.input)
				require.NoError(t, err)
			}

			resolved, err := ref.Canonical.Resolve(parsed)
			require.NoError(t, err)

			resolvedPtrs := make([]*ref.Resolved, len(resolved))
			for i := range resolved {
				resolvedPtrs[i] = &resolved[i]
			}

			result, err := formatter.Format(resolvedPtrs)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNLetterFormatter_Format(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		style    string
		input    string
		expected string
	}{
		{"2letter Genesis", "2letter", "Genesis 1:1", "Gn 1:1"},
		{"3letter Genesis", "3letter", "Genesis 1:1", "Gen 1:1"},
		{"2letter. Genesis", "2letter.", "Genesis 1:1", "Gn. 1:1"},
		{"3letter. Genesis", "3letter.", "Genesis 1:1", "Gen. 1:1"},
		{"2letter Exodus", "2letter", "Exodus 1:1", "Ex 1:1"},
		{"3letter Exodus", "3letter", "Exodus 1:1", "Exo 1:1"},
		{"2letter numbered book", "2letter", "1 John 3:16", "1 Jn 3:16"},
		{"3letter numbered book", "3letter", "1 John 3:16", "1 Jhn 3:16"},
		{"2letter. numbered book", "2letter.", "2 Peter 1:1", "2 Pt. 1:1"},
		{"3letter. numbered book", "3letter.", "2 Peter 1:1", "2 Pet. 1:1"},
		{"2letter multiple", "2letter", "Genesis 1:1; Romans 8:28", "Gn 1:1; Rm 8:28"},
		{"2letter verse range", "2letter", "John 3:16-18", "Jn 3:16-18"},
		{"3letter chapter", "3letter", "Psalm 23", "Psm 23"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			formatter, err := ref.GetFormatter(tt.style)
			require.NoError(t, err)

			var parsed ref.Absolute
			parsed, err = ref.ParseProper(tt.input)
			if err != nil {
				parsed, err = ref.ParseMultiple(tt.input)
				require.NoError(t, err)
			}

			resolved, err := ref.Canonical.Resolve(parsed)
			require.NoError(t, err)

			resolvedPtrs := make([]*ref.Resolved, len(resolved))
			for i := range resolved {
				resolvedPtrs[i] = &resolved[i]
			}

			result, err := formatter.Format(resolvedPtrs)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatResolvedWithName_SingleChapterBook(t *testing.T) {
	t.Parallel()

	formatter, err := ref.GetFormatter("2letter")
	require.NoError(t, err)

	parsed, err := ref.ParseProper("Philemon 5")
	require.NoError(t, err)

	resolved, err := ref.Canonical.Resolve(parsed)
	require.NoError(t, err)

	resolvedPtrs := make([]*ref.Resolved, len(resolved))
	for i := range resolved {
		resolvedPtrs[i] = &resolved[i]
	}

	result, err := formatter.Format(resolvedPtrs)
	require.NoError(t, err)
	assert.Equal(t, "Pm 5", result)
}
