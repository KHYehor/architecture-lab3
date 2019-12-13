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

// var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func checkSentence(prev byte, cur byte) bool {
	if string(cur) != " " {
		return false
	}
	match, _ := regexp.MatchString("[!?.]", string(prev))
	if match {
		return true
	}
	return false
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
	var prev byte = ' '
	buffer := make([]byte, 1)
	for {
		// Read from file to buffer
		_, err := fileRead.Read(buffer)
		// Error handling
		if err != nil {
			if err != io.EOF {
				return err
			}
			fmt.Println(err)
			break
		}
		// Check if sentence has finished
		if checkSentence(prev, buffer[0]) {
			break
		}
		// Replace current char to previous
		prev = buffer[0]
		// Write char to write buffer
		_, err = writeScanner.WriteString(string(buffer[0]))
		if err != nil {
			return err
		}
		// move from buffer to file
		writeScanner.Flush()
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
