package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func changeFile(text string) string {
	r, _ := regexp.Compile(`[^.!?]*[.!?]`)
	return r.FindString(text)
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

func copyScan(fileRead *os.File, writeScanner *bufio.Writer) error {
	buffer := make([]byte, 1)
	var data = ""
	for {
		bytesread, err := fileRead.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		data += string(buffer[:bytesread])
		var changedData = changeFile(data)
		if changedData != "" {
			data = changedData
			break
		}
	}
	_, err := writeScanner.WriteString(data)
	if err != nil {
		return err
	}
	// move from buffer to file
	writeScanner.Flush()
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
	// create file if not exists
	fileWrite, err := os.Create(output + strings.TrimSuffix(filename, filepath.Ext(filename)) + ".res")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fileWrite.Close()
	// create write scanner for it
	writeScanner := bufio.NewWriter(fileWrite)

	err = copyScan(fileRead, writeScanner)
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
