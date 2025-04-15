package core

import (
	"fmt"
	"github.com/dop251/goja"
	"testing"
)

type LuckySheet struct {
	Ct struct {
		Fa string `json:"fa"`
		T  string `json:"t"`
	} `json:"ct"`
	V  int    `json:"v"`
	M  int    `json:"m"`
	Bg string `json:"bg"`
	Ff int    `json:"ff"`
	Fc string `json:"fc"`
	Bl int    `json:"bl"`
	It int    `json:"it"`
	Fs int    `json:"fs"`
	Cl int    `json:"cl"`
	Ht int    `json:"ht"`
	Vt int    `json:"vt"`
	Tr int    `json:"tr"`
	Tb int    `json:"tb"`
	Ps struct {
		Left   int    `json:"left"`
		Top    int    `json:"top"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Value  string `json:"value"`
		Isshow bool   `json:"isshow"`
	} `json:"ps"`
	F string `json:"f"`
}

func TestScript(t *testing.T) {
	vm := goja.New()
	// 定义JavaScript对象
	script := `
		var obj = { a: [1, 2, 3], b: [4, 5, 6] };
		var result = [];
		for (var key in obj) {
			var list = obj[key];
			for (var i = 0; i < list.length; i++) {
				result.push({ key: key, value: list[i] });
			}
		}
		result; // 返回结果数组
	`
	v, err := vm.RunString(script)
	if err != nil {
		panic(err)
	}

	// 导出结果为[]interface{}
	result := v.Export().([]interface{})
	for _, item := range result {
		entry := item.(map[string]interface{})
		fmt.Printf("键: %s, 值: %v\n", entry["key"], entry["value"])
	}
}
