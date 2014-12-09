package main

import (
	"bufio"
	"log"
	"os"
)

type PriceFile struct {
	Path string
}

func (pf PriceFile) lastLine() string {
	file, err := os.Open(pf.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}
	return line
}

func (pf PriceFile) addLine(line string) {
	file, err := os.OpenFile(pf.Path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(line + "\n")
}
