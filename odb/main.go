// Commandline tool to fetch and print odb post
package main

// TODO: Add output format options

import (
	"flag"
	"fmt"
	"github.com/saintfish/odb"
	"os"
	"strconv"
)

var language = flag.String("language", "en", "Language of odb website. Currently supported languages are en, zh-hans and zh-hant")

var languageCodeMap = map[string]odb.Language{
	"en":      odb.English,
	"zh-hans": odb.SimplifiedChinese,
	"zh-hant": odb.TraditionalChinese,
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage odb [year] [month] [day]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 3 {
		usage()
	}
	year, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse year: %s", args[0])
		os.Exit(2)
	}
	month, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse month: %s", args[1])
		os.Exit(2)
	}
	day, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse day: %s", args[2])
		os.Exit(2)
	}
	l, ok := languageCodeMap[*language]
	if !ok {
		fmt.Fprintf(os.Stderr, "Bad language code: %s", *language)
		os.Exit(2)
	}
	o := odb.NewOdb(l)
	p, err := o.GetPost(year, month, day)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in fetch the page: %s", err)
		os.Exit(1)
	}
	fmt.Printf("%+v", p)
}
