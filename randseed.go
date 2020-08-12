package main

import (
    "fmt"
    "math/rand"
    "time"
)
        
func main() {
    i:=10
    for i<=10{
    rand.Seed(time.Now().UnixNano())
    min := 5
    max := 10
    fmt.Println(rand.Intn(max - min + 1) + min)
    i=i-1
    if i==0{
    break}
    }
}
