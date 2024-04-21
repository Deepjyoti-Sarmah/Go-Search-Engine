package search

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type CrawlData struct {
	Url          string
	Success      bool
	ResponseCode int
	CrawlData    ParseBody
}

type ParseBody struct {
	CrawlTime       time.Duration
	PageTitle       string
	PageDescription string
	Heading         string
	Links           Links
}

type Links struct {
	Internal []string
	External []string
}

func runCrawl(inputUrl string) CrawlData {
	resp, err := http.Get(inputUrl)
	baseUrl, err := url.Parse(inputUrl)
	// check id error or if response is empty
	if err != nil || resp == nil {
		fmt.Println("something went wrong fetching the body")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: 0, CrawlData: ParseBody{}}
	}
	defer resp.Body.Close()
	//check for 200
	if resp.StatusCode != 200 {
		fmt.Println("non 200 code found")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParseBody{}}
	}
	//check for html
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		//response html
		data, err := parseBody(resp.Body, baseUrl)
		if err != nil {
			return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParseBody{}}
		}
		return CrawlData{Url: inputUrl, Success: true, ResponseCode: resp.StatusCode, CrawlData: data}
	} else {
		//response is not html
		fmt.Println("non html response found")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParseBody{}}
	}
}

func parseBody(body io.Reader, baseUrl *url.URL) (ParseBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return ParseBody{}, err
	}
	start := time.Now()
	//Get links
	links := getLinks(doc, baseUrl)
	//Get page Title and description
	title, desc := getPageData(doc)
	//Get H1 tags
	headings := getPageHeadings(doc)
	//Return the time & date
	end := time.Now()
	return ParseBody{
		CrawlTime:       end.Sub(start),
		PageTitle:       title,
		PageDescription: desc,
		Heading:         headings,
		Links:           links,
	}, nil
}

// Dept first search, recursive function for scanning the html tree
func getLinks(node *html.Node, baseUrl *url.URL) Links {
	links := Links{}
	if node == nil {
		return links
	}

	var findLinks func(*html.Node)
	findLinks = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					url, err := url.Parse(attr.Val)
					if err != nil || strings.HasPrefix(url.String(), "#") || strings.HasPrefix(url.String(), "mail") || strings.HasPrefix(url.String(), "tel") || strings.HasPrefix(url.String(), "javascript") || strings.HasPrefix(url.String(), ".md") {
						continue
					}
					if url.IsAbs() {
						if isSameHost(url.String(), baseUrl.String()) {
							links.Internal = append(links.Internal, url.String())
						} else {
							links.External = append(links.External, url.String())
						}
					} else {
						rel := baseUrl.ResolveReference(url)
						links.Internal = append(links.Internal, rel.String())
					}
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			findLinks(child)
		}
	}
	findLinks(node)
	return links
}

func isSameHost(absoluteUrl string, baseUrl string) bool {
	absUrl, err := url.Parse(absoluteUrl)
	if err != nil {
		return false
	}
	baseUrlParsed, err := url.Parse(baseUrl)
	if err != nil {
		return false
	}
	return absUrl.Host == baseUrlParsed.Host
}

func getPageData(node *html.Node) (string, string) {
	if node == nil {
		return "", ""
	}
	//find title and description
	title, desc := "", ""
	var findMetaAndTitle func(*html.Node)

	findMetaAndTitle = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" {
			//check if empty
			if node.FirstChild == nil {
				title = ""
			} else {
				title = node.FirstChild.Data
			}
		} else if node.Type == html.ElementNode && node.Data == "meta" {
			var name, content string
			for _, attr := range node.Attr {
				if attr.Key == "name" {
					name = attr.Val
				} else if attr.Key == "content" {
					content = attr.Val
				}
			}
			if name == "description" {
				desc = content
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findMetaAndTitle(child)
	}
	findMetaAndTitle(node)
	return title, desc
}

func getPageHeadings(node *html.Node) string {
	if node == nil {
		return ""
	}

	var heading strings.Builder
	var findH1 func(*html.Node)

	findH1 = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "h1" {
			//Check if node is empty
			if node.FirstChild != nil {
				heading.WriteString(node.FirstChild.Data)
				heading.WriteString(", ")
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			findH1(child)
		}
	}
	//remove the last comma
	return strings.TrimSuffix(heading.String(), ",")
}
