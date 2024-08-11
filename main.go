package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

var commands map[string]command

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
	exit                bool
	mapIndex            int
	defaultMapIncrement int
	mapNext             string
	mapBefore           string
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

func FetchAndUnmarshal[T any](url string) (T, error) {
	var result T
	resp, err := http.Get(url)
	if err != nil {
		return result, fmt.Errorf("failed to fetch data from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	// fmt.Println("result should be returned", result)
	return result, nil
}

func locMap(s *state) error {

	endpoint := s.mapNext
	if endpoint == "" {
		endpoint = "http://pokeapi.co/api/v2/location"
	}

	locRes, err := FetchAndUnmarshal[locationResponse](endpoint)
	if err != nil {
		return err
	}

	// fmt.Println(locRes.Next, locRes.Previous)
	// fmt.Printf("%+v \n", locList)

	for _, location := range locRes.Results {
		fmt.Println(location.Name)
	}

	s.mapNext = locRes.Next
	s.mapBefore = endpoint

	return nil
}

func locMapB(s *state) error {

	endpoint := s.mapBefore
	if endpoint == "" {
		endpoint = "http://pokeapi.co/api/v2/location"
	}

	locRes, err := FetchAndUnmarshal[locationResponse](endpoint)
	if err != nil {
		return err
	}

	for _, location := range locRes.Results {
		fmt.Println(location.Name)
	}

	s.mapNext = locRes.Next
	s.mapBefore = locRes.Previous

	return nil
}

var ss state

func init() {
	ss = state{false, 0, 20, "", ""}
	commands = map[string]command{
		"help":  {"help", "display help message", helpCom},
		"exit":  {"exit", "exit the program", exitCom},
		"reset": {"reset", "reset the state of program if any", tbd},
		"map":   {"map", "list location, if without argument call again to call next available location", locMap},
		"mapb":  {"mapb", "reset the state of program if any", locMapB},
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

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

		// if input == "exit" || input == "EXIT" {
		// 	fmt.Println("ma'a salamah")
		// 	break
		// }
		if ss.exit {
			break
		}
	}

}
