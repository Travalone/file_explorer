package utils

import (
	"sort"
	"strings"
)

func StringsSort(list []string) {
	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i]) < strings.ToLower(list[j])
	})
}

func CheckListsEqual(list1 []string, list2 []string) bool {
	if len(list1) != len(list2) {
		return false
	}
	for i, item := range list1 {
		if item != list2[i] {
			return false
		}
	}
	return true
}

func MergeLists(list1 []string, list2 []string) []string {
	set := map[string]struct{}{}
	for _, item := range list1 {
		set[item] = struct{}{}
	}
	for _, item := range list2 {
		set[item] = struct{}{}
	}

	res := make([]string, 0, len(set))
	for key := range set {
		res = append(res, key)
	}
	return res
}

func ToInterfaces(list []string) []interface{} {
	res := make([]interface{}, len(list))
	for i, item := range list {
		res[i] = item
	}
	return res
}

// MoveListItems 原地批量修改list内指定元素的位置
// 不检查checkIndexList内index合法和重复
func MoveListItems(srcList []interface{}, checkIndexList []int, offset int) {
	if len(srcList) == 0 || len(checkIndexList) == 0 || offset == 0 || len(srcList) == len(checkIndexList) {
		return
	}

	checkIndexSet := make(map[int]struct{})
	for _, checkIndex := range checkIndexList {
		checkIndexSet[checkIndex] = struct{}{}
	}

	// 获取指定元素新位置
	sort.Slice(checkIndexList, func(i, j int) bool {
		if offset < 0 {
			return checkIndexList[i] < checkIndexList[j]
		}
		return checkIndexList[i] > checkIndexList[j]
	})
	dstList := make([]interface{}, len(srcList))
	top, bottom := 0, len(srcList)-1
	for _, oldIndex := range checkIndexList {
		newIndex := oldIndex + offset
		if newIndex < top {
			newIndex = top
			top++
		} else if newIndex > bottom {
			newIndex = bottom
			bottom--
		}
		dstList[newIndex] = srcList[oldIndex]
	}

	// 未指定元素依次占位
	newIndex := 0
	for oldIndex := range srcList {
		// 跳过被选中元素
		if _, ok := checkIndexSet[oldIndex]; ok {
			continue
		}
		// 查找未填充位置
		for dstList[newIndex] != nil {
			newIndex++
		}
		dstList[newIndex] = srcList[oldIndex]
	}

	for i, item := range dstList {
		srcList[i] = item
	}
}

func ListContainsItem(srcList []string, item string) bool {
	for _, srcItem := range srcList {
		if item == srcItem {
			return true
		}
	}
	return false
}
