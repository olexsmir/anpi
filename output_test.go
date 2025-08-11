package main

import (
	"bytes"
	"strings"
	"testing"

	"olexsmir.xyz/x/is"
)

func TestToTSV_withExample(t *testing.T) {
	entry := &Entry{
		Word: "skip",
		POSBlocks: []POSBlock{
			{
				POS: "verb",
				IPA: "skɪp",
				Senses: []Sense{
					{
						Definition: "to move lightly and quickly",
						Examples:   []string{"she loves to skip down the hall"},
					},
				},
			},
		},
	}
	is.Equal(t, "she loves to <b>skip</b> down the hall\tto move lightly and quickly\t/skɪp/", toTSV(entry))
}

func TestToTSV_withoutExample(t *testing.T) {
	entry := &Entry{
		Word: "hello",
		POSBlocks: []POSBlock{
			{
				POS: "exclamation",
				IPA: "/heˈləʊ/",
				Senses: []Sense{
					{
						Definition: "used when meeting someone",
					},
				},
			},
		},
	}
	is.Equal(t, "<b>hello</b>\tused when meeting someone\t/heˈləʊ/", toTSV(entry))

	is.Equal(t, "", toTSV(&Entry{}))

	is.Equal(t, "", toTSV(&Entry{POSBlocks: []POSBlock{{}}}))
}

func TestToTSV_ipaWrapping(t *testing.T) {
	entry := &Entry{
		Word: "test",
		POSBlocks: []POSBlock{
			{
				IPA: "/test/",
				Senses: []Sense{
					{Definition: "a test"},
				},
			},
		},
	}
	is.Equal(t, "<b>test</b>\ta test\t/test/", toTSV(entry))
}

func TestToTSV_usIPA(t *testing.T) {
	entry := &Entry{
		Word: "water",
		POSBlocks: []POSBlock{
			{
				IPA: "ˈwɑː.t̬ɚ",
				Senses: []Sense{
					{Definition: "a clear liquid"},
				},
			},
		},
	}
	is.Equal(t, "<b>water</b>\ta clear liquid\t/ˈwɑː.t̬ɚ/", toTSV(entry))
}

func TestBoldWord_caseInsensitive(t *testing.T) {
	is.Equal(t, "she loves to <b>skip</b> down", boldWord("she loves to skip down", "skip"))
}

func TestBoldWord_inflectedExample(t *testing.T) {
	is.Equal(t, "she <b>skipped</b> down", boldWord("she skipped down", "skip"))
}

func TestBoldWord_mixedCase(t *testing.T) {
	is.Equal(t, "<b>Hello</b> there", boldWord("Hello there", "hello"))
}

func TestBoldWord_noMatch(t *testing.T) {
	is.Equal(t, "just passing through", boldWord("just passing through", "skip"))
}

func TestBoldWord_notSubstring(t *testing.T) {
	is.Equal(t, "<b>skipping</b> class", boldWord("skipping class", "skip"))
}

func TestBoldWord_matchesAllInstances(t *testing.T) {
	is.Equal(t, "<b>wait</b> <b>wait</b> <b>wait</b>", boldWord("wait wait wait", "wait"))
}

func TestBoldWord_notSubstringOfLargerWord(t *testing.T) {
	is.Equal(t, "<b>waiter</b>", boldWord("waiter", "wait"))
}

func TestBoldWord_inflectedLookup(t *testing.T) {
	is.Equal(t, "A suspect has been <b>detained</b> by the police", boldWord("A suspect has been detained by the police", "detaining"))
	is.Equal(t, "she <b>skipped</b> down the hall", boldWord("she skipped down the hall", "skipping"))
	is.Equal(t, "he <b>passes</b> the ball", boldWord("he passes the ball", "passing"))
	is.Equal(t, "she <b>liked</b> it", boldWord("she liked it", "liking"))
	is.Equal(t, "he <b>walks</b> fast", boldWord("he walks fast", "walking"))
	is.Equal(t, "we <b>loved</b> it", boldWord("we loved it", "loving"))
}

func TestWriteAll(t *testing.T) {
	entry := &Entry{
		Word: "skip",
		POSBlocks: []POSBlock{
			{
				POS: "verb",
				IPA: "skɪp",
				Senses: []Sense{
					{Definition: "move with quick light steps", Examples: []string{"she skipped down the path", "skip over the puddle"}},
					{Definition: "omit or pass over", Examples: []string{"skip the boring parts"}},
				},
			},
			{
				POS: "noun",
				IPA: "skɪp",
				Senses: []Sense{
					{Definition: "a light bouncing movement", Examples: []string{"a skip of excitement"}},
				},
			},
		},
	}

	var buf bytes.Buffer
	writeAll(entry, &buf)

	got := buf.String()
	is.Equal(t, true, strings.Contains(got, "[verb] /skɪp/"))
	is.Equal(t, true, strings.Contains(got, "[noun] /skɪp/"))
	is.Equal(t, true, strings.Contains(got, "1. move with quick light steps"))
	is.Equal(t, true, strings.Contains(got, "2. omit or pass over"))
	is.Equal(t, true, strings.Contains(got, "   - she skipped down the path"))
	is.Equal(t, true, strings.Contains(got, "   - skip over the puddle"))
	is.Equal(t, true, strings.Contains(got, "   - skip the boring parts"))
	is.Equal(t, true, strings.Contains(got, "1. a light bouncing movement"))
	is.Equal(t, true, strings.Contains(got, "   - a skip of excitement"))
}

func TestWriteAll_emptyEntry(t *testing.T) {
	var buf bytes.Buffer
	writeAll(&Entry{}, &buf)
	is.Equal(t, "", buf.String())
}

func TestWriteAll_noSenses(t *testing.T) {
	var buf bytes.Buffer
	writeAll(&Entry{POSBlocks: []POSBlock{{POS: "verb", IPA: "skɪp"}}}, &buf)
	is.Equal(t, "[verb] /skɪp/\n", buf.String())
}

func TestWriteAll_separatesBlocksWithBlankLine(t *testing.T) {
	entry := &Entry{
		Word: "test",
		POSBlocks: []POSBlock{
			{POS: "verb", IPA: "v", Senses: []Sense{{Definition: "def1"}}},
			{POS: "noun", IPA: "n", Senses: []Sense{{Definition: "def2"}}},
		},
	}
	var buf bytes.Buffer
	writeAll(entry, &buf)
	t.Logf("got:\n%s\n", buf.String())
	is.Equal(t, true, strings.Contains(buf.String(), "[verb] /v/\n  1. def1\n\n[noun] /n/"))
}
