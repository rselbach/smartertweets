// Just an experiment... Public domain, enjoy it.
package main

import (
	"http"
	"regexp"
	"twitterstream" // https://github.com/hoisie/{httplib.go,twitterstream}
	"url"
)

var (
	urlRegExp = regexp.MustCompile("https?://[a-zA-Z0-9][a-zA-Z0-9\\.]*[a-zA-Z](:[0-9]+)?(/[^ ?&]*)*(\\?[^ ]+&?)*(#[a-zA-Z0-9]*)?")
	rtRegExp  = regexp.MustCompile("RT @")
)

const (
	accessToken = "FACEBOOK ACCESS TOKEN"
	appId = "FACEBOOK APP ID"
	appSecret = "FACEBOOK APP SECRET"
	postURL = "https://graph.facebook.com/PROFILE_ID/feed"
	linkURL = "https://graph.facebook.com/PROFILE_ID/links"
)

func PostFacebookLink(text string, urls []string) {
	u := make(url.Values)
	u.Add("access_token", accessToken)
	u.Add("link", url.QueryEscape(urls[0]))
	u.Add("message", url.QueryEscape(text))
	http.PostForm(linkURL, u)
}

func PostFacebookStatus(text string) {
	u := make(url.Values)
	u.Set("access_token", accessToken)
	u.Set("message", url.QueryEscape(text))
	http.PostForm(postURL, u)
}

func main() {
	stream := make(chan *twitterstream.Tweet)
	client := twitterstream.NewClient("TWITTER USER", "TWITTER PASSWORD")
	err := client.Follow([]int64{8865192 /*USER IDs TO FOLLOW*/}, stream)
	if err != nil {
		println(err.String())
	}
	for {
		tw := <-stream
		// if it's not a reply
		if tw.In_reply_to_user_id == 0 {
			// and not a retweet
			if rtRegExp.FindString(tw.Text) == "" {

				// this is a valid tweet to export
				// check if it's got links
				us := urlRegExp.FindAllString(tw.Text, -1)
				if len(us) > 0 {
					PostFacebookLink(tw.Text, us)
				} else {
					PostFacebookStatus(tw.Text)
				}
			}
		}
	}
}
