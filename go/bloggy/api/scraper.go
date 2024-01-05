package api

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/sebito91/bootdotdev/go/bloggy/internal/database"
)

// StartScraping is a goroutine that handles the article scraping from various sources
// within the 'feeds' database table. This function will itself kick off up to `concurrency`
// goroutines to fetch deduplicated sources from their locations.
func (ac *apiConfig) StartScraping(concurrency int, sleepInterval time.Duration) {
	log.Printf("starting to scrape records using %d goroutines set to poll every %s\n", concurrency, sleepInterval)
	ticker := time.NewTicker(sleepInterval)

	feedsArgs := database.GetNextFeedsToFetchParams{
		Limit: int32(concurrency),
	}

	// kick off the ticker to start our collection
	for ; ; <-ticker.C {
		feedsArgs.LastFetchedAt = time.Now().Add(-sleepInterval)
		feeds, err := ac.DB.GetNextFeedsToFetch(context.Background(), feedsArgs)
		if err != nil {
			log.Printf("could not fetch feeds: %s\n", err)
			continue
		}

		log.Printf("found %d feeds to fetch\n", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go ac.fetchFeed(wg, feed)
		}
		wg.Wait()
	}
}

// fetchFeed is a helper function to retrieve the XML from the given RSS feed. Once
// an attempt to fetch is made each feed is updated in the `feeds` table.
func (ac *apiConfig) fetchFeed(wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := ac.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: time.Now(),
		ID:            feed.ID,
	})
	if err != nil {
		log.Printf("could not mark feed %s: %s\n", feed.ID, err)
		return
	}

	feedItems, err := ac.scrapeFeed(feed.Url)
	if err != nil {
		log.Printf("could not scrape feed %s from url %s: %s\n", feed.ID, feed.Url, err)
		return
	}

	log.Printf("fetched feed %s from url %s with %d items\n", feedItems.Channel.Title, feedItems.Channel.Link, len(feedItems.Channel.Item))

	for _, feedItem := range feedItems.Channel.Item {
		log.Printf("feed %s - item %s with date %s and description %s", feedItem.Title, feedItem.Link, feedItem.PubDate, feedItem.Description)
	}
}

// scrapeFeed will do the actual http call out to the provided URL and parse out the RSS information
func (ac *apiConfig) scrapeFeed(feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(feedURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}

	return &rssFeed, nil
}

// RSSFeed is made up of the channel descriptions and individual RSSItem items
type RSSFeed struct {
	Channel struct {
		Title         string    `xml:"title"`
		Link          string    `xml:"link"`
		Description   string    `xml:"description"`
		Generator     string    `xml:"generator,omitempty"`
		Language      string    `xml:"language"`
		LastBuildDate string    `xml:"lastBuildDate,omitempty"`
		Item          []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem is an individual item within an RSS feed
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Description string `xml:"description,omitempty"`
}
