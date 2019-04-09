package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

func main() {
	type Words struct {
		Word1 string `gorm:"index;type:varchar(512)"`
		Word2 string `gorm:"type:varchar(512)"`
	}

	db, err := gorm.Open("sqlite3", "words.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var word string
	var tweet string

	for {
		word = ""
		tweet = ""
		for {
			var words Words
			db.Raw("SELECT * FROM words WHERE word1 = ? ORDER BY RANDOM() LIMIT 1", word).Scan(&words)
			word = words.Word2
			tweet = tweet + word
			if word == "" {
				break
			}
		}

		if len(tweet) < 140 {
			break
		}
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	twitterApi := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))

	if _, err := twitterApi.PostTweet(tweet, nil); err != nil {
		panic(err)
	}

}
