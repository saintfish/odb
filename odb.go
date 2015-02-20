// A library to fetch Our Daily Bread online and parse it
package odb

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/bible.go/bible"
	"github.com/saintfish/brave"
	"net/http"
	"regexp"
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

type param struct {
	domain                        string
	bibleBoxPrefix, bibleBoxSplit string
}

var languageParamMap = map[Language]param{
	English: param{
		domain:         "odb.org",
		bibleBoxPrefix: "Read: ",
		bibleBoxSplit:  " | Bible in a Year: ",
	},
	SimplifiedChinese: param{
		domain:         "simplified-odb.org",
		bibleBoxPrefix: "读经: ",
		bibleBoxSplit:  " | 全年读经: ",
	},
	TraditionalChinese: param{
		domain:         "traditional-odb.org",
		bibleBoxPrefix: "讀經: ",
		bibleBoxSplit:  " | 全年讀經: ",
	},
}

// Odb represents an odb website accessor
type Odb struct {
	param
	httpClient *http.Client
}

// Create an accessor to a odb website in a specific language
func NewOdb(l Language) *Odb {
	return &Odb{languageParamMap[l], &http.Client{}}
}

func NewOdbWithHttpClient(l Language, c *http.Client) *Odb {
	return &Odb{languageParamMap[l], c}
}

// Get a post in specific date
func (odb *Odb) GetPost(year, month, day int) (*Post, error) {
	calendarUrl := fmt.Sprintf("http://%s/%4d/%02d/%02d?calendar-redirect=true", odb.domain, year, month, day)
	q, err := fetch(odb.httpClient, calendarUrl)
	if err != nil {
		return nil, fmt.Errorf("error in crawl post page %s: #v", calendarUrl, err)
	}
	title := text(q.Find("h2.entry-title"))
	goldenVerse := text(q.Find(".verse-box"))
	bibleBox := text(q.Find(".passage-box"))
	bibleVerse := ""
	if strings.HasPrefix(bibleBox, odb.bibleBoxPrefix) {
		bibleVerse, _ = token(bibleBox[len(odb.bibleBoxPrefix):], odb.bibleBoxSplit)
	}
	passages := []string{}
	q.Find(".post-content > p").Each(func(_ int, s *goquery.Selection) {
		passages = append(passages, text(s))
	})
	poemText, _ := q.Find(".entry-content .poem-box").Html()
	poem := splitParagraphs(poemText)
	thought := text(q.Find(".entry-content .thought-box"))

	refMatch, bibleVerseRef := brave.ParseChineseFull(bibleVerse)
	if !refMatch {
		bibleVerseRef = nil
	}

	p := &Post{
		Year:           year,
		Month:          month,
		Day:            day,
		Url:            q.Url.String(),
		Title:          title,
		BibleVerse:     bibleVerse,
		BibleVerseRef:  bibleVerseRef,
		GoldenVerse:    goldenVerse,
		Passages:       passages,
		Poem:           poem,
		Thought:        thought,
	}
	return p, nil
}

func fetch(c *http.Client, url string) (*goquery.Document, error) {
	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unsuccessful status code %d", res.StatusCode)
	}
	q, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	q.Url = res.Request.URL
	return q, nil
}

func text(q *goquery.Selection) string {
	return strings.Trim(q.Text(), "\u2029 \n\t")
}

func token(s, sep string) (token, remainder string) {
	split := strings.SplitN(s, sep, 2)
	if len(split) < 2 {
		return "", split[0]
	}
	return split[0], split[1]
}

var newLine = regexp.MustCompile("((\\s|\u2029)*<br\\s*/?>(\\s|\u2029)*)+")

func splitParagraphs(html string) []string {
	result := []string{}
	for _, s := range newLine.Split(html, -1) {
		if len(s) != 0 {
			result = append(result, s)
		}
	}
	return result
}
