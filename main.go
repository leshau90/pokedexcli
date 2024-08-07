package main

import (
	"bufio"
	"fmt"
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
type command struct{
  name string
  desc string  
  do func()error
}

func tbd()error {
  fmt.Println("this will be developed soon")
  return nil
}

func helpCom()error{
   fmt.Println("this is help")
   return nil
}

func init (){
  commands = map[string]command{
    "help"  : {"help","display help message",tbd},
    "exit"  : {"exit","exit the program",tbd},
    "reset" : {"reset","reset the state of program if any",tbd},
  }
}

func main() {
  reader := bufio.NewReader(os.Stdin)

  for {
    fmt.Print( Green +"pokedex > "+Reset)
    input, _ := reader.ReadString('\n')
    input =  strings.TrimSpace(input)

    if input == "exit" || input == "EXIT" {
      fmt.Println("ma'a salamah")
      break
    }
    fmt.Println("You entered:",input)
  }

}
