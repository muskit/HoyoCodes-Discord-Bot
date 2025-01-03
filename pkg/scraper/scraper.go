package scraper

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// Given a Project Tactics article containing MiHoYo game codes,
// return a map of codes and their description, as well as
// the datetime which the data was updated.
func scrapePJT(url string, identifierText string) (map[string]string, string) {
	// scraped data
	activeCodes := make(map[string]string)
	datetime := ""

	c := colly.NewCollector(colly.AllowedDomains("www.pockettactics.com"))

	// --- callback setup ---
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %s...\n", url)
	})

	// populate codes
	c.OnHTML("strong", func(h *colly.HTMLElement) {
		if strings.Contains(h.Text, identifierText) {
			log.Printf("FOUND HEADER: %s\n", h.Text)
			if strings.Contains(h.Text, "expire") {
				log.Println("Appears to have expired according to header; stopping...")
				return
			}

			log.Println("Gathering codes...")

			listContainer := h.DOM.Parent().Next()
			list := listContainer.Children()

			for i, elem := range list.Nodes {
				entry := elem.FirstChild
				key := entry.FirstChild.Data
				desc := string([]rune(entry.NextSibling.Data)[3:])

				activeCodes[key] = desc
				log.Printf("%d: [%s] (%s)\n", i, key, desc)
			}
		} else {
			// log.Printf("Didn't find \"%s\"\n", identifierText)
		}
	})

	// populate datetime
	c.OnHTML("time", func(h *colly.HTMLElement) {
		if h.DOM.HasClass("updated") {
			datetime = h.Attr("datetime")
			log.Printf("Update datetime: %s\n", datetime)
		}
	})

	// begin scrape
	c.Visit(url)

	// TODO: check that data to return is good

	log.Println("done")
	return activeCodes, datetime
}

func ScrapeHI3() (map[string]string, string) {
	log.Println("--- [HONKAI IMPACT] ---")
	return scrapePJT(
		"https://www.pockettactics.com/honkai-impact/codes",
		"Here are all the new Honkai Impact codes",
	)
}

func ScrapeGI() (map[string]string, string) {
	log.Println("--- [GENSHIN IMPACT] ---")
	return scrapePJT(
		"https://www.pockettactics.com/genshin-impact/codes",
		"Here are all of the new Genshin Impact codes",
	)
}

func ScrapeHSR() (map[string]string, string) {
	log.Println("--- [HONKAI STAR RAIL] ---")
	return scrapePJT(
		"https://www.pockettactics.com/honkai-star-rail/codes",
		"Here are all of the new Honkai Star Rail codes",
	)
}

func ScrapeHSRLive() (map[string]string, string) {
	log.Println("--- [HONKAI STAR RAIL LIVESTREAM CODES] ---")
	return scrapePJT(
		"https://www.pockettactics.com/honkai-star-rail/codes",
		"livestream codes",
	)
}

func ScrapeZZZ() (map[string]string, string) {
	log.Println("---[ZENLESS ZONE ZERO] ---")
	return scrapePJT(
		"https://www.pockettactics.com/zenless-zone-zero/codes",
		"Here are all of the new Zenless Zone Zero codes",
	)
}
