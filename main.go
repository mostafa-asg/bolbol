package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	"github.com/mostafa-asg/bolbol/db"
)

type questionItem struct {
	question string
	answer   string
}

type place struct {
	start int
	end   int
}

func printQuestion(current int, total int, question string) {
	fmt.Println(Bold(Blue(strconv.Itoa(current)+"/"+strconv.Itoa(total)+".")), question)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Please specify a file")
		return
	}

	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		println("Error", err.Error())
		return
	}
	fileInfo, _ := os.Stat(filename)

	sentences := strings.Split(string(content), "\n\n")

	rand.Seed(time.Now().UnixNano())

	startTime := time.Now()
	question := 0
	help := 0
	mistakes := 0
	notKnow := 0
	nowKnowWords := ""
	totalQuestions := len(sentences)
	inputReader := bufio.NewReader(os.Stdin)
	answered := false

	dbConn, err := db.Open()
	if err != nil {
		fmt.Println("Cannot open database", err.Error())
	} else {
		defer db.Close(dbConn)
	}

	for len(sentences) > 0 {

		// remove the \n from the last item
		sentences[len(sentences)-1] = strings.TrimSpace(sentences[len(sentences)-1])

		// select a random sentence
		i := rand.Int() % len(sentences)

		q := removeStars(sentences[i])

		question++
		printQuestion(question, totalQuestions, q.question)
		help = 0

		for {
			fmt.Print(Blue("answer: "))

			var input string
			input, _ = inputReader.ReadString('\n')
			if len(input) > 0 {
				input = input[0 : len(input)-1] // remove the `\n`
			}

			if strings.Compare(q.answer, input) == 0 {
				fmt.Println(Green("correct ✔\n"))
				answered = true
				break
			} else if input == ":h" || input == ":H" {
				help = help + 1
				if len(q.answer) == help {
					notKnow++
					nowKnowWords += q.answer + "\n"
					fmt.Println("The answer is: ", Green(q.answer))
					break
				} else {
					fmt.Println("Answer starts with: ", Red(q.answer[0:help]))
				}
			} else if input == ":s" || input == ":S" {
				notKnow++
				nowKnowWords += q.answer + "\n"
				fmt.Println("The answer is: ", Green(q.answer))
				answered = false
				break
			} else {
				mistakes++
				fmt.Println(Red("wrong ✗\n"))
			}
		}

		if dbConn != nil {
			go func(word string, answered bool) {
				db.InsertWordReport(dbConn, word, fileInfo.Name(), answered)
			}(q.answer, answered)
		}

		// remove this sentence, so next time it won't show
		sentences = remove(i, sentences)
	}

	fmt.Printf("Total questions: %d\n", Green(totalQuestions))

	if mistakes > 0 {
		fmt.Printf("mistakes: %d\n", Red(mistakes))
	} else {
		fmt.Printf("mistakes: %s\n", Green("nothing"))
	}

	if notKnow > 0 {
		fmt.Printf("Do not know: %d\n", Red(notKnow))
	} else {
		fmt.Printf("Do not know: %s\n", Green("nothing"))
	}

	if nowKnowWords != "" {
		fmt.Printf("Word(s) you didn't know:\n%s", Red(nowKnowWords))
	}

	fmt.Printf("Time spent: %v\n", Blue(time.Since(startTime)))
	println("Finish :)")

	seconds := time.Now().Unix() - startTime.Unix()

	if dbConn != nil {
		err = db.InsertFileReport(dbConn, fileInfo.Name(), mistakes, notKnow, seconds)
		if err != nil {
			log.Println(err)
		}
	}
}

func contains(arr []string, key string) bool {
	for _, item := range arr {
		if item == key {
			return true
		}
	}
	return false
}

// Find all words that mark with star
// One of them choose randomly and convert to ________
// For the rest, remove only star
func removeStars(sentence string) questionItem {
	places := make([]place, 0)

	foundFirstStar := false
	firstStarIndex := 0

	// find the position of *
	for i := 0; i < len(sentence); i++ {
		if sentence[i] == '*' {
			if !foundFirstStar {
				foundFirstStar = true
				firstStarIndex = i
			} else {
				foundFirstStar = false
				places = append(places, place{firstStarIndex, i})
			}
		}
	}

	if len(places) == 0 {
		return questionItem{
			question: sentence,
			answer:   "",
		}
	}

	hidePlaceIndex := rand.Int() % len(places)
	var word = sentence[places[hidePlaceIndex].start+1 : places[hidePlaceIndex].end]

	q := strings.Replace(sentence, "*"+word+"*", "__________", -1)
	q = strings.Replace(q, "*", "", -1)
	q = strings.Replace(q, "|:s|", "("+scramble(word)+")", -1)

	return questionItem{
		question: q,
		answer:   word,
	}
}

func scramble(word string) string {
	var s = len(word)
	var newWord = make([]byte, s)

	// First copy all chars in their right position
	for i := 0; i < s; i++ {
		newWord[i] = word[i]
	}

	// Swap 2 chars randomly
	for i := 1; i <= 25; i++ {
		var a = rand.Int() % s
		var b = rand.Int() % s

		var tmp = newWord[a]
		newWord[a] = newWord[b]
		newWord[b] = tmp
	}

	return string(newWord)
}

func remove(i int, items []string) []string {
	size := len(items)

	if size > 0 {
		items[i] = items[size-1] // copy last item
		return items[:size-1]    // remove last item
	}

	return items
}
