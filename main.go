package main

import (
	"github.com/tijnstolwijk/boomernieuws/lib/pages"
	"fmt"
	"os"
	"os/exec"
	"runtime"

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
        addr := args[2]
        pages.SavePage(addr, path)
        runtime.Goexit()
      }
    }

    addr := args[0]
    pages.PrintPage(addr)
    runtime.Goexit()
  }

  //Standaardpagina
  curPage := pages.PrintPage("101")

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
    curPage = getNewPage(input, curPage)
  }
}

func getNewPage(input string, curPage pages.Teletekst_pagina) pages.Teletekst_pagina{
  var pagina pages.Teletekst_pagina
  switch input{
    case "h":
      if curPage.PrevPage == "" {
        fmt.Printf("Page does not exist\r\n") 
        return pages.PrintPage(curPage.SelfPage)
      }
      pagina = pages.PrintPage(curPage.PrevPage)
    case "l":
      if curPage.NextPage == "" {
        fmt.Printf("Page does not exist\r\n") 
        return pages.PrintPage(curPage.SelfPage)
      }
      pagina = pages.PrintPage(curPage.NextPage)
    case "j":
      if curPage.NextSubPage == "" {
        fmt.Printf("Page does not exist\r\n") 
        return pages.PrintPage(curPage.SelfPage)
      }
      pagina = pages.PrintPage(curPage.NextSubPage)
    case "k":
      if curPage.PrevSubPage == "" {
        fmt.Printf("Page does not exist\r\n") 
        return pages.PrintPage(curPage.SelfPage)
      }
      pagina = pages.PrintPage(curPage.PrevSubPage)
    default:
      pagina = pages.PrintPage(input)
  }
  return pagina
}

func clearScreen(){
  cmd := exec.Command("clear")
  if runtime.GOOS == "windows"{
        cmd = exec.Command("cmd", "/c", "cls")
  }
  cmd.Stdout = os.Stdout
  cmd.Run()
}
