package commands

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

// Send a random file from a folder in the bucket
func SendRandomFileFromFolder(s *discordgo.Session, channel string, folder string) {
	resp, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(strings.ToLower(folder)),
	})
	if err != nil {
		fmt.Println("Unable to get items from bucket", err)
		s.ChannelMessageSend(channel, Something_Broke)
		return
	}

	// Skip the empty placeholder objects the S3 console makes for folders (their names end with "/")
	var items []string
	for _, object := range resp.Contents {
		key := aws.ToString(object.Key)
		if !strings.HasSuffix(key, "/") {
			items = append(items, key)
		}
	}
	if len(items) == 0 {
		fmt.Println("No files in bucket folder", folder)
		s.ChannelMessageSend(channel, Something_Broke)
		return
	}

	SendFileFromS3(s, channel, GetRandomItem(items))
}

// Send a file from the bucket straight to Discord, without saving it to disk first.
// Note that the item name is converted to lowercase in here.
// Returns the sent message, or nil if something went wrong
func SendFileFromS3(s *discordgo.Session, channel string, item string) *discordgo.Message {
	object, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(strings.ToLower(item)),
	})
	if err != nil {
		fmt.Println("Unable to get an item from bucket", item, err)
		s.ChannelMessageSend(channel, Something_Broke)
		return nil
	}
	defer object.Body.Close()

	message, err := s.ChannelFileSend(channel, path.Base(item), object.Body)
	if err != nil {
		fmt.Println("error sending a file", item, err)
		return nil
	}
	return message
}
