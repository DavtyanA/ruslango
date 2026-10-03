package commands

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Suggest a game to play
func WhatToPlay(authorID string) string {
	var responses []string
	name := Discord_names[authorID]
	if name == "" {
		name = "сука"
	}
	//special cases
	switch authorID {
	case Enderlord_ID:
		return "Хеллоу ебать ты кто"
	case Makich_ID:
		return "DOTA 2"
	case Ruslan_ID:
		responses = append(responses, "Ты чё, долбоёб, сам с собой разговариваешь?")
	case Japi_ID, Seagull_ID, Tanya_ID:
		responses = append(responses, `Пиздец ну давай я еще для девушек буду ответочки писать, в оригинале же 
													чтоб удобно было ответочка оканчивается на "заебал"`)
	}
	responses = append(responses, "DOTA 2", "Apex", "Овер", "Пубг", "Пираты", "Падающие Ребята", "Свояк", "Амогус", "Бля, "+name+", ты заебал. Ты сам решить не можешь?")
	response := GetRandomItem(responses)
	return response
}

// Get a random story
func StoryTelling() string {
	responses := []string{"Бля, мне когда 5 лет было, у меня была серьезная болезнь, из-за которой " +
		"я теперь жирный. В смысле ты не веришь? ой блять долбаеб я жирный не из-за болезни, " +
		"а из-за операции, если бы мне не сделали операцию, я бы умер. Дебил хочешь у мамы моей спроси",

		"Я бы встретился с тобой и обсудил это в жизни, так как за мной тут следят, " +
			"но похуй, я сейчас попытаюсь. Все, получилось. Так вот, короче за мной теперь " +
			"охотится парень моей бывшей Он в общем торгует " +
			"людьми и оружием Я хуй знает, че мне с ним делать, он мне угрожает и " +
			"пытается через левых людей инфу на меня найти. Че думаешь мне делать?",
		// \\uD83D\\uDE02\\uD83D\\uDE02\\uD83D\\uDE02 \\uD83D\\uDE06  do something later maybe idk.....
		"У меня брат уволился, сейчас работает на бирже Нью-Йорк Тайм Сквер, " +
			"знаешь такую? Ну это короче самая крутая биржа. Он взял меня к себе, я у него там " +
			"типо сись админ. бля, тупой т9, я хотел сказать сисадмин. Ну вот для меня это типо " +
			"стажировки + опыт работы",

		//need to change this one from svoyak
		"В следующий раз, когда я тебя увижу, Олег, я тебе разобью ебало, серьезно, " +
			"отвечай за слова свои. Вот докажи то, что ты ээ у тебя не получается. Давай, " +
			"alt+F10 ... Вот серьезно, ты такое ебло сидишь выебываешься, думаешь, " +
			"что тебе ничего не будет, скажи, яйца иметь надо, ты выглядишь, как петух блять, говоришь, как петух блять, " +
			"говоришь, как ебанный петух, почему нельзя быть как хоть немножечко иметь яйца," +
			"что-то говорить",

		"Я впервые в инете встретил педофила, причем такого явно, это как-то смешно " +
			"сидеть угарать над ним. ахахха. он типо написал моей девушке, я решил ради прикола " +
			"с ним поговорить, и он сейчас типо да ладно че такого, что у тебя есть парень, давай " +
			"встретимся, так далее. Я такой Ммм заебись. Бля может сказ может самому с этим челом" +
			" ради прикола встретиться"}
	return GetRandomItem(responses)
}

// Without a timeout, a stuck joke site would make the request wait forever
var httpClient = &http.Client{Timeout: 10 * time.Second}

// Get a random joke
func GetRandomAnecdote() string {
	api_url := "http://anecdotica.ru/api"
	skey := os.Getenv("ANEK_SECRET_KEY")
	method := "getRandItem" //need 100 rublikov

	q := url.Values{}
	q.Set("pid", os.Getenv("ANEK_PID"))
	q.Set("method", method)
	q.Set("uts", fmt.Sprint(time.Now().Unix()))
	q.Set("genre", "1")
	q.Set("lang", "1")
	q.Set("format", "txt")
	q.Set("charset", "utf-8")
	q.Set("markup", "1")
	q.Set("hash", GetMD5Hash(q.Encode()+skey))
	// these ones will be added when we reach 100 rubliks
	// q.Set("category", "json")
	// q.Set("note", "0")
	// q.Set("wlist", "0")
	// q.Set("censor", "0")

	final_url := api_url + "?" + q.Encode()
	resp, err := httpClient.Get(final_url)
	if err != nil {
		fmt.Println("error getting anek", err)
		return "Не удалось получить анек. Сервис поломался("
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Не удалось получить анек. Сервис поломался("
	}
	return string(body)
}

// Check the input before calling the actual roll function
func RollInput(msg string) string {
	message := strings.Split(msg, " ")
	if len(message) > 1 {
		num, err := strconv.Atoi(message[1])
		// negative numbers (and the very biggest one) used to crash the whole bot
		if err != nil || num < 0 || num == math.MaxInt {
			return "ептвою мать пиши нормально `ролл 5`, `ролляй 100`, `roll 228` нахуй мне твои буквы"
		} else {
			return Roll(num)
		}
	} else {
		return Roll(100)
	}
}

// roll a random number between 0 and input (including input)
// returns a string because the other function wants a string returned.
// It's a question of making the code more compact and readable or more correct
func Roll(input int) string {
	return strconv.Itoa(rand.Intn(input + 1))
}
