package core

import (
	"bytes"
	"github.com/xuri/excelize/v2"
)

func ImportExcel() {
	excelize.OpenReader(bytes.NewReader(nil))
}
