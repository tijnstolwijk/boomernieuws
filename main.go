package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
  "os/exec"
	"slices"
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
  args := os.Args[1:]
  if len(args) != 0{
    if args[0] == "--output"{
      if len(args) == 3{
        path := args[1]
        page := args[2]
        writePageToFile(path, page)
        os.Exit(0)
      }
    }
    
    page := args[0]
    printPage(page)
    os.Exit(0)
  }
  printPage("101")
  
  var input string
  for {
    fmt.Print("Pagina: ")
    fmt.Scan(&input)
    if input == "q"{os.Exit(0)}
    clearScreen()
    printPage(input)
    }
  }



func printPage(page string) {
  url := fmt.Sprintf("https://teletekst-data.nos.nl/json/%v", page)
  response, err := http.Get(url)
  if err != nil {
    fmt.Printf("Error getting to teletekst\n")
    os.Exit(1)
  }
  var pagina teletekst_pagina 
  if response.ContentLength == 0 {
    fmt.Printf("Page does not exist\n") 
    return
  }
  err = json.NewDecoder(response.Body).Decode(&pagina)
  lines := strings.Split(pagina.Content, "\n")
  for i := 0; i < len(lines); i++ {
    processed_line := processHTML(lines[i])
    processed_line = replaceBlockCharsR(processed_line)
    processed_line = replaceSpecialChars(processed_line)
    fmt.Printf("%v\n", processed_line)
  }
}

func writePageToFile(path string, page string){
  url := fmt.Sprintf("https://teletekst-data.nos.nl/json/%v", page)
  response, err := http.Get(url)
  if err != nil {
    fmt.Printf("Error getting to teletekst\n")
    os.Exit(1)
  }
  var pagina teletekst_pagina 
  if response.ContentLength == 0 {
    fmt.Printf("Page does not exist\n") 
    return
  }
  file, err := os.Create(path)
  if err != nil{
    panic(err)
  }
  defer file.Close()

  err = json.NewDecoder(response.Body).Decode(&pagina)
  lines := strings.Split(pagina.Content, "\n")
  for i := 0; i < len(lines); i++ {
    processed_line := processHTML(lines[i])
    processed_line = replaceBlockCharsR(processed_line)
    processed_line = replaceSpecialChars(processed_line)
    _, err := file.WriteString(fmt.Sprintf("%v\n", processed_line))
    if err != nil{
      panic(err)
    }
  }
  file.Sync()
}


func clearScreen(){
  cmd := exec.Command("clear")
  if runtime.GOOS == "windows"{
        cmd = exec.Command("cmd", "/c", "cls")
  }
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
      }
      if chars[i] == ">"{
        end = i
        chars = slices.Delete(chars, begin, end + 1)
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
  conversionMap := map[string]string {
    "&Agrave;": "À",
    "&Aacute;": "Á",
    "&Acirc;": "Â",
    "&Atilde;": "Ã",
    "&Auml;": "Ä",
    "&Aring;": "Å",
    "&AElig;": "Æ",
    "&Ccedil;": "Ç",
    "&Egrave;": "È",
    "&Eacute;": "É",
    "&Ecirc;": "Ê",
    "&Euml;": "Ë",
    "&Igrave;": "Ì",
    "&Iacute;": "Í",
    "&Icirc;": "Î",
    "&Iuml;": "Ï",
    "&ETH;": "Ð",
    "&Ntilde;": "Ñ",
    "&Ograve;": "Ò",
    "&Oacute;": "Ó",
    "&Ocirc;": "Ô",
    "&Otilde;": "Õ",
    "&Ouml;": "Ö",
    "&Oslash;": "Ø",
    "&Ugrave;": "Ù",
    "&Uacute;": "Ú",
    "&Ucirc;": "Û",
    "&Uuml;": "Ü",
    "&Yacute;": "Ý",
    "&THORN;": "Þ",
    "&szlig;": "ß",
    "&agrave;": "à",
    "&aacute;": "á",
    "&acirc;": "â",
    "&atilde;": "ã",
    "&auml;": "ä",
    "&aring;": "å",
    "&aelig;": "æ",
    "&ccedil;": "ç",
    "&egrave;": "è",
    "&eacute;": "é",
    "&ecirc;": "ê",
    "&euml;": "ë",
    "&igrave;": "ì",
    "&iacute;": "í",
    "&icirc;": "î",
    "&iuml;": "ï",
    "&eth;": "ð",
    "&ntilde;": "ñ",
    "&ograve;": "ò",
    "&oacute;": "ó",
    "&ocirc;": "ô",
    "&otilde;": "õ",
    "&ouml;": "ö",
    "&oslash;": "ø",
    "&ugrave;": "ù",
    "&uacute;": "ú",
    "&ucirc;": "û",
    "&uuml;": "ü",
    "&yacute;": "ý",
    "&thorn;": "þ",
    "&yuml;": "ÿ",
  }

  for k, v := range conversionMap {
    if strings.Contains(line, k){
      line = strings.ReplaceAll(line, k, v) 
    }
  }
  return line
}

