package test

import (
	"os"
	"testing"
)

func TestDeleteFile(t *testing.T) {
	path := "/Users/bytedance/Downloads/test.txt"
	os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)

	//err := os.Remove(path)
	//if err != nil {
	//	return
	//}
}
