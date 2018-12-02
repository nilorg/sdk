package store

import (
	"database/sql"
	"reflect"
)

// ScanRowsToMaps 扫描行
func ScanRowsToMaps(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// 临时缓存，用于保存从rows中读取出来的一行。
	buff := make([]interface{}, len(cols))
	for i := range cols {
		// 初始化key的value
		var value interface{}
		buff[i] = &value
	}
	var data []map[string]interface{}
	for rows.Next() {
		// 给临时缓存赋值
		if err := rows.Scan(buff...); err != nil {
			return nil, err
		}
		// 行
		line := make(map[string]interface{}, len(cols))
		for i, v := range cols {
			if buff[i] == nil {
				continue
			}
			// 反射
			value := reflect.Indirect(reflect.ValueOf(buff[i]))
			line[v] = value.Interface()
		}
		data = append(data, line)
	}
	return data, nil
}
