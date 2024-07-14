package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	inputDir := "img"
	outputDir := "thumbnails"

	now := time.Now()
	processImages(inputDir, outputDir)
	fmt.Println(time.Now().Sub(now).Milliseconds())
}

func processImages(inputDir, outputDir string) {

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	err = os.MkdirAll(outputDir, os.ModePerm)

	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			go processThumbnails(filepath.Join(inputDir, file.Name()), outputDir, &wg)
		}
	}

	wg.Wait()

	fmt.Println("Img processed:", len(files))
}

func processThumbnails(fileInput, outputDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	files, err := os.Open(fileInput)

	if err != nil {
		panic(err)
	}

	defer files.Close()

	img, _, err := image.Decode(files)
	if err != nil {
		panic(err)
	}

	thumbnail := resize.Resize(100, 0, img, resize.Lanczos3)

	outputPath := filepath.Join(outputDir, filepath.Base(fileInput))

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}

	defer outputFile.Close()

	switch strings.ToLower(filepath.Ext(fileInput)) {
	case ".jpg", ".jpeg":
		err := jpeg.Encode(outputFile, thumbnail, nil)
		if err != nil {
			panic(err)
		}
	case ".png":
		err := png.Encode(outputFile, thumbnail)
		if err != nil {
			panic(err)
		}
	}
}
