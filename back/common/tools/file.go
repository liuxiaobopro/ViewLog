package tools

import (
	"bufio"
	"os"
	"strings"
)

/**
 * 读取文件
 *
 * @param filepath 文件路径
 * @param m 读取行数
 * @param n 跳过行数
 * @return []string
 * @return error
 */
func ReadFileLine(filepath string, m, n int) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0, m)

	// 跳过前n-1行
	for i := 1; i <= n && scanner.Scan(); i++ {
	}

	for i := 1; i <= m && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	str := strings.Join(lines, "\n")
	return str, nil
}
