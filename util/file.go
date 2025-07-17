package util

import (
	"os"
	"strings"
)

// Determine whether the file exists
func FileExists(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		return !fi.IsDir() // If the file exists and is not a directory, return true
	}
	return os.IsNotExist(err) // If the error indicates that the file does not exist, return true
}

// Create folders by the file path
func MKdirWithFilePath(filePath string)error {
	paths := strings.Split(filePath, "/") // 返回一個切片(Slice)
	paths[len(paths)-1] = "" // 將最後一個元素設為空字串，這樣就不會創建最後一個文件

	for i, v := range paths {
		if i == len(paths)-1 {
			break // 如果是最後一個元素，則跳出迴圈
		}
		if i != 0 {
			paths[len(paths)-1] += "/" // 在每個路徑前加上斜線，除了第一個元素
		}
		paths[len(paths)-1] += v // 將當前元素添加到最後一個元素中
	}

	return  os.MkdirAll(paths[len(paths)-1], 0775) // 使用 os.MkdirAll 創建所有必要的目錄，並設置權限為 0775
}