package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io"
	"os"
	"os/exec"
)

func main() {
	lambda.Start(Main)
}

func Main() {
	type Words struct {
		Word1 string `gorm:"index;type:varchar(512)"`
		Word2 string `gorm:"type:varchar(512)"`
	}

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
	})

	s3file, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String("/words.tar.xz"),
	})
	if err != nil {
		panic(err)
	}
	defer s3file.Body.Close()

	file, err := os.Create("/tmp/words.tar.xz")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := io.Copy(file, s3file.Body); err != nil {
		panic(err)
	}

	if err := exec.Command("sh", "-c", "cd /tmp && tar Jxvf words.tar.xz").Run(); err != nil {
		panic(err)
	}

	db, err := gorm.Open("sqlite3", "/tmp/words.db")
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
			if word == "" {
				break
			}
			word = words.Word2
			tweet = tweet + word
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

	if err := exec.Command("sh", "-c", "rm -rf /tmp/*").Run(); err != nil {
		panic(err)
	}

}
