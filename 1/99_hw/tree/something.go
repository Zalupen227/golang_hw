package main

import "fmt"

func myFunc() {
	mySlice := make([]int, 0, 20)
	myMap := make(map[int]string)
	myMap = map[int]string{
		2: ",",
	}
	fmt.Println(mySlice)
}
