// A library to fetch Our Daily Bread online and parse it
package odb

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

// The metadata of a odb page
type Post struct {
	Year, Month, Day int
	Url              string
	Title            string
	BibleVerse       string
	GoldenVerse      string
	Passages         []string
	Poem             string
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
	domain string
}

// Create an accessor to a odb website in a specific language
func NewOdb(l Language) *Odb {
	return &Odb{languageDomainMap[l]}
}

// Get a post in specific date
func (odb *Odb) GetPost(year, month, day int) (*Post, error) {
	q, err := fetch(fmt.Sprintf("http://%s/%d/%02d/%02d/", odb.domain, year, month, day))
	if err != nil {
		return nil, fmt.Errorf("Error in crawl list page: #v", err)
	}
	link := q.Find(".entry-title > a")
	url, _ := link.Attr("href")
	title := link.Text()
	if url == "" {
		return nil, errors.New("Unable to find the post url")
	}
	q, err = fetch(url)
	if err != nil {
		return nil, fmt.Errorf("Error in crawl post page: #v", err)
	}
	metaBoxes := q.Find(".entry-content .side-box .meta-box")
	bibleVerse := metaBoxes.Eq(0).Text()
	goldenVerse := metaBoxes.Eq(1).Text()
	passages := []string{}
	q.Find(".entry-content > p").Each(func(_ int, s *goquery.Selection) {
		passages = append(passages, s.Text())
	})
	poem, _ := q.Find(".entry-content .poem-box").Html()
	thought := q.Find(".entry-content .thought-box").Text()

	p := &Post{
		Year:        year,
		Month:       month,
		Day:         day,
		Url:         url,
		Title:       title,
		BibleVerse:  bibleVerse,
		GoldenVerse: goldenVerse,
		Passages:    passages,
		Poem:        poem,
		Thought:     thought,
	}
	return p, nil
}
