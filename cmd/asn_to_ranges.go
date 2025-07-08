package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/samber/lo"
	"log"
	"net/netip"
	"os"
	"strings"
)

func main() {
	var (
		as   string
		path string
	)

	flag.StringVar(&as, "as", "AS00000", "ASN to scrap")
	flag.StringVar(&path, "save", "as.txt", "File to save the output to")
	flag.Parse()

	ranges := scrapRanges(as)
	if len(ranges) == 0 {
		panic("ranges list is empty")
	}
	_ = os.Remove(path)
	_ = os.WriteFile(path, []byte(strings.Join(lo.Map(ranges, func(v netip.Prefix, _ int) string {
		return v.String()
	}), "\n")), 0644)
	log.Printf("Saved %d ranges to file %s", len(ranges), path)
}

// Credit: https://github.com/FDUTCH/bedrock_server_scanner/blob/b0945539e4b61e82302a427f372b726e93a310f7/scanner/ranges.go#L45
func scrapRanges(as string) []netip.Prefix {
	var prefixes []netip.Prefix
	c := colly.NewCollector()
	c.OnHTML("a[class]", func(e *colly.HTMLElement) {
		if e.Attr("class") == "charcoal-link " {
			prefix, err := netip.ParsePrefix(e.Text)
			if err == nil {
				addr := prefix.Addr()
				if addr.Is4() {
					prefixes = append(prefixes, prefix)
				}
			}
		}
	})
	_ = c.Visit(fmt.Sprintf("https://ipinfo.io/%s", as))
	return prefixes
}
