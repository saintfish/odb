package odb

import (
	"testing"
	"time"
)

func getPost(t *testing.T, l Language, year, month, day int) *Post {
	odb := NewOdb(l)
	p, err := odb.GetPost(year, month, day)
	if err != nil {
		t.Errorf("Error in NewOdb(%v).GetPost(%d, %d, %d): %v", l, year, month, day, err)
	}
	return p
}

func forDate(t *testing.T, p *Post, year, month, day int) {
	if p.Year != year || p.Month != month || p.Day != day {
		t.Errorf("Post is not for date %d/%d/%d: %#v", year, month, day, p)
	}
}

func hasData(t *testing.T, p *Post) bool {
	if len(p.Url) == 0 ||
		len(p.Title) == 0 ||
		len(p.BibleVerse) == 0 ||
		len(p.GoldenVerse) == 0 ||
		len(p.Passages) == 0 ||
		len(p.Poem) == 0 ||
		len(p.Thought) == 0 {
		t.Errorf("Post has incomplete data: %#v", p)
		return false
	}
	return true
}

func hasBibleRef(t *testing.T, p *Post) bool {
	if len(p.BibleVerseRef) == 0 {
		t.Errorf("Unable to detect verse ref from %s", p.BibleVerse)
		return false
	}
	return true
}

func TestToday(t *testing.T) {
	y, m, d := time.Now().Date()
	for _, l := range []Language{English, SimplifiedChinese, TraditionalChinese} {
		p := getPost(t, l, y, int(m), d)
		if p != nil {
			forDate(t, p, y, int(m), d)
			hasData(t, p)
			if l == SimplifiedChinese || l == TraditionalChinese {
				hasBibleRef(t, p)
			}
		}
	}
}

func Test7DaysLater(t *testing.T) {
	y, m, d := time.Now().AddDate(0, 0, 7).Date()
	for _, l := range []Language{English, SimplifiedChinese, TraditionalChinese} {
		p := getPost(t, l, y, int(m), d)
		if p != nil {
			forDate(t, p, y, int(m), d)
			if !hasData(t, p) {
				t.Errorf("Language %v is not unavailable in less than 7 days", l)
			}
		}
	}
}
