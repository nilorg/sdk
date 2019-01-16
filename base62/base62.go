package base62

import (
	"math"
	"strings"
)

const (
	code62    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	codeLenth = 62
)

var edoc = map[string]int{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16, "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23, "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30, "v": 31, "w": 32, "x": 33, "y": 34, "z": 35, "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42, "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49, "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56, "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61}

// Encode 编码 整数 为 base62 字符串
func Encode(number int) string {
	if number == 0 {
		return "0"
	}
	result := make([]byte, 0)
	for number > 0 {
		round := number / codeLenth
		remain := number % codeLenth
		result = append(result, code62[remain])
		number = round
	}
	return string(result)
}

// Decode 解码字符串为整数
func Decode(str string) int {
	str = strings.TrimSpace(str)
	result := 0
	for index, char := range []byte(str) {
		result += edoc[string(char)] * int(math.Pow(codeLenth, float64(index)))
	}
	return result
}
