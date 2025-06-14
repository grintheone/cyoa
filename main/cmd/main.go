package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

type (
	Path  string
	Paths map[Path]Story
)

type Story struct {
	Title       string   `json:"title"`
	Description []string `json:"story"`
	Options     []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func traverse(paths Paths, startPath Path) {
	start := paths[startPath]

	title := color.New(color.FgCyan)
	description := color.New(color.FgHiWhite)
	option := color.New(color.FgYellow)

	fmt.Println()

	for _, ch := range start.Title {
		title.Fprintf(os.Stdout, "%v", string(ch))
		time.Sleep(time.Millisecond * 40)
	}

	fmt.Println()
	fmt.Println()

	for _, pg := range start.Description {
		for _, ch := range pg {
			description.Fprintf(os.Stdout, "%v", string(ch))
			time.Sleep(time.Millisecond * 40)
		}
		fmt.Println()
		fmt.Println()
	}

	for i, opt := range start.Options {
		option.Fprintf(os.Stdout, "Option %d: ", i+1)
		for _, ch := range opt.Text {
			option.Fprintf(os.Stdout, "%v", string(ch))
			time.Sleep(time.Millisecond * 40)
		}
		fmt.Println()
	}

	fmt.Println()

	if len(start.Options) > 0 {
		fmt.Println("Choose an option number from your keyboard (ESC anytime to quit):")

		var availableKeys []string
		var nextArc Path

		for i := range start.Options {
			availableKeys = append(availableKeys, strconv.Itoa(i+1))
		}

		for {
			char, key, err := keyboard.GetSingleKey()
			if err != nil {
				log.Fatal(err)
			}

			if key == keyboard.KeyEsc {
				break
			}

			if contains := slices.Contains(availableKeys, string(char)); !contains {
				fmt.Printf("The key you pressed is not in the option list! Available keys are: %v\n", strings.Join(availableKeys, ", "))
				continue
			} else {
				index, err := strconv.Atoi(string(char))
				if err != nil {
					fmt.Fprintf(os.Stderr, "Your key is invalid somehow: %v", err)
					continue
				}
				nextArc = Path(start.Options[index-1].Arc)
				break
			}
		}

		traverse(paths, nextArc)
	} else {
		var playAgain bool

		for {
			fmt.Printf("Congratulations! You've reached the end of this story! Press ENTER to play again or ESC to quit:\n")

			_, key, err := keyboard.GetSingleKey()
			if err != nil {
				log.Fatal(err)
			}

			if key == keyboard.KeyEsc {
				break
			}

			if key == keyboard.KeyEnter {
				playAgain = true
				break
			}
		}

		if playAgain {
			traverse(paths, "intro")
		}
	}
}

func main() {
	paths := make(Paths)

	f, err := os.Open("gophers.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&paths); err != nil {
		log.Fatal(err)
	}

	traverse(paths, "intro")
}
