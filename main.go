package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"
)

func clearScreen() {
	fmt.Print("\033c")
}

func main() {
	var useTimer bool
	for {
		fmt.Println("Voulez-vous jouer avec un timer ? (Oui/Non)")
		var timerChoice string
		fmt.Scanln(&timerChoice)
		if strings.ToLower(timerChoice) == "oui" {
			useTimer = true
		} else {
			useTimer = false
		}

		if useTimer {
			fmt.Println("Vous avez choisi de jouer avec un timer.")
		} else {
			fmt.Println("Vous avez choisi de jouer sans timer.")
		}

		playGame(useTimer)

		fmt.Println("Voulez-vous rejouer ? (Oui/Non)")
		var replay string
		fmt.Scanln(&replay)
		if strings.ToLower(replay) != "oui" {
			return
		}
	}
}

func playGame(useTimer bool) {
	clearScreen()
	red := color.New(color.FgRed)
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	purple := color.New(color.FgHiMagenta)
	yellow := color.New(color.FgHiYellow)
	wordsList, err := ioutil.ReadFile("words.txt")
	if err != nil {
		log.Fatal(err)
	}

	words := strings.Split(string(wordsList), "\n")
	word := words[rand.Intn(len(words))]

	maxAttempts := 10
	attemptsLeft := maxAttempts
	usedLetters := make(map[rune]bool)

	word = strings.ToLower(word)

	correct := make([]rune, len(word))
	for i := range correct {
		correct[i] = '_'
	}

	miss := ""
	timeLimit := 60 * time.Second

	hangmanStagesBytes, err := ioutil.ReadFile("hangman.txt")
	if err != nil {
		log.Fatal(err)
	}
	hangmanStages := strings.Split(string(hangmanStagesBytes), "\n\n")

	blue.Println("Vous entrez dans le jeu du pendu !")
	blue.Println("Quentin, notre cobaye, est prêt pour la pendaison !")
	blue.Println("Sauf si vous réussissez à trouver le mot avant !")
	blue.Printf("Vous avez %d tentatives pour trouver ce mot.\n", maxAttempts)
	time.Sleep(5 * time.Second)

	startTime := time.Now()

	for attemptsLeft > 0 {
		elapsedTime := time.Since(startTime)

		if useTimer {
			fmt.Printf("Temps écoulé : %s\n", elapsedTime)
			if elapsedTime > timeLimit {
				red.Println("Temps écoulé ! Vous avez perdu.")
				return
			}
		}

		clearScreen()
		fmt.Printf("Mot à deviner: %s\n", string(correct))
		purple.Printf("Lettres manquées : %s\n", miss)
		purple.Println(hangmanStages[maxAttempts-attemptsLeft])
		purple.Printf("Tentatives restantes : %d\n", attemptsLeft)
		purple.Println("Veuillez entrer une lettre :")

		reader := bufio.NewReader(os.Stdin)
		letter, _, err := reader.ReadRune()

		if err != nil {
			red.Println("Erreur de lecture.")
			time.Sleep(1 * time.Second)
			continue
		}

		letter = unicode.ToLower(letter)

		if !unicode.IsLetter(letter) || len([]rune(string(letter))) != 1 {
			red.Println("Veuillez entrer une seule lettre de l'alphabet.")
			time.Sleep(1 * time.Second)
			continue
		}

		if _, ok := usedLetters[letter]; ok {
			yellow.Printf("Vous avez déjà essayé la lettre '%c'.\n", letter)
			time.Sleep(1 * time.Second)
			continue
		}
		usedLetters[letter] = true

		correctGuess := false
		for i, char := range word {
			if char == letter {
				correct[i] = letter
				correctGuess = true
			}
		}
		if !correctGuess {
			miss += string(letter)
			attemptsLeft--
		}
		fmt.Printf("Mot à deviner: %s\n", string(correct))
		if string(correct) == word {
			green.Printf("Vous avez gagné ! Le mot était: %s\n", word)
			return
		}
	}

	red.Printf("Vous avez perdu ! Le mot était: %s\n", word)
}
