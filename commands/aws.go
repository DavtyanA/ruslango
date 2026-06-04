package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/bwmarrin/discordgo"
)

var bucket string
var s3Client *s3.Client
var region = "us-east-2"

func init() {
	bucket = "ruslanbot"

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Println("error loading AWS config", err)
		return
	}

	s3Client = s3.NewFromConfig(cfg)

}

// Note that the item name is converted to lowercase in here
func downloadFromS3Bucket(item string) (string, error) {
	downloader := manager.NewDownloader(s3Client)
	ss := strings.Split(item, "/")
	fileName := ss[len(ss)-1]
	file, err := os.Create("tmp/" + fileName)
	if err != nil {
		fmt.Println("error creating a file", fileName, err)
		return "", err
	}
	defer file.Close()
	numBytes, err := downloader.Download(context.TODO(), file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(strings.ToLower(item)),
		})
	if err != nil {
		fmt.Println("error downloading the file", item, err)
		return file.Name(), err
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return file.Name(), nil
}

func downloadFromS3BucketFolder(folder string) (string, error) {
	resp, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(strings.ToLower(folder)),
	})
	if err != nil {
		return "Could not get bucket, pls contact Oleg Ermolaev", err
	}

	items := resp.Contents
	item := *GetRandomItem(items).Key
	return downloadFromS3Bucket(item)
}

func SendRandomFileFromFolder(s *discordgo.Session, channel string, folder string) {

	//download file from s3 bucket and folder, returns the name of the file
	fileName, err := downloadFromS3BucketFolder(folder)
	if err != nil {
		fmt.Println("Unable to get items from bucket", err)
		s.ChannelMessageSend(channel, "ойой чета паламалась( Напиши Ендерлолу он там посмотрит че поломалось")
	} else {
		sendFile(s, channel, fileName)
	}
}

func SendFileFromS3(s *discordgo.Session, channel string, item string) {
	fileName, err := downloadFromS3Bucket(item)
	if err != nil {
		fmt.Println("Unable to get an item from bucket", err)
		s.ChannelMessageSend(channel, "ойой чета паламалась( Напиши Ендерлолу он там посмотрит че поломалось")
	} else {
		sendFile(s, channel, fileName)
	}
}

func sendFile(s *discordgo.Session, channel string, fileName string) {
	//open file to give it to discord
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error opening a file "+file.Name(), err)
		s.ChannelMessageSend(channel, "ойой чета паламалась( Напиши Ендерлолу он там посмотрит че поломалось")
	} else {
		s.ChannelFileSend(channel, filepath.Base(file.Name()), file)
		//to not leave any leftovers
		os.Remove(file.Name())
	}
}
