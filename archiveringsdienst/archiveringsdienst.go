package main

import (
  "time"
  "os"
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
  timestamp := time.Now().Unix()
  os.Mkdir(fmt.Sprintf("teletekst/%v", timestamp), os.FileMode(0777))
  for nextPage != "" {
    //downloads page then derives from it the nextpage and the subpage
    page := download(nextPage, timestamp)
    nextPage = page.NextPage
    nextSubPage := page.NextSubPage
    //Download all subpages until there is no subpage left
    for nextSubPage != "" {
      nextSubPage = download(nextSubPage, timestamp).NextSubPage
    }
  }

}

func download(pageAddr string, timestamp int64) pages.Teletekst_pagina{
  //downloads page and returns that page
  //TODO: we only want to download if the page isn't the same in the last snapshot
  page := pages.FetchPage(pageAddr)
  fmt.Printf("Downloading %v\n", pageAddr)
  path := fmt.Sprintf("teletekst/%v/%v.txt", timestamp, pageAddr)
  pages.PrintPage(page)
  pages.WritePageToFile(path, page)
  return page
}

