package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var (
	inputPath  = "./inputs/"
	outputPath = "./outputs/"
	sheetName  = "Sheet1"
	maxOnce    = 60000
)

func main() {
	f, err := excelize.OpenFile(inputPath + "sql01.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get all the rows
	rows, _ := f.GetRows(sheetName)
	rows = rows[1:]
	totalRow := len(rows)
	sqlTemplate := GetTemplate()

	for i := 0; i < totalRow; i += maxOnce {
		var b strings.Builder
		last := i + maxOnce
		if last > totalRow {
			last = totalRow
		}

		for j := i; j < last; j++ {
			fmt.Fprintf(&b, "('%s'),\r\n", rows[j][0])
		}
		content := b.String()
		content = content[0 : len(content)-3]
		CreateFile(fmt.Sprintf("%ssql01_%d.sql", outputPath, i+1), fmt.Sprintf(sqlTemplate, content))
	}
}

func GetTemplate() string {
	content, err := ioutil.ReadFile(inputPath + "sql01.txt")
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func CreateFile(fileName, content string) {
	f2, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	_, err = f2.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
}
