package core

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"testing"
	"text/template"
)

// 示例数据结构
type Employee struct {
	Name       string
	Age        int
	Department string
}

func TestTemplate(t *testing.T) {
	data := Employee{
		Name:       "张三",
		Age:        28,
		Department: "技术部",
	}

	// 填充模板并导出
	err := FillExcelTemplate("C:\\Users\\aixot\\Desktop\\template.xlsx", "C:\\Users\\aixot\\Desktop\\output.xlsx", data)
	if err != nil {
		fmt.Println("导出失败:", err)
		return
	}
	fmt.Println("导出成功!")
}
func FillExcelTemplate(templatePath, outputPath string, data interface{}) error {
	// 打开模板文件
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// 遍历所有工作表
	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return err
		}

		// 遍历有数据的行和列
		for rowIdx, row := range rows {
			for colIdx, cellValue := range row {
				// 转换行列索引为单元格坐标（如 "A1"）
				axis, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
				if err != nil {
					return err
				}

				// 检查单元格是否包含模板语法
				if containsTemplate(cellValue) {
					// 渲染模板
					rendered, err := renderTemplate(cellValue, data)
					if err != nil {
						continue // 或记录错误
					}

					// 更新单元格值（保留原有样式）
					f.SetCellValue(sheet, axis, rendered)
				}
			}
		}
	}

	// 保存为新文件
	if err := f.SaveAs(outputPath); err != nil {
		return err
	}
	return nil
}

// 检查是否包含模板占位符
func containsTemplate(s string) bool {
	return len(s) > 4 && s[0:2] == "{{" && s[len(s)-2:] == "}}"
}

// 使用text/template渲染内容
func renderTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("cell").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
