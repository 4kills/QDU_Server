package main

import (
	"math"
	"strings"
)

func encodeBase64(u []byte) string {
	return bitsToBase64(uuidToBits(u))
}

func bitsToBase64(bits []bool) string {
	var sS []string
	sS = append(sS, base64Encoding(bitsToDez(bits[:2])))
	for i := 0; i < 21; i++ {
		sS = append(sS, base64Encoding(bitsToDez(bits[2+6*i:8+6*i])))
	}
	return strings.Join(sS, "")
}

func base64Encoding(num int) string {
	const codeSet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	base := int(len(codeSet)) // 64
	var encoded []string

	if num == 0 {
		return string(codeSet[0])
	}
	for num != 0 {
		remainder := num % base
		encoded = append([]string{codeSet[remainder : remainder+1]}, strings.Join(encoded, ""))
		num = num / base
	}

	return strings.Join(encoded, "")
}

func bitsToDez(bits []bool) int { // 00 0110
	num := 0
	n := float64(len(bits)) - 1
	for _, val := range bits {
		if val {
			num += int(math.Pow(2, n))
		}
		n--
	}
	return num
}

func uuidToBits(u []byte) []bool {
	var bits []bool
	for _, val := range u {
		bitArr := byteToBits(val)
		bits = append(bits, bitArr[:]...)
	}
	return bits
}

func byteToBits(by byte) [8]bool {
	var bits [8]bool
	pos := 7
	for n := 0; n < 8; n++ {
		mask := math.Pow(2, float64(n))
		if (by & byte(mask)) != 0 {
			bits[pos] = true
		}
		pos--
	}
	return bits
}
