package main

import (
  "slices"
  "fmt"
)

func main(){
  intList := []int{0,1,2,3,4,5}
  for i := 0; i < len(intList); i++{
    intList = slices.Delete(intList, len(intList)-1, len(intList)-1)
  }
  fmt.Printf("%v\n", intList)
}
