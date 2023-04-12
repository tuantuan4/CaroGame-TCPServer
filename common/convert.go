package common

import (
	"fmt"
	"strconv"
)

func StringToInt(number string) int {
	marks, err := strconv.Atoi(number)

	if err != nil {
		fmt.Println("Error during conversion")
		return -1
	}
	return marks
}

func IntToString(number int) string {
	return strconv.Itoa(number)
}
