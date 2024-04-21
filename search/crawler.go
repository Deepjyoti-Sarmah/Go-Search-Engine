package search

import "time"


type CrawlData struct {
	Url string 
	Success bool
	ResponseCode int
	CrawlData ParseBody
}

type ParseBody struct {
	CrawlTime time.Duration
	PageTitle string
	PageDescription string
	Heading string
	Links Links
}

type Links struct {
	Internal []string
	External []string
}
