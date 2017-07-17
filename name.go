package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	doubleVowel Suffix = "doubleVowel"
	singleVowel Suffix = "singleVowel"
	singleConsonant Suffix = "singleConsonant"
	none Suffix = ""
)

var (
	probabilities = map[NameState]TProb{
		NameState{0, none}: TProb{.4, .6, 0},
		NameState{0, singleConsonant}: TProb{1, 0, 0},
		NameState{1, doubleVowel}: TProb{0, .7, .3},
		NameState{1, singleVowel}: TProb{.1, .8, .1},
		NameState{1, singleConsonant}: TProb{.9, 0, .1},
		NameState{2, doubleVowel}: TProb{0, .1, .9},
		NameState{2, singleVowel}: TProb{.1, .3, .6},
		NameState{2, singleConsonant}: TProb{.6, 0, .4},
		NameState{3, doubleVowel}: TProb{0, 0, 1},
		NameState{3, singleVowel}: TProb{.05, .1, .85},
		NameState{3, singleConsonant}: TProb{0, 0, 1},
		NameState{4, singleConsonant}: TProb{0, 0, 1},
	}
	vowelPhonemes = []Phoneme{}
	consonantPhonemes = []Phoneme{}
)

type Suffix string

type Phoneme struct {
	String string
	Type string
	Available bool
	Spelled bool
}

type NameState struct {
	Syllables int
	Suffix Suffix
}

// TProb represents the transition probablities of a suffix.
type TProb struct {
	Vowel, Consonant, End float32
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Missing path to phonemes file")
	}
	phonemesPath := os.Args[1]

	if len(os.Args) < 3 {
		log.Fatalf("Missing name count")
	}
	count, err := strconv.Atoi(os.Args[2])
	if err != nil || count < 1 {
		log.Fatalf("Name count is not valid")
	}

	f, err := os.Open(phonemesPath)
	defer f.Close()
	if err != nil {
		log.Fatalf("Failed opening phonemes file: %v\n", err)
	}

	phonemeStrs, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	if len(phonemeStrs) < 2 {
		log.Fatalf("len(phonemes) == 0")
	}
	// Convert to Phoneme data strucutures.
	phonemes := make([]Phoneme, len(phonemeStrs)-1)
	for i, row := range phonemeStrs[1:] {
		phonemes[i].String = row[0]
		phonemes[i].Type = row[1]
		phonemes[i].Available = row[2] == "true"
		phonemes[i].Spelled = row[3] == "true"
	}

	commonPhonemes := []Phoneme{}
	for _, p := range phonemes {
		//if p.Available && p.Spelled {
		if p.Available {
			commonPhonemes = append(commonPhonemes, p)
		}
	}

	fmt.Printf("Common phonemes: ")
	for _, p := range commonPhonemes {
		fmt.Print(p.String + " ")
		if p.Type == "vowel" {
			vowelPhonemes = append(vowelPhonemes, p)
		} else {
			consonantPhonemes = append(consonantPhonemes, p)
		}
	}
	fmt.Printf("\n")

	if len(vowelPhonemes) == 0 || len(consonantPhonemes) == 0 {
		log.Fatal("Not enough phonemes")
	}

	rand.Seed(time.Now().Unix())
	for i := 0; i < count; i++ {
		fmt.Printf("Name: ")
		name := []Phoneme{}
		p, more := pickPhoneme(probabilities[getState(name)])
		for more {
			fmt.Printf("%s", strings.Replace(p.String, "/", " ", -1))
			name = append(name, p)
			p, more = pickPhoneme(probabilities[getState(name)])
		}
		fmt.Printf("\n")
	}
}

func pickPhoneme(prob TProb) (Phoneme, bool) {
	i := rand.Float32()
	if i >= prob.Vowel+prob.Consonant {
		return Phoneme{}, false
	} else if i < prob.Vowel {
		return vowelPhonemes[rand.Intn(len(vowelPhonemes))], true
	} else {
		return consonantPhonemes[rand.Intn(len(consonantPhonemes))], true
	}
}

// getState gets the NameState of a name. `name` must have at least one Phoneme.
func getState(name []Phoneme) NameState {
	state := NameState{0, none}
	if len(name) == 0 {
		return state
	}
	for i, p := range name {
		if p.Type == "vowel" {
			if i == 0 || name[i-1].Type != "vowel" {
				state.Syllables++
			}
		}
	}
	if name[len(name)-1].Type == "vowel" {
		state.Suffix = singleVowel
		if len(name) > 1 && name[len(name)-2].Type == "vowel" {
			state.Suffix = doubleVowel
		}
	} else {
		state.Suffix = singleConsonant
	}
	return state
}
