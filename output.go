package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

func toTSV(entry *Entry) string {
	if len(entry.POSBlocks) == 0 || len(entry.POSBlocks[0].Senses) == 0 {
		return ""
	}
	block := entry.POSBlocks[0]
	sense := block.Senses[0]

	col1 := ""
	if len(sense.Examples) > 0 {
		col1 = boldWord(sense.Examples[0], entry.Word, entry.RedirectWord)
	} else {
		col1 = fmt.Sprintf("<b>%s</b>", entry.Word)
	}
	col2 := sense.Definition
	col3 := block.IPA
	if col3 != "" && !strings.HasPrefix(col3, "/") {
		col3 = "/" + col3 + "/"
	}
	return fmt.Sprintf("%s\t%s\t%s", col1, col2, col3)
}

func boldWord(s, word string, extra ...string) string {
	words := append([]string{word}, extra...)
	seen := map[string]bool{}
	for _, w := range words {
		if w == "" || seen[w] {
			continue
		}
		seen[w] = true
		stem := w
		for _, suffix := range []string{"ing", "ed", "es", "d", "s"} {
			if len(w) > len(suffix) && strings.HasSuffix(w, suffix) {
				stem = strings.TrimSuffix(w, suffix)
				break
			}
		}
		if len(stem) < 3 {
			stem = w
		}
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(stem) + `\w*\b`)
		if result := re.ReplaceAllString(s, `<b>$0</b>`); result != s {
			return result
		}
	}
	return s
}

func writeAll(entry *Entry, w io.Writer) {
	for pi, block := range entry.POSBlocks {
		if pi > 0 {
			_, _ = fmt.Fprintln(w)
		}
		_, _ = fmt.Fprintf(w, "[%s] /%s/\n", block.POS, block.IPA)
		for si, sense := range block.Senses {
			if si > 0 {
				_, _ = fmt.Fprintln(w)
			}
			_, _ = fmt.Fprintf(w, "  %d. %s\n", si+1, sense.Definition)
			for _, ex := range sense.Examples {
				_, _ = fmt.Fprintf(w, "     - %s\n", ex)
			}
		}
	}
}
