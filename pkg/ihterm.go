package pkg

import (
	"net/url"
	"log"
	"net/http"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

const (
	ihUrl = "https://www.indiehackers.com/"
	ihPostSelector = "div.thread-row div.ember-view"
	ihPostUpvotesSelector = "div.thread-row.ember-view div.thread-voter.ember-view div.thread-voter__text div.thread-voter__count"
	ihPostTitleSelector = "div.thread-row.ember-view div.thread__details a.thread__title.ember-view"

)

// IHTerm ...
type IHTerm struct {
	Url *url.URL
}

func (iht *IHTerm) Posts() ([]IHPost, error) {
	res, err := http.Get(iht.Url.String())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}


	// posts := make([]IHPost, 0)
	fmt.Println("finding")
	fmt.Println(doc.Find("div.homepage__thread-list").Children().Length())
	doc.Find("div.homepage__thread-list").Children().Each(func(i int, s *goquery.Selection) {
		fmt.Println(i)
		userName := strings.TrimSpace(s.Find("span.user-link__username").Text())
		fmt.Println(userName)
		// posts = append(posts, )
		// For each item found, get the
		// upvotes := s.Find(ihPostUpvotesSelector).Text()
		// title := s.Find(ihPostTitleSelector).Text()
		// fmt.Println(title)
		// fmt.Println(upvotes)
	})

	return nil, nil
}

type IHPost struct {
	Upvotes int
	Title string
	Url *url.URL
	Author string
	ApproxTime string
	Comments int
}

// NewIHTerm ...
func NewIHTerm() *IHTerm {
	ihUrl, err := url.Parse(ihUrl)
	if err != nil {
		log.Fatal(err)
	}
	return &IHTerm{
		Url: ihUrl,
	}
}
