package main

import (
	"fmt"
	"io/ioutil"
)

// Resources struct containg file list
type Resources struct {
	Images []string
	Sounds []string
}

// GetFileList return list of file of certain file
func GetFileList() Resources {
	result := Resources{}
	result.Images = readDirectory(imagesPath)
	result.Sounds = readDirectory(soundsPath)

	return result
}

func readDirectory(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	var result []string

	for _, file := range files {
		result = append(result, file.Name())
	}

	return result
}
