package common

import (
	"bufio"
	"os"
)

// ReadLines 文件按行读取
func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 跳过空行
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}
