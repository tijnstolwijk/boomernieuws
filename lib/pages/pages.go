package pages


import (
  "boomernieuws/lib/process"
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
  return pagina
}

func PrintPage(pagina Teletekst_pagina) {
  lines := strings.Split(pagina.Content, "\n")
  for i := 0; i < len(lines); i++ {
    processed_line := process.ProcessHTML(lines[i])
    processed_line = process.ReplaceBlockCharsR(processed_line)
    processed_line = process.ReplaceSpecialChars(processed_line)
    fmt.Printf("%v\r\n", processed_line)
  }
}

func WritePageToFile(path string, pagina Teletekst_pagina){
  file, err := os.Create(path)
  if err != nil{
    panic(err)
  }
  defer file.Close()

  lines := strings.Split(pagina.Content, "\n")
  for i := 0; i < len(lines); i++ {
    processed_line := process.ProcessHTML(lines[i])
    processed_line = process.ReplaceBlockCharsR(processed_line)
    processed_line = process.ReplaceSpecialChars(processed_line)
    _, err := file.WriteString(fmt.Sprintf("%v\n", processed_line))
    if err != nil{
      panic(err)
    }
  }
  file.Sync()
}

