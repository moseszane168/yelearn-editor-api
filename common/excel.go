/**
 * Excel处理工具类
 */

package common

import (
	"crf-mold/base"
	minioutil "crf-mold/common/minio"
	"crf-mold/dao"
	"crf-mold/model"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/xuri/excelize/v2"
)

/**
 * 导出Excel
 */
func ExportExcel(fileName, sheetName string, vo interface{}) minio.UploadInfo {
	f := excelize.NewFile()

	if sheetName == "" {
		sheetName = "Sheet1"
	}

	// Create a new sheet.
	index := f.NewSheet(sheetName)

	_, tags, values := base.GetVoFieldTypesValues(vo, "excel", "dict")

	// A1-Z1:请求头，先只用26列，多了有问题.TODO
	var rowStart byte = 'A'

	// excel表头
	var curRow int64 = 1 // 行
	for i := 0; i < len(tags); i++ {
		index := rowStart + byte(i)
		v := tags[i]
		cellStr := string(index) + strconv.Itoa(int(curRow))
		f.SetCellValue(sheetName, cellStr, v["excel"])
	}

	// excel内容
	for i := 0; i < len(values); i++ {
		curRow++
		row := values[i]
		for j := 0; j < len(row); j++ {
			index := rowStart + byte(j)
			dictKey := tags[j]["dict"]
			var v interface{}
			if dictKey == "" {
				v = row[j]
			} else {
				var m model.DictProperty
				if err := dao.GetConn().Table("dict_property").Where("is_deleted = 'N' and group_code = ? and `key` = ?", dictKey, row[j]).First(&m).Error; err != nil {
					v = row[j]
				} else {
					// TODO:根据语言判断导出什么数据
					v = m.ValueCn
				}
			}

			cellStr := string(index) + strconv.Itoa(int(curRow))
			f.SetCellValue(sheetName, cellStr, v)
		}
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs(fileName); err != nil {
		panic(err)
	}

	defer func() {
		// 删除本地文件
		err := os.Remove(fileName)
		if err != nil {
			panic("删除文件失败!")
		}
	}()

	// 上传minio
	info := minioutil.UploadMinio(fileName, fileName)

	return info
}

/**
 * 解析Excel文件
 */
func ParseExcel(fileName string, vo interface{}, sb *strings.Builder, faultCountMap *map[int]string) []interface{} {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		panic(err)
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	// 获得对象的所有字段和字典tag
	fields, tags, _ := base.GetVoFieldTypesValues(vo, "excel", "dict")

	// 获得所有行
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return []interface{}{}
	}

	// 返回结果
	var res []interface{}

	// 每一列对应vo的对象
	var colMapperVoField map[int]string = make(map[int]string)

	var headers []string
	// 每一行
	for i, row := range rows {
		// 创建一个vo对象
		voValue := reflect.New(reflect.TypeOf(vo))
		voPtr := voValue.Interface()

		// 每一列
		for j := 0; j < len(row); j++ {
			// 第一行是头部，需要额外处理,填充colMapperVoField
			if i == 0 {
				// 表头的内容
				header := strings.TrimSpace(row[j])
				for index, tag := range tags {
					txt := tag["excel"]
					if txt == header {
						colMapperVoField[j] = fields[index]
						break
					}
				}
				headers = append(headers, header)
			} else {
				// 其它列
				fieldName := colMapperVoField[j]
				if fieldName == "" {
					continue
				}

				content := strings.TrimSpace(row[j])
				if err := base.SetStructField(voPtr, fieldName, content); err != nil {
					m := fmt.Sprintf("第%d行:【%s】格式错误.", i+1, headers[j])
					sb.WriteString(m)
					sb.WriteString("\n")

					if faultCountMap != nil {
						(*faultCountMap)[i+1] = ""
					}
				}
			}
		}

		if i != 0 {
			// 将vo加到数组
			voActual := voValue.Elem().Interface()
			res = append(res, voActual)
		}
	}

	return res
}
