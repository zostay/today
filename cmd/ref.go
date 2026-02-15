package cmd

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text/esv"
)

var refCmd = &cobra.Command{
	Use:   "ref [references...]",
	Short: "Convert Bible references to various output styles",
	Long: `Convert Bible references to various output styles.

Takes references from command-line arguments or from standard input (one per line).
Outputs formatted references according to the specified style.

Available styles:
  canonical - Full book name (e.g., "John 3:16")
  abbr      - Preferred abbreviation (e.g., "Jn. 3:16")
  2letter   - First 2-letter abbreviation (e.g., "Jn 3:16")
  3letter   - First 3-letter abbreviation (e.g., "Jhn 3:16")
  2letter.  - First 2-letter abbreviation with period (e.g., "Jn. 3:16")
  3letter.  - First 3-letter abbreviation with period (e.g., "Jhn. 3:16")`,
	Args: cobra.ArbitraryArgs,
	RunE: RunRef,
}

var (
	refStyle      string
	refListStyles bool
	refStat       string
)

func init() {
	refCmd.Flags().StringVarP(&refStyle, "style", "s", "canonical", "Output style for references")
	refCmd.Flags().BoolVar(&refListStyles, "list-styles", false, "List available styles and exit")
	refCmd.Flags().StringVar(&refStat, "stat", "off", "Show statistics (off|ref|esv)")
	refCmd.Flags().Lookup("stat").NoOptDefVal = "ref"
}

func RunRef(cmd *cobra.Command, args []string) error {
	// Handle --list-styles
	if refListStyles {
		for _, style := range ref.GetAvailableStyles() {
			fmt.Fprintln(cmd.OutOrStdout(), style)
		}
		return nil
	}

	// Get formatter
	formatter, err := ref.GetFormatter(refStyle)
	if err != nil {
		return fmt.Errorf("invalid style: %w", err)
	}

	// Determine input source
	var references []string
	if len(args) > 0 {
		references = args
	} else {
		scanner := bufio.NewScanner(cmd.InOrStdin())
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				references = append(references, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	// Process each reference
	for _, refStr := range references {
		if err := processReference(cmd, formatter, refStr); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error processing %q: %v\n", refStr, err)
			continue
		}
	}

	return nil
}

func processReference(cmd *cobra.Command, formatter ref.RefFormatter, refStr string) error {
	// Try parsing as Proper first, then Multiple
	var parsed ref.Absolute
	var err error

	parsed, err = ref.ParseProper(refStr)
	if err != nil {
		parsed, err = ref.ParseMultiple(refStr)
		if err != nil {
			return fmt.Errorf("failed to parse: %w", err)
		}
	}

	// Validate
	if err := parsed.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Resolve
	resolved, err := ref.Canonical.Resolve(parsed)
	if err != nil {
		return fmt.Errorf("resolution failed: %w", err)
	}

	// Convert []Resolved to []*Resolved
	resolvedPtrs := make([]*ref.Resolved, len(resolved))
	for i := range resolved {
		resolvedPtrs[i] = &resolved[i]
	}

	// Format
	formatted, err := formatter.Format(resolvedPtrs)
	if err != nil {
		return fmt.Errorf("formatting failed: %w", err)
	}

	// Output formatted reference
	fmt.Fprintln(cmd.OutOrStdout(), formatted)

	// Output stats if requested
	switch refStat {
	case "off":
		// No stats
	case "ref":
		stats := ref.CalculateRefStats(resolvedPtrs)
		printRefStats(cmd, stats)
	case "esv":
		ec, err := esv.NewFromEnvironment()
		if err != nil {
			return fmt.Errorf("failed to initialize ESV client: %w", err)
		}
		stats, err := ref.CalculateESVStats(cmd.Context(), resolvedPtrs, ec)
		if err != nil {
			return fmt.Errorf("failed to calculate ESV stats: %w", err)
		}
		printESVStats(cmd, stats)
	default:
		return fmt.Errorf("invalid stat mode: %q (expected off, ref, or esv)", refStat)
	}

	return nil
}

func printRefStats(cmd *cobra.Command, stats *ref.RefStats) {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "  Book: %s\n", stats.Book)

	if len(stats.ChapterRanges) > 0 {
		if len(stats.ChapterRanges) == 1 {
			fmt.Fprintf(out, "  Chapter: %s\n", stats.ChapterRanges[0])
		} else {
			fmt.Fprintf(out, "  Chapters: %s\n", stats.ChapterRanges[0])
			for _, ch := range stats.ChapterRanges[1:] {
				fmt.Fprintf(out, "            %s\n", ch)
			}
		}
	}

	fmt.Fprintln(out, "  Verse Ranges:")
	for _, vr := range stats.VerseRanges {
		fmt.Fprintf(out, "    %s\n", vr)
	}

	fmt.Fprintf(out, "  Chapter Count: %d\n", stats.ChapterCount)
	fmt.Fprintf(out, "  Verse Count: %d\n", stats.VerseCount)
}

func printESVStats(cmd *cobra.Command, stats *ref.ESVStats) {
	// Print basic ref stats first
	printRefStats(cmd, &stats.RefStats)

	// Print ESV-specific stats
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "  Paragraphs: %d\n", stats.Paragraphs)
	fmt.Fprintf(out, "  Lines: %d\n", stats.Lines)
	fmt.Fprintf(out, "  Words: %d\n", stats.Words)
	fmt.Fprintf(out, "  Characters: %d\n", stats.Characters)
}
