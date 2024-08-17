package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type locationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type command struct {
	name string
	desc string
	do   func(s *state) error
}

type state struct {
	mapIndex            int
	defaultMapIncrement int
	mapNext             string
	mapBefore           string
	exit                bool
	testing             bool
}

func tbd(s *state) error {
	fmt.Println("this will be developed soon")
	return nil
}

func helpCom(s *state) error {
	fmt.Println("this is help")
	return nil
}

func exitCom(s *state) error {
	fmt.Println("ma'a salaamah")
	s.exit = true
	return nil
}

// func init() {

// }
var commands map[string]command
var ss state
var responseCache map[string]stash

func main() {
	reader := bufio.NewReader(os.Stdin)

	normalInit()

	go func() {
		for {
			time.Sleep(15 * time.Second)
			cleanCache(20 * time.Second)
		}
	}()

	for {
		fmt.Print(Green + "pokedex > " + Reset)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		// fmt.Println("You entered:", input)

		if comm, ok := commands[input]; ok {
			comErr := comm.do(&ss)
			if comErr != nil {
				fmt.Println("Exitting because of error")
				fmt.Println(comErr)
				ss.exit = true
			}
		}

		if ss.exit {
			break
		}
	}
}
func normalInit() {

	responseCache = make(map[string]stash)
	ss = state{0, 20, "", "", false, false}
	commands = map[string]command{
		"help":  {"help", "display help message", helpCom},
		"exit":  {"exit", "exit the program", exitCom},
		"reset": {"reset", "reset the state of program if any", tbd},
		"map":   {"map", "list location, if without argument call again to call next available location", locMap},
		"mapb":  {"mapb", "reset the state of program if any", locMapB},
	}

}
