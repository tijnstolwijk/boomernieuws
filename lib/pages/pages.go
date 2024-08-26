package pages

import (
	"github.com/tijnstolwijk/boomernieuws/lib/process"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)
type Teletekst_pagina struct{
  PrevPage string
  NextPage string
  PrevSubPage string
  NextSubPage string
  FastTextLinks  []fastTextLink
  Content string
  SelfPage string
}

type fastTextLink struct {
  Title string
  Page string
}

func FetchPage(page string) Teletekst_pagina{
  var pagina Teletekst_pagina
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
  pagina.SelfPage = page
  return pagina
}

func ParsePage(pagina Teletekst_pagina) string{
  lines := strings.Split(pagina.Content, "\n")
  for i := 0; i < len(lines); i++ {
    processed_line := process.ProcessHTML(lines[i])
    processed_line = process.ReplaceBlockCharsR(processed_line)
    processed_line = process.ReplaceSpecialChars(processed_line)
    lines[i] = processed_line
  }
  return strings.Join(lines, "\r\n")
}

func SaveText(path string, text string){
  file, err := os.Create(path)
  if err != nil{
    panic(err)
  }
  defer file.Close()
  
  _, err = file.WriteString(text)
  if err != nil{
    panic(err)
  }
  file.Sync()
}


//High-level functions used by boomernieuws

func SavePage(pageAddr string, path string) {
  page := FetchPage(pageAddr)
  text := ParsePage(page)
  SaveText(path, text)
}

func PrintPage(pageAddr string) Teletekst_pagina{
  page := FetchPage(pageAddr)
  text := ParsePage(page)
  fmt.Print(text)
  return page
}
