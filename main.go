package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var (
	inputPath    = "./inputs/logs/"
	outputPath   = "./outputs/"
	sheetName    = "Sheet1"
	scanKeywords = []string{"Exception"}
)

func main() {
	// set excel
	excel := excelize.NewFile()

	// read log file
	dir, err := ioutil.ReadDir(inputPath)
	if err != nil {
		panic(err)
	}

	// set titles
	_ = excel.SetCellStr(sheetName, "A1", "File Name")
	for inedx, keyword := range scanKeywords {
		axis, _ := excelize.CoordinatesToCellName(inedx+2, 1)
		_ = excel.SetCellStr(sheetName, axis, keyword)
	}

	// scan files
	for index, fileInfo := range dir {
		fmt.Println(inputPath + fileInfo.Name())
		file, err := os.Open(inputPath + fileInfo.Name())
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		// init keywordCounts
		keywordCounts := make(map[string]int)
		for _, keyword := range scanKeywords {
			keywordCounts[keyword] = 0
		}

		// search count
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			for _, keyword := range scanKeywords {
				if strings.Contains(text, keyword) {
					keywordCounts[keyword]++
				}
			}
		}

		fileNames := strings.Split(file.Name(), "/")
		fileName := fileNames[len(fileNames)-1]
		_ = excel.SetCellStr(sheetName, fmt.Sprintf("A%d", index+2), fileName)
		kindex := 0
		for _, keywordCount := range keywordCounts {
			axis, _ := excelize.CoordinatesToCellName(kindex+2, index+2)
			_ = excel.SetCellInt(sheetName, axis, keywordCount)
			kindex++
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	}

	// save excel
	if err = excel.SaveAs(fmt.Sprintf("%s%s.xlsx", outputPath, time.Now().Format("20060102 030405"))); err != nil {
		fmt.Println(err)
	}
}
