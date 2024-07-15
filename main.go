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
  "golang.org/x/term"
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
  //runtime.Goexit() is handig omdat alle defer statements worden aangeroepen, maar speelt niet goed met terminals
  defer os.Exit(0)

  //Voor wanneer er cmdline options zijn
  args := os.Args[1:]
  if len(args) != 0{
    if args[0] == "--output"{
      if len(args) == 3{
        path := args[1]
        page := args[2]
        writePageToFile(path, page)
        runtime.Goexit()
      }
    }

    page := args[0]
    printPage(fetchPage(page))
    runtime.Goexit()
  }

  //Standaardpagina
  curPage := fetchPage("101")
  printPage(curPage)

  //Raw is nodig voor de interactieve modus
  oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
  if err != nil {panic(err)}
  defer term.Restore(int(os.Stdin.Fd()), oldState)

  var input string
  for {
    term.MakeRaw(int(os.Stdin.Fd()))

    byte := make([]byte, 1)
    _, err2 := os.Stdin.Read(byte)
    if err2 != nil {panic(err2)}
    input = string(byte)
     
    if input == ":"{
      //"command mode" ingaan, niet "interactive mode"
      term.Restore(int(os.Stdin.Fd()), oldState)
      fmt.Print(":")
      fmt.Scanln(&input)
      if input == "q"{runtime.Goexit()}
      if input == "\n"{break}
    }
    clearScreen()
    newPage := getNewPage(input, curPage)

    // Als de nieuwe pagina lege content heeft, blijven we op de huidige pagina
    if newPage.Content != ""{
      curPage = newPage
    }
    printPage(curPage)
  }
}

// Het hele idee is dat niet bestaande pagina's altijd als een lege struct worden gegeven
func getNewPage(input string, curPage teletekst_pagina) teletekst_pagina{
  var pagina teletekst_pagina
  switch input{
    case "h":
      if curPage.PrevPage == "" {
        fmt.Printf("Page does not exist\r\n") 
        return pagina
      }
      pagina = fetchPage(curPage.PrevPage)
    case "l":
      if curPage.NextPage == "" {
        fmt.Printf("Page does not exist\r\n") 
        return pagina
      }
      pagina = fetchPage(curPage.NextPage)
    case "j":
      if curPage.NextSubPage == "" {
        fmt.Printf("Page does not exist\r\n") 
        return pagina
      }
      pagina = fetchPage(curPage.NextSubPage)
    case "k":
      if curPage.PrevSubPage == "" {
        fmt.Printf("Page does not exist\r\n") 
        return pagina
      }
      pagina = fetchPage(curPage.PrevSubPage)
    default:
      pagina = fetchPage(input)
  }
  return pagina
}


func fetchPage(page string) teletekst_pagina{
  var pagina teletekst_pagina
  url := fmt.Sprintf("https://teletekst-data.nos.nl/json/%v", page)
  response, err := http.Get(url)
  if err != nil {
    fmt.Printf("Error getting to teletekst\r\n")
    return pagina
  }
  if response.ContentLength == 0 {
    fmt.Printf("Page does not exist\r\n") 
    return pagina
  }
  err = json.NewDecoder(response.Body).Decode(&pagina)
  return pagina
}

func printPage(pagina teletekst_pagina) {
  lines := strings.Split(pagina.Content, "\n")
  for i := 0; i < len(lines); i++ {
    processed_line := processHTML(lines[i])
    processed_line = replaceBlockCharsR(processed_line)
    processed_line = replaceSpecialChars(processed_line)
    fmt.Printf("%v\r\n", processed_line)
  }
}

func writePageToFile(path string, page string){
  url := fmt.Sprintf("https://teletekst-data.nos.nl/json/%v", page)
  response, err := http.Get(url)
  if err != nil {
    fmt.Printf("Error getting to teletekst\r\n")
    runtime.Goexit()
  }
  var pagina teletekst_pagina 
  if response.ContentLength == 0 {
    fmt.Printf("Page does not exist\r\n") 
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

