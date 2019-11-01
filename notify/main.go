package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		println("Please specifiy a directory path")
		os.Exit(1)
	}
	dirPath := os.Args[1]

	files, err := getFiles(dirPath)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	// randomly select a file
	selectedFile := files[rand.Intn(len(files))].Name()
	fileBytes, err := ioutil.ReadFile(filepath.Join(dirPath, selectedFile))
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	sentences := strings.Split(string(fileBytes), "\n\n")
	// randomly select a sentence
	sentenceIndex := rand.Intn(len(sentences))
	sentence := sentences[sentenceIndex]

	if len(os.Args) == 3 && os.Args[2] == "notify" {
		exec.Command("notify-send", "-t", "60000", sentence).Run()
	} else {
		println(sentence)
	}
}

func getFiles(dirPath string) ([]os.FileInfo, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	var result []os.FileInfo

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			result = append(result, fileInfo)
		}
	}

	return fileInfos, nil
}
