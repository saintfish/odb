// A library to fetch Our Daily Bread online and parse it
package odb

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/bible.go/bible"
	"github.com/saintfish/brave"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// The metadata of a odb page
type Post struct {
	Year, Month, Day int
	Url              string
	Title            string
	BibleVerse       string
	BibleVerseRef    bible.RefRangeList
	GoldenVerse      string
	Passages         []string
	Poem             []string
	Thought          string
}

// Language of odb website
type Language int

const (
	English Language = iota
	SimplifiedChinese
	TraditionalChinese
)

var languageDomainMap = map[Language]string{
	English:            "odb.org",
	SimplifiedChinese:  "simplified-odb.org",
	TraditionalChinese: "traditional-odb.org",
}

// Odb represents an odb website accessor
type Odb struct {
	domain     string
	httpClient *http.Client
}

// Create an accessor to a odb website in a specific language
func NewOdb(l Language) *Odb {
	return &Odb{languageDomainMap[l], &http.Client{}}
}

func NewOdbWithHttpClient(l Language, c *http.Client) *Odb {
	return &Odb{languageDomainMap[l], c}
}

func (odb *Odb) ListPost(year, month int) (map[int]string, error) {
	dayUrlMap := make(map[int]string)
	url := fmt.Sprintf("http://%s/%d/%02d/", odb.domain, year, month)
	q, err := fetch(odb.httpClient, url)
	if err != nil {
		return nil, fmt.Errorf("Error in crawl list page %s: #v", url, err)
	}
	links := q.Find("#wp-calendar a")
	links.Each(func(i int, link *goquery.Selection) {
		dayStr := strings.TrimSpace(link.Text())
		href, _ := link.Attr("href")
		if n, err := strconv.ParseInt(dayStr, 10, 8); err == nil {
			dayUrlMap[int(n)] = href
		}
	})
	return dayUrlMap, nil
}

// Get a post in specific date
func (odb *Odb) GetPost(year, month, day int) (*Post, error) {
	dayUrlMap, err := odb.ListPost(year, month)
	if err != nil {
		return nil, err
	}
	url, found := dayUrlMap[day]
	if !found {
		return nil, fmt.Errorf("cannot find %d/%d/%d", year, month, day)
	}
	q, err := fetch(odb.httpClient, url)
	if err != nil {
		return nil, fmt.Errorf("error in crawl post page %s: #v", url, err)
	}
	title := q.Find(".entry-title").Text()
	metaBoxes := q.Find(".entry-content .side-box .meta-box")
	bibleVerse := metaBoxes.Eq(0).Text()
	goldenVerse := metaBoxes.Eq(1).Text()
	passages := []string{}
	q.Find(".entry-content > p").Each(func(_ int, s *goquery.Selection) {
		passages = append(passages, s.Text())
	})
	poemText, _ := q.Find(".entry-content .poem-box").Html()
	poem := splitPassages(poemText)
	thought := q.Find(".entry-content .thought-box").Text()

	refMatch, bibleVerseRef := brave.ParseChineseFull(bibleVerse)
	if !refMatch {
		bibleVerseRef = nil
	}

	p := &Post{
		Year:          year,
		Month:         month,
		Day:           day,
		Url:           url,
		Title:         title,
		BibleVerse:    bibleVerse,
		BibleVerseRef: bibleVerseRef,
		GoldenVerse:   goldenVerse,
		Passages:      passages,
		Poem:          poem,
		Thought:       thought,
	}
	return p, nil
}

var newLine = regexp.MustCompile("[\\n\\r]+|<br\\s*/?>")

func splitPassages(html string) []string {
	result := []string{}
	for _, s := range newLine.Split(html, -1) {
		if len(s) != 0 {
			result = append(result, s)
		}
	}
	return result
}
