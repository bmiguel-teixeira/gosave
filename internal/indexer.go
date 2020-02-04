package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
)

type FileData struct {
	index int
	path  string
}

func main() {
	var index int = 0
	var myFiles []FileData = []FileData{}

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				temp := FileData{index, path}
				myFiles = append(myFiles, temp)
				index++
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	//setup paralelism
	var wg sync.WaitGroup

	chunks := 12
	arraySize := len(myFiles)
	var chunkSize int = (int)(math.Floor(float64(arraySize / chunks)))
	finalChunk := arraySize - (chunkSize * chunks)

	for i := 0; i < chunks; i++ {
		wg.Add(1)
		log.Println(fmt.Sprintf("launching: %d", i))
		go PrintIt(&wg, i, myFiles[i*chunkSize:(i+1)*chunkSize-1])
	}
	if finalChunk > 0 {
		PrintIt(&wg, chunks, myFiles[arraySize-finalChunk-1:arraySize-1])
	}
	wg.Wait()
}

func PrintIt(wg *sync.WaitGroup, index int, files []FileData) {
	defer wg.Done()
	for _, elem := range files {
		fmt.Sprintf("%d -> %s -> %s", elem.index, elem.path, HashIt(elem.path))
		//log.Println(txt)
	}
	log.Println(fmt.Sprintf("worker done: %d| Processed: %d", index, len(files)))
}

func HashIt(filename string) string {
	hasher := sha256.New()
	s, err := ioutil.ReadFile(filename)
	hasher.Write(s)
	if err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}
