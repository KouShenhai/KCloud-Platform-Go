package core

import (
	"fmt"
	"github.com/dop251/goja"
	"testing"
	"time"
)

// 预编译JS函数提升性能
var precompiled *goja.Program

func init() {
	// 初始化阶段预编译脚本
	jsCode := `
    function processMap(map) {
        let sum = 0;
        for (const key in map) {
            sum += map[key];
        }
        return {
            total: sum,
            avg: sum / Object.keys(map).length,
            timestamp: Date.now()
        };
    }`

	var err error
	precompiled, err = goja.Compile("mapProcessor.js", jsCode, true)
	if err != nil {
		panic(err)
	}
}

func TestJs(t *testing.T) {
	// 1. 创建可复用的虚拟机实例池
	vmPool := make(chan *goja.Runtime, 5)
	for i := 0; i < 5; i++ {
		vm := goja.New()
		vm.RunProgram(precompiled) // 加载预编译代码
		vmPool <- vm
	}

	// 2. 准备测试数据
	data := map[string]int{
		"a": 10,
		"b": 20,
		"c": 30,
	}

	// 3. 执行JS处理
	start := time.Now()
	result := executeJS(vmPool, data)
	elapsed := time.Since(start)

	fmt.Printf("结果: %+v\n", result)
	fmt.Printf("耗时: %v\n", elapsed.Milliseconds())
}

// 带缓存的执行函数
func executeJS(pool chan *goja.Runtime, data map[string]int) map[string]interface{} {
	vm := <-pool
	defer func() { pool <- vm }()

	// 将Go Map转换为JS对象（优化点）
	jsMap := vm.NewObject()
	for k, v := range data {
		jsMap.Set(k, v)
	}

	// 调用函数
	fn, ok := goja.AssertFunction(vm.Get("processMap"))
	if !ok {
		panic("函数未找到")
	}

	// 执行并获取结果
	res, err := fn(goja.Undefined(), jsMap)
	if err != nil {
		panic(err)
	}

	// 将JS对象转换为Go Map
	return exportMap(res.ToObject(vm))
}

// 深度转换JS对象到Go Map（递归处理嵌套结构）
func exportMap(obj *goja.Object) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range obj.Keys() {
		val := obj.Get(key)
		if subObj, ok := val.(*goja.Object); ok {
			// 处理嵌套对象
			result[key] = exportMap(subObj)
		} else {
			result[key] = val.Export()
		}
	}
	return result
}
