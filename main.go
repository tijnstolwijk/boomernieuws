package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
  "os/exec"
	"slices"
	"strconv"
	"strings"
  "runtime"
	"github.com/h2so5/goback/regexp"
)

type teletekst_pagina struct{
  PrevPage string
  NextPage string
  PrevSubPage string
  NextSubPage string
  FastTextLinks  []fastTextLink
  Content string

}

type fastTextLink struct {
  Title string
  Page string
}

func main(){
  printPage(101)
  
  var input string
  for {
    fmt.Print("Pagina: ")
    fmt.Scan(&input)
    pageNum, err := strconv.ParseInt(input, 10, 0)
    if err == nil {
      clearScreen()
      printPage(int(pageNum))
    }else{
      if input == "q"{
        os.Exit(0)
      }
    }
  }
}

func printPage(page int) {
  url := fmt.Sprintf("https://teletekst-data.nos.nl/json/%d", page)
  response, err := http.Get(url)
  if err != nil {
    fmt.Printf("Error getting to teletekst\n")
    os.Exit(1)
  }
  var pagina teletekst_pagina 
  err = json.NewDecoder(response.Body).Decode(&pagina)
  lines := strings.Split(pagina.Content, "\n")
  for i := 0; i < len(lines); i++ {
    processed_line := processHTML(lines[i])
    processed_line = replaceBlockCharsR(processed_line)
    fmt.Printf("%v\n", processed_line)
  }
}

func clearScreen(){
  var cmd *exec.Cmd
  if runtime.GOOS == "windows"{
        cmd = exec.Command("cmd", "/c", "cls")
  }
  cmd = exec.Command("clear")
  cmd.Stdout = os.Stdout
  cmd.Run()
}

// Text processing functions
func processHTML(line string) string{
  // Remove all HTML tags
  chars := strings.Split(line, "")

  var begin int
  var end int

  for {
    if !slices.Contains(chars, "<") && !slices.Contains(chars, ">"){
      break
    }

    for i := 0; i < len(chars); i++{
      if chars[i] == "<"{
        begin = i
        //fmt.Printf("%v found at i=%d\n", chars[i], i)
      }
      if chars[i] == ">"{
        end = i
        //fmt.Printf("%v found at i=%d\n", chars[i], i)
        //fmt.Printf("Deleting from %v to %v\n", chars[begin], chars[end])
        chars = slices.Delete(chars, begin, end + 1)
        //fmt.Printf(strings.Join(chars, "")+"\n")
        break
      }
    }

  }
  return strings.Join(chars, "")
}

func replaceBlockCharsR(line string) string {
  blockChars := regexp.MustCompile("&#xF0(?!20;).{3}")
  result := blockChars.ReplaceAll([]byte(line), []byte("%"))
  blackBlockChars := regexp.MustCompile("&#xF020;")
  result = blackBlockChars.ReplaceAll(result, []byte(" "))
  return string(result)
}

func replaceSpecialChars(line string) string{
  // https://www.html.am/reference/html-special-characters.cfm ISO 8859-1 section
  return ""
}


