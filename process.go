package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func changeFile(text string) string {
	r, _ := regexp.Compile(`[^.!?]*[.!?]`)
	found := r.FindString(text)
	return strings.Replace(text, found, "\""+found+"\"", 1)
}

func getCLIArgs() (string, string) {
	inputFolder := os.Args[1]
	outputFolder := os.Args[2]
	return inputFolder, outputFolder
}

func listDirectory(directory string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if path != directory {
			files = append(files, strings.TrimPrefix(path, directory))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func copyScan(readScanner *bufio.Scanner, writeScanner *bufio.Writer) error {
	var parsed = false
	for readScanner.Scan() {
		var text = readScanner.Text()
		if !parsed {
			text = changeFile(text)
			parsed = true
		}
		_, err := writeScanner.WriteString(text)
		if err != nil {
			return err
		}
		writeScanner.WriteByte('\n')
		// move from buffer to file
		writeScanner.Flush()
	}

	if err := readScanner.Err(); err != nil {
		return err
	}
	return nil
}

func copyFile(input string, output string, filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	// open src file
	fileRead, err := os.Open(input + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fileRead.Close()
	// create read scanner for it
	readScanner := bufio.NewScanner(fileRead)
	// create file if not exists
	fileWrite, err := os.Create(output + strings.TrimSuffix(filename, filepath.Ext(filename)) + ".res")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fileWrite.Close()
	// create write scanner for it
	writeScanner := bufio.NewWriter(fileWrite)

	err = copyScan(readScanner, writeScanner)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func copyDirectory(input string, output string) error {
	var wg sync.WaitGroup
	files, err := listDirectory(input)
	if err != nil {
		return err
	}
	wg.Add(len(files))
	for _, file := range files {
		go copyFile(input, output, file, &wg)
	}
	wg.Wait()
	fmt.Print("Total number of processed files: ")
	fmt.Print(len(files))
	fmt.Print("\n")
	return nil
}

func main() {
	input, output := getCLIArgs()
	err := copyDirectory(input, output)
	if err != nil {
		fmt.Println(err)
	}
}
