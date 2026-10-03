package commands

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

// Get a random item from an array, works with dynamic types since Go 1.18! Zaebumba!
// (no need for len - 1: rand.Intn(n) already returns a number from 0 to n-1)
func GetRandomItem[T any](inputarray []T) T {
	return inputarray[rand.Intn(len(inputarray))]
}

// String contains to lower case
func StringContains(S string, sub string) bool {
	return strings.Contains(strings.ToLower(S), strings.ToLower(sub))
}

// String starts with to lower case
func StringStartsWith(S string, sub string) bool {
	return strings.HasPrefix(strings.ToLower(S), strings.ToLower(sub))
}

// String contains but for array of substrings (for convenience)
func StringStartsWithArray(S string, subs []string) bool {
	for _, s := range subs {
		if StringStartsWith(S, s) {
			return true
		}
	}
	return false
}

// String contains but for array of substrings (for convenience)
func StringContainsArray(S string, subs []string) bool {
	for _, s := range subs {
		if StringContains(S, s) {
			return true
		}
	}
	return false
}

// Like StringContainsArray, but only matches whole words, so "бан" doesn't fire on "банк" or "кабан".
// Go's regexp \b only knows English letters, so this compares the words themselves instead
func StringContainsWordArray(S string, words []string) bool {
	text := " " + onlyWords(S) + " "
	for _, w := range words {
		if strings.Contains(text, " "+onlyWords(w)+" ") {
			return true
		}
	}
	return false
}

// Lowercase the text and keep only its words, separated by single spaces
func onlyWords(S string) string {
	words := strings.FieldsFunc(strings.ToLower(S), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	return strings.Join(words, " ")
}

// Mega story created by Rustam ili Vova ya xz
func MegaStory(s *discordgo.Session, channel string) {

	SendFileFromS3(s, channel, Pictures_Folder_Other+"daiproidu.jpg")
	time.Sleep(1500 * time.Millisecond)
	SendFileFromS3(s, channel, Pictures_Folder_Other+"xuya.jpg")
	time.Sleep(1500 * time.Millisecond)
	SendFileFromS3(s, channel, Pictures_Folder_Other+"poebalu.jpg")
	time.Sleep(1500 * time.Millisecond)
	SendFileFromS3(s, channel, Pictures_Folder_Other+"choblyatb.jpg")
	time.Sleep(1500 * time.Millisecond)
	SendFileFromS3(s, channel, Pictures_Folder_Other+"razebu.jpg")
	time.Sleep(1500 * time.Millisecond)
	SendFileFromS3(s, channel, Pictures_Folder_Other+"willsee.jpg")
	time.Sleep(1500 * time.Millisecond)
	SendFileFromS3(s, channel, Pictures_Folder_Other+"taashaa.jpg")
	time.Sleep(1500 * time.Millisecond)
	SendFileFromS3(s, channel, Pictures_Folder_Other+"bilyateblo.jpg")
	time.Sleep(1500 * time.Millisecond)
	SendFileFromS3(s, channel, Pictures_Folder_Other+"che tam.jpg")
	time.Sleep(1500 * time.Millisecond)
	SendFileFromS3(s, channel, Pictures_Folder_Other+"blyatb.jpg")
}

// Get MD5 hash
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// For easier checking
func IsEnderlord(id string) bool {
	return id == Enderlord_ID
}
