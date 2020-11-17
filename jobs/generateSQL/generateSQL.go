package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/NedHsu/golang-readlog/utilities"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/xormplus/xorm"
)

var (
	inputPath  = "./inputs/"
	outputPath = "./outputs/"
	sheetName  = "Sheet1"
	maxOnce    = 60000
	config     = utilities.InitConfigure()
	columns    = []string{
		"userName", "t1", "t2", "t3",
	}
)

func main() {
	engine, err := xorm.NewMSSQL("mssql", config.GetString("MSSQL.ConnectionString"))
	CheckErr(err)
	defer engine.Close()

	f, err := excelize.OpenFile(inputPath + "sql01.xlsx")
	CheckErr(err)

	fout := excelize.NewFile()

	// Headers
	for i, col := range columns {
		axis, _ := excelize.CoordinatesToCellName(i+1, 1)
		CheckErr(fout.SetCellValue(sheetName, axis, col))
	}

	// Get all the rows
	rows, _ := f.GetRows(sheetName)
	totalRow := len(rows)
	sqlTemplate := GetTemplate()
	rowNum := 1
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
		sqlScript := fmt.Sprintf(sqlTemplate, content)
		CreateFile(fmt.Sprintf("%ssql01_%d.sql", outputPath, i+1), sqlScript)
		results, err := engine.QueryInterface(sqlScript)
		CheckErr(err)
		for ri, record := range results {
			for ci, col := range columns {
				axis, _ := excelize.CoordinatesToCellName(ci+1, rowNum+ri+1)
				CheckErr(fout.SetCellValue(sheetName, axis, record[col]))
			}
		}
		rowNum += len(results)
	}
	CheckErr(fout.SaveAs(fmt.Sprintf("%s%s.xlsx", outputPath, time.Now().Format("20060102 030405"))))
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
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
