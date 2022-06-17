package slice

import (
	"fmt"
	"testing"
)

func TestProduct(t *testing.T) {
	attrs := [][]interface{}{
		{
			"女款",
			"男款",
		},
		{
			"卫衣",
			"短袖",
		},
		{
			"S码",
			"L码",
		},
		{
			"蓝色",
			"红色",
		},
	}
	result := Product(attrs...)
	for _, vs := range result {
		for _, v := range vs {
			fmt.Printf("%v ", v)
		}
		fmt.Println("\n===================")
	}
}
