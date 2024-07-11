package main

import (
  "encoding/json"
  "net/http"
  "os"
  "fmt"
  "slices"
  "strings"
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
  url := "https://teletekst-data.nos.nl/json/101" 
  response, err := http.Get(url)
  if err != nil {
    fmt.Printf("Error getting to teletekst\n")
    os.Exit(1)
  }
  var pagina teletekst_pagina 
  err = json.NewDecoder(response.Body).Decode(&pagina)
  lines := strings.Split(pagina.Content, "\n")
  //for i := 0; i < len(lines); i++ {
  //  processed_line := processHTML(lines[i])
  //  fmt.Printf(processed_line)
  //}
  fmt.Printf(lines[0] + "\n")
  fmt.Printf(processHTML(lines[0]) + "\n")

}

func processHTML(line string) string{
  // Remove all HTML tags
  chars := strings.Split(line, "")
  var begin int
  var end int
  for {
    if !slices.Contains(chars, "<"){
      break
    }
    for i := 0; i < len(chars); i++{
      if chars[i] == "<"{
        begin = i
      }
      if chars[i] == ">"{
        end = i
        chars = slices.Delete(chars, begin, end + 1)
      }
    }

  }
  return strings.Join(chars, "")
}

