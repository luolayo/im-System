package util

import (
	"bufio"
	Logger "im-System/logger"
	"os"
	"strings"
)

func InputString() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		logger := Logger.NewLogger(Logger.ErrorLevel)
		logger.Error("Error reading input: %s", err)
	}
	input = strings.TrimSpace(input)
	return input
}
