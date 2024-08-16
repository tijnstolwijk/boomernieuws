package main

import (
	"os"
	"strconv"
	"time"
  "slices"
	"boomernieuws/lib/pages"
	"fmt"
)

func main(){
  // Elk uur ?
  archiving := true
  args := os.Args[1:]
  if len(args) == 3 {
    if args[0] == "read" {
      archiving = false 
      timestamp := args[1]
      addr := args[2]
      if timestamp == "latest" {
        fmt.Print(readLatestPage(addr))
      } else {
        if timestampInt, err := strconv.ParseInt(timestamp, 0, 64); err == nil{
          fmt.Print(getPageWithTimestamp(addr, timestampInt))
        }else {fmt.Print("Invalid time, either a timestamp or latest")}
      }
      os.Exit(0)
    } 
  }
  if archiving { 
   ticker := time.NewTicker(time.Hour)
    quit := make(chan struct{})
    for {
      select {
        case <- ticker.C:
          downloadAll()
        case <- quit:
          ticker.Stop()
          break
      }
    } 
  }
}

func downloadAll(){
  fmt.Println("Running new download round")
  nextPage := "100"
  timestamp := time.Now().Unix()
  os.Mkdir(fmt.Sprintf("teletekst/%v", timestamp), os.FileMode(0777))
  // All this code does is crawl through all paths to get every page
  for nextPage != "" {
    page := download(nextPage, timestamp)
    nextPage = page.NextPage
    nextSubPage := page.NextSubPage
    for nextSubPage != "" {
      nextSubPage = download(nextSubPage, timestamp).NextSubPage
    }
  }

}
func download(pageAddr string, timestamp int64) pages.Teletekst_pagina{
  //downloads page and returns that page
  //TODO: we only want to download if the page isn't the same as the latest version of the same page
  //We fetch a page and process it, but don't save it to the system just yet
  //First we compare to the latest page that existed
  //We don't want to process our pages twice, once to compare, once to save
  //This requires us to overhaul the page management file
  page := pages.FetchPage(pageAddr)
  processedPage := pages.ParsePage(page)
  latestPage := readLatestPage(pageAddr)
  if processedPage != latestPage || latestPage == "Gabbagool? ova heeeeree.."{
    fmt.Printf("Downloading %v\n", pageAddr)
    path := fmt.Sprintf("teletekst/%v/%v.txt", timestamp, pageAddr)
    fmt.Print(processedPage)
    pages.SaveText(path, processedPage)

  } else {
    fmt.Println("Duplicate!!")
  }
  return page
}

func readLatestPage(pageAddr string) string {
  //traverse all timestamps, check if there is a page with our pageAddr
  //if not, we take the next timestamp
  timestamps := getAllTimestamps() 
  for i := 0; i < len(timestamps); i++{
    path := fmt.Sprintf("./teletekst/%v/%v.txt", timestamps[i], pageAddr)
    if _, err := os.Stat(path); err == nil{
      //this page exists (it must be the latest)
      bytes, err := os.ReadFile(path)
      if err != nil {
        fmt.Printf("Error reading file: %v\r\n", err.Error())
      }
      text := string(bytes)
      return text
    }
   }
  //if no page is found, scream bloody murder (jk we just download the page just in case)
  return "Gabbagool? ova heeeeree.."
}

func getAllTimestamps() []int64{
  // get all directories that could be a timestamp (can be cast to int)
  // then sort them from highest to lowest so we know how to traverse
  // the directories when we search for the latest page
  timestampDirs := []int64{}
  dirs, err := os.ReadDir("./teletekst/")
  if err != nil {
    fmt.Printf("Error reading directory: %v\r\n", err.Error())
  }
  for _, entry := range dirs {
    if timestampDir, err := strconv.ParseInt(entry.Name(), 0, 64); err == nil{
      timestampDirs = append(timestampDirs, timestampDir)
    }
  }
  slices.Sort(timestampDirs)
  slices.Reverse(timestampDirs)
  return timestampDirs
}

func getPageWithTimestamp(addr string, timestamp int64) string{
  timestamps := getAllTimestamps()
  if timestamp < timestamps[len(timestamps) - 1] {
    return "Too early, we haven't recorded that far back\r\n"
  }
  if timestamp > time.Now().Unix() {
    return "This is in the future stupid\r\n"
  }
  for i := 0; i < len(timestamps); i++{
    //This means we found the correct snapshot for the given timestamp
    if timestamp >= timestamps[i] {
      // We can now search for our page excluding all the higher timestamps
      timeAccurateTimestamps := timestamps[i:]
      for j := 0; j < len(timeAccurateTimestamps); j++{
        path := fmt.Sprintf("./teletekst/%v/%v.txt", timeAccurateTimestamps[j], addr)
        if _, err := os.Stat(path); err == nil{
          // This is the most recent page for this time
          bytes, err := os.ReadFile(path)
          if err != nil {
            fmt.Printf("Error reading file: %v\r\n", err.Error())
          }
          text := string(bytes)
          return text
        }
      }
    }
  }
  return "Gabbagool? ova heeeeree..\r\n"
}
