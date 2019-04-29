package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
)

type questionItem struct {
	question string
	answer   string
}

type place struct {
	start int
	end   int
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

	sentences := strings.Split(string(content), "\n\n")

	rand.Seed(time.Now().UnixNano())

	line := 0
	help := 0
	for len(sentences) > 0 {

		// remove the \n from the last item
		sentences[len(sentences)-1] = strings.TrimSpace(sentences[len(sentences)-1])

		// select a random sentence
		i := rand.Int() % len(sentences)

		q := removeStars(sentences[i])

		line++
		fmt.Println(Bold(Blue(strconv.Itoa(line)+".")), q.question)
		help = 0

		for {
			fmt.Print(Blue("answer: "))

			var input string
			fmt.Scanln(&input)
			if q.answer == input {
				fmt.Println(Green("correct ✔\n"))
				break
			} else if input == ":h" || input == ":H" {
				help = help + 1
				if len(q.answer) == help {
					fmt.Println("The answer is: ", Green(q.answer))
					break
				} else {
					fmt.Println("Answer starts with: ", Red(q.answer[0:help]))
				}
			} else if input == ":s" || input == ":S" {
				fmt.Println("The answer is: ", Green(q.answer))
				break
			} else {
				fmt.Println(Red("wrong ✗\n"))
			}
		}

		// remove this sentence, so next time it won't show
		sentences = remove(i, sentences)
	}

	println("Finish :)")
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

	return questionItem{
		question: q,
		answer:   word,
	}
}

func remove(i int, items []string) []string {
	size := len(items)

	if size > 0 {
		items[i] = items[size-1] // copy last item
		return items[:size-1]    // remove last item
	}

	return items
}
