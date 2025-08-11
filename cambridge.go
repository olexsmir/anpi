package main

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

const (
	cambridgeURL = "https://dictionary.cambridge.org/dictionary/english/"
	userAgent    = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type CambridgeDictionary struct {
	baseURL string
	ua      string
}

func NewCambridgeDictionary() *CambridgeDictionary {
	return &CambridgeDictionary{baseURL: cambridgeURL, ua: userAgent}
}

func (d *CambridgeDictionary) fetch(word string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", d.baseURL+strings.ToLower(word), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", d.ua)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("word not found (HTTP %d)", resp.StatusCode)
	}
	return resp.Body, nil
}

func redirectBase(entry *Entry) string {
	if len(entry.POSBlocks) == 0 || len(entry.POSBlocks[0].Senses) == 0 {
		return ""
	}
	def := strings.Join(strings.Fields(entry.POSBlocks[0].Senses[0].Definition), " ")
	prefixes := []string{
		"past simple of ", "past tense of ", "past participle of ",
		"present participle of ", "third person singular of ", "3rd person singular of ",
	}
	for _, p := range prefixes {
		if i := strings.Index(strings.ToLower(def), p); i >= 0 {
			base := strings.TrimSpace(def[i+len(p):])
			if space := strings.IndexAny(base, " \t\n"); space > 0 {
				base = base[:space]
			}
			return strings.ToLower(base)
		}
	}
	return ""
}

func (d *CambridgeDictionary) Lookup(word string) (*Entry, error) {
	body, err := d.fetch(word)
	if err != nil {
		return nil, err
	}
	defer body.Close() //nolint:errcheck
	entry, err := d.Scrape(word, body)
	if err != nil {
		return nil, err
	}
	if base := redirectBase(entry); base != "" && base != word {
		redirected, err := d.Lookup(base)
		if err != nil {
			return nil, err
		}
		redirected.Word = word
		redirected.RedirectWord = base
		return redirected, nil
	}
	return entry, nil
}

func (d *CambridgeDictionary) Scrape(word string, r io.Reader) (*Entry, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	cald4 := findFirstDescendantWithAttr(doc, "div", "data-id", "cald4")
	if cald4 == nil {
		return nil, fmt.Errorf("word not found")
	}

	bodies := findAllDescendantsWithClass(cald4, "div", "entry-body__el")
	if len(bodies) == 0 {
		return nil, fmt.Errorf("word not found")
	}

	var blocks []POSBlock
	for _, body := range bodies {
		blocks = append(blocks, parsePOSBlock(body))
	}
	return &Entry{Word: word, POSBlocks: blocks}, nil
}

func hasClass(n *html.Node, class string) bool {
	for _, a := range n.Attr {
		if a.Key == "class" {
			if slices.Contains(strings.Fields(a.Val), class) {
				return true
			}
		}
	}
	return false
}

func textContent(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.TrimSpace(b.String())
}

func findFirstDescendantWithClass(n *html.Node, tag, class string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag && hasClass(n, class) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findFirstDescendantWithClass(c, tag, class); found != nil {
			return found
		}
	}
	return nil
}

func findAllDescendantsWithClass(n *html.Node, tag, class string) []*html.Node {
	var res []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag && hasClass(n, class) {
			res = append(res, n)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return res
}

func findFirstDescendantWithAttr(n *html.Node, tag, key, val string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		for _, a := range n.Attr {
			if a.Key == key && a.Val == val {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findFirstDescendantWithAttr(c, tag, key, val); found != nil {
			return found
		}
	}
	return nil
}

func parseIPA(ph *html.Node, region string) string {
	el := findFirstDescendantWithClass(ph, "span", region)
	if el == nil {
		return ""
	}
	ipaEl := findFirstDescendantWithClass(el, "span", "ipa")
	if ipaEl == nil {
		return ""
	}
	return textContent(ipaEl)
}

func parsePOSBlock(body *html.Node) POSBlock {
	ph := findFirstDescendantWithClass(body, "div", "pos-header")

	pos := ""
	if el := findFirstDescendantWithClass(ph, "span", "pos"); el != nil {
		pos = textContent(el)
	}

	ipa := parseIPA(ph, "us")
	if ipa == "" {
		ipa = parseIPA(ph, "uk")
	}

	defBlocks := findAllDescendantsWithClass(body, "div", "def-block")
	var senses []Sense
	for _, db := range defBlocks {
		def := ""
		if el := findFirstDescendantWithClass(db, "div", "def"); el != nil {
			def = strings.TrimRight(textContent(el), ":")
		}
		exampDivs := findAllDescendantsWithClass(db, "div", "examp")
		var exs []string
		for i, ed := range exampDivs {
			if i >= 3 {
				break
			}
			if eg := findFirstDescendantWithClass(ed, "span", "eg"); eg != nil {
				exs = append(exs, textContent(eg))
			}
		}
		senses = append(senses, Sense{Definition: def, Examples: exs})
	}

	return POSBlock{POS: pos, IPA: ipa, Senses: senses}
}
