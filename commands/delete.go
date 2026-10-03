package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Delete(s *discordgo.Session, channel string, messageobj *discordgo.MessageCreate) {
	message := messageobj.Content
	authorID := messageobj.Author.ID
	//To divide message into separate words
	msg := strings.Split(message, " ")
	//only accept the word and the messages count
	if len(msg) != 2 {
		s.ChannelMessageSend(channel, Delete_Usage)
		return
	}
	//if the second element is a positive number, proceed, otherwise alert the user
	//(0 or less would make Discord fall back to its default of 50 messages)
	delete_count, err := strconv.Atoi(msg[1])
	if err != nil || delete_count < 1 {
		s.ChannelMessageSend(channel, Delete_Usage)
		return
	}

	//Second to last message (the last one is always the delete command)
	last_messages, err := s.ChannelMessages(channel, 2, "", "", "")
	if err == nil && len(last_messages) == 2 {
		switch last_messages[1].Content {
		//if someone is oxyel
		case Delete_Success, Delete_FuckYou:
			if !IsEnderlord(authorID) {
				s.ChannelMessageSend(channel, Delete_FuckYou)
				return
			}
		}
	}

	//After negotiation with Oleg, I came to number 79, no on vse ravno pidoras (po facts)
	if !IsEnderlord(authorID) && (delete_count > 79) {
		s.ChannelMessageSend(channel, "А не дохуя ли?")
		return
	}
	//Discord gives out at most 100 messages at a time, and one of them is the delete command
	delete_count = min(delete_count, 99)

	//message count + the delete command
	messages_to_delete, err := s.ChannelMessages(channel, delete_count+1, "", "", "")
	if err != nil {
		fmt.Println("error getting messages to delete", err)
		s.ChannelMessageSend(channel, Something_Broke)
		return
	}
	var ids []string
	for _, m := range messages_to_delete {
		ids = append(ids, m.ID)
	}
	err = s.ChannelMessagesBulkDelete(channel, ids)
	if err != nil {
		fmt.Println("error deleting messages", err)
		s.ChannelMessageSend(channel, Delete_Failed)
		return
	}
	author := messageobj.Author.Username
	fmt.Println(fmt.Sprint(author, " Has deleted ", strconv.Itoa(len(ids)-1), " messages:"))
	s.ChannelMessageSend(channel, Delete_Success)

	//Deleted messages logging
	//because printing takes a long time, put it after everything's deleted
	//I should look into threading or async processes for this
	count := 1
	for i := len(messages_to_delete) - 1; i >= 1; i-- {
		m := messages_to_delete[i]
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprint("\nauthor: ", m.Author.Username, "\n"))
		sb.WriteString(fmt.Sprint("message ", count, ": ", m.Content))
		if len(m.Attachments) > 0 {
			sb.WriteString(m.Attachments[0].Filename)
		}
		fmt.Println(sb.String())
		count++
	}
}
