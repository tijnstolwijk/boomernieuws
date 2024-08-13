package main

import (
  //"time"
  //"os"
  //"syscall"
  //"os/signal"
  "boomernieuws/lib/pages"
  "fmt"
)

func main(){
  downloadAll()
  // Elk uur ?

  //ticker := time.NewTicker(time.Hour)
  //quit := make(chan struct{})
  //go func(){
  //  for {
  //    select {
  //      case <- ticker.C:
  //        go downloadAll()
  //      case <- quit:
  //        ticker.Stop()
  //        break
  //    }
  //  } 
  //}()

  ////
  //c := make(chan os.Signal)
  //signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  //go func() {
  //    <-c
  //    close(quit)
  //}()
}

func downloadAll(){
  nextPage := "100"

  for nextPage != "" {
    page := download(nextPage)
    nextPage = page.NextPage
    nextSubPage := page.NextSubPage
    //Download all subpages until there is no subpage left
    for nextSubPage != "" {
      nextSubPage = download(nextSubPage).NextSubPage
    }
  }

}

func download(pageAddr string) pages.Teletekst_pagina{
  page := pages.FetchPage(pageAddr)
  fmt.Printf("Downloading %v\n", pageAddr)
  path := fmt.Sprintf("teletekst/%v.txt", pageAddr)
  pages.PrintPage(page)
  pages.WritePageToFile(path, page)
  return page
}
