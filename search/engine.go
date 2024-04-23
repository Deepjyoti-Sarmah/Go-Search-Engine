package search

import (
	"fmt"
	"time"

	"github.com/Deepjyoti-Sarmah/GolangSearchEngine/db"
)

func RunEngine() {
	fmt.Println("started search engine crawl...")
	defer fmt.Println("search engine crawl has finished")
	//Get settings
	settings := &db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		fmt.Println("setting went wrong getting the settings")
		return
	}
	// check the search engine is on
	if !settings.SearchOn {
		fmt.Println("search is turned off")
		return
	}

	crawl := &db.CrawledUrl{}
	nextUrls, err := crawl.GetNextCrawleUrls(int(settings.Amount))
	if err != nil {
		fmt.Println("something went wrong getting next urls")
		return
	}

	newUrls := []db.CrawledUrl{}
	testedTime := time.Now()
	for _, next := range nextUrls {
		results := runCrawl(next.Url)
		if !results.Success {
			err := next.UpdatedUrl(db.CrawledUrl{
				ID:              next.ID,
				Url:             next.Url,
				Success:         false,
				CrawlDuration:   results.CrawlData.CrawlTime,
				ResponseCode:    results.ResponseCode,
				PageTitle:       results.CrawlData.PageTitle,
				PageDescription: results.CrawlData.PageDescription,
				Heading:         results.CrawlData.Heading,
				LastTested:      &testedTime,
			})
			if err != nil {
				fmt.Println("something went wrong updating a failed url")
			}
			continue
		}
		//Success
		err := next.UpdatedUrl(db.CrawledUrl{
			ID:              next.ID,
			Url:             next.Url,
			Success:         results.Success,
			CrawlDuration:   results.CrawlData.CrawlTime,
			ResponseCode:    results.ResponseCode,
			PageTitle:       results.CrawlData.PageTitle,
			PageDescription: results.CrawlData.PageDescription,
			Heading:         results.CrawlData.Heading,
			LastTested:      &testedTime,
		})
		if err != nil {
			fmt.Println("something went wrong updating a Success url")
			fmt.Println(next.Url)
		}

		for _, newUrl := range results.CrawlData.Links.External {
			newUrls = append(newUrls, db.CrawledUrl{Url: newUrl})
		}
	} // end of range
	if !settings.AddNew {
		return
	}
	//Insert new urls
	for _, newUrl := range nextUrls {
		err := newUrl.Save()
		if err != nil {
			fmt.Println("Something went wrong add the new url to the database")
		}
	}
	fmt.Printf("\n Added %d new urls to the database", len(newUrls))
}

func RunIndex() {
	fmt.Println("started search indexing...")
	defer fmt.Println("search indexing has finished")
	crawled := &db.CrawledUrl{}
	notIndexed, err := crawled.GetNotIndex()
	if err != nil {
		return
	}

	idx := make(Index)
	idx.Add(notIndexed)
	searchIndex := &db.SearchIndex{}
	err = searchIndex.Save(idx, notIndexed)
	if err != nil {
		fmt.Println(err)
		fmt.Println("something went wrong saving the index")
		return
	}
	err = crawled.SetIndexedTrue(notIndexed)
	if err != nil {
		fmt.Println("somethign went wrong updating the indexed urls")
		return
	}
}
