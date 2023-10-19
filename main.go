package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

func main() {
	for {
		playGame()

		fmt.Println("Voulez-vous rejouer ? (Oui/Non)")
		var replay string
		fmt.Scanf("%s", &replay)
		if strings.ToLower(replay) != "oui" {
			break
		}
	}
}

func playGame() {
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

	fmt.Println("Vous entrez dans le jeu du pendu !")
	fmt.Println("Quentin, notre cobaye, est prêt pour la pendaison !")
	fmt.Println("Sauf si vous réussissez à trouver le mot avant !")
	fmt.Printf("Vous avez %d tentatives pour trouver ce mot.\n", maxAttempts)

	startTime := time.Now()

	for attemptsLeft > 0 {
		elapsedTime := time.Since(startTime)

		fmt.Printf("Temps écoulé : %s\n", elapsedTime)

		if elapsedTime > timeLimit {
			fmt.Println("Temps écoulé ! Vous avez perdu.")
			return
		}

		fmt.Printf("Lettres manquées : %s\n", miss)
		fmt.Println(hangmanStages[maxAttempts-attemptsLeft])
		fmt.Printf("Tentatives restantes : %d\n", attemptsLeft)
		fmt.Println("Veuillez entrer une lettre :")

		var letter rune
		fmt.Scanf("%c\n", &letter)

		if letter < 'a' || letter > 'z' {
			fmt.Println("Veuillez entrer une lettre de l'alphabet.")
			continue
		}

		letter = unicode.ToLower(letter)

		if _, ok := usedLetters[letter]; ok {
			fmt.Printf("Vous avez déjà essayé la lettre '%c'.\n", letter)
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
			fmt.Printf("Vous avez gagné ! Le mot était: %s\n", word)
			return
		}
	}

	fmt.Printf("Vous avez perdu ! Le mot était: %s\n", word)
}
