package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/fs"
	"os"
)

// 增加合并单元格的文件导出
func WriteExcelMerge(datas map[string][][]string) (string, error) {
	file := excelize.NewFile()

	for sheetName, data := range datas {
		index, _ := file.NewSheet(sheetName)
		for i, row := range data {
			for j, val := range row {
				// 列行数字索引转excel坐标索引
				cellName, _ := excelize.CoordinatesToCellName(j+1, i+1)
				// 设置，写入
				file.SetCellValue(sheetName, cellName, val)
			}
		}
		// 创建表格
		file.SetActiveSheet(index)
	}

	for sheet, _ := range datas {
		// 获取表格中的所有行
		rows, err := file.GetRows(sheet)
		if err != nil {
			return "", err
		}

		for i := 1; i < len(rows); i++ {
			currRow := rows[i]
			prevRow := rows[i-1]

			// 判断前3列是否相等
			if currRow[0] == prevRow[0] &&
				currRow[1] == prevRow[1] &&
				currRow[2] == prevRow[2] {

				// 合并相邻单元格
				for j := 1; j <= 3; j++ {
					cellName1, _ := excelize.CoordinatesToCellName(j, i) // 列，行
					cellName2, _ := excelize.CoordinatesToCellName(j, i+1)
					err := file.MergeCell(sheet, cellName1, cellName2)
					if err != nil {
						return "", err
					}
				}
			}
		}
	}

	filename := "aaa" + ".xlsx"

	// 创建目录
	_, err := os.ReadDir("aaa/")
	if err != nil {
		// 不存在就创建
		err = os.MkdirAll("aaa/", fs.ModePerm)
		if err != nil {
			return "", err
		}
	}

	file.DeleteSheet("Sheet1")

	err = file.SaveAs("aaa/" + filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func main() {

	datas := make(map[string][][]string)
	datas["xxx"] = [][]string{{"a", "b", "c", "D1", "E1", "G1", "H1", "I1"}, {"a", "b", "c", "D1"}, {"a", "b", "c", "D1"}}
	datas["yyy"] = [][]string{{"X1", "Y1", "Z1"}, {"X2", "Y2", "Z2"}, {"X3", "Y3", "Z3"}}
	fileName, err := WriteExcelMerge(datas)
	if err != nil {
		fmt.Println("Write excel error: ", err)
		return
	}

	fmt.Println("Write excel success, file name is: ", fileName)
}
