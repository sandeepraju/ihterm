package pkg

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	ihBaseURL = "https://www.indiehackers.com/"

	ihPostSelector              = "div.homepage__thread-list"
	ihPostAuthorNameSelector    = "div.thread__metadata span.user-link__username"
	ihPostAuthorProfileSelector = "div.thread__metadata div a"
	ihPostPublishDateSelector   = "div.thread__metadata a.thread__date"
	ihPostCommentSelector       = "div.thread__metadata a.thread__reply-count"
	ihPostTitleSelector         = "div.thread__details a.thread__title"
	ihPostUpvotesSelector       = "div.thread-voter div.thread-voter__text div.thread-voter__count"

	topNPosts = 15
)

// IHTerm ...
type IHTerm struct {
	BaseURL *url.URL
	Title   string
	Posts   []IHPost
}

func (iht *IHTerm) downloadPosts() {

	// Fetch the document
	res, err := http.Get(iht.BaseURL.String())
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalln(fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	iht.Posts = make([]IHPost, 0)
	doc.Find(ihPostSelector).Children().Slice(0, topNPosts).Each(func(i int, s *goquery.Selection) {
		ihPost := IHPost{}
		ihPost.Author = &IHPostAuthor{}

		// Save author's name
		ihPost.Author.Name = strings.TrimSpace(s.Find(ihPostAuthorNameSelector).Text())

		// Save author's profile
		if authorProfile, ok := s.Find(ihPostAuthorProfileSelector).Attr("href"); ok {
			part, err := url.Parse(authorProfile)
			if err != nil {
				log.Fatalln(err.Error())
			}
			ihPost.Author.Profile = iht.BaseURL.ResolveReference(part)
		}

		// Save post's time
		ihPost.ApproxTime = strings.TrimSpace(s.Find(ihPostPublishDateSelector).Text())

		// Save post's comment URL
		if postCommentURL, ok := s.Find(ihPostPublishDateSelector).Attr("href"); ok {
			part, err := url.Parse(postCommentURL)
			if err != nil {
				log.Fatalln(err.Error())
			}
			ihPost.CommentURL = iht.BaseURL.ResolveReference(part)
		}

		// Save post's comments
		commentCount := strings.TrimSpace(strings.Split(strings.TrimSpace(s.Find(ihPostCommentSelector).Text()), " ")[0])
		comments, err := strconv.Atoi(commentCount)
		if err != nil {
			log.Fatalln(err)
		}
		ihPost.Comments = comments

		// Save post's title
		ihPost.Title = strings.TrimSpace(s.Find(ihPostTitleSelector).Text())

		// Save post's url
		if postURL, ok := s.Find(ihPostTitleSelector).Attr("href"); ok {
			part, err := url.Parse(postURL)
			if err != nil {
				log.Fatalln(err)
			}

			if part.IsAbs() {
				ihPost.URL = part
			} else {
				ihPost.URL = iht.BaseURL.ResolveReference(part)
			}
		}

		// Save post's upvotes
		upvotes, err := strconv.Atoi(strings.TrimSpace(s.Find(ihPostUpvotesSelector).Text()))
		if err != nil {
			log.Fatalln(err.Error())
		}
		ihPost.Upvotes = upvotes

		// Add the constructed post to the list of posts
		iht.Posts = append(iht.Posts, ihPost)
	})
}

func (iht *IHTerm) BitBarString() string {
	output := fmt.Sprintf("%s\n---\n", iht.Title)
	for _, post := range iht.Posts {
		output += post.BitBarString()
	}

	return output
}

type IHPostAuthor struct {
	Name    string
	Profile *url.URL
}

type IHPost struct {
	Author     *IHPostAuthor
	Upvotes    int
	Title      string
	URL        *url.URL
	ApproxTime string
	Comments   int
	CommentURL *url.URL
}

// BitBarString ...
func (ihp *IHPost) BitBarString() string {
	return fmt.Sprintf("%s | href=%s\n", ihp.Title, ihp.URL.String()) +
		fmt.Sprintf("Upvotes: %d Comments: %d | href=%s\n---\n", ihp.Upvotes, ihp.Comments, ihp.CommentURL.String())
}

// NewIHTerm ...
func NewIHTerm() *IHTerm {
	ihURL, err := url.Parse(ihBaseURL)
	if err != nil {
		log.Fatal(err)
	}

	iht := &IHTerm{
		BaseURL: ihURL,
		Posts:   make([]IHPost, 0),
		Title:   "IH",
	}

	// pre-fetch the posts
	iht.downloadPosts()

	return iht
}
