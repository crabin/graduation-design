/*
 * @author Crabin
 */

package utils

import (
	"encoding/json"
)

// 将结构体转换为字节切片
func StructToBytes(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// 将字节切片转换为结构体
func BytesToStruct(b []byte, v interface{}) error {
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}
