package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fileName := "birthday_001.txt"
	// => Birthday - 1 of 4.txt
	newName, err := match(fileName)
	if err != nil {
		fmt.Println("no match")
		os.Exit(1)
	}
	fmt.Println(newName)
}

// match returns the new name for the given file name
// or an error if no match is found

func match(fileName string) (string, error) {
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]
	fileName = strings.Join(pieces[0:len(pieces)-1], ".")
	pieces = strings.Split(fileName, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match our pattern", fileName)
	}
	// Birthday - 1.txt
	return fmt.Sprintf("%s - %d.%s", name, number, ext), nil

}
