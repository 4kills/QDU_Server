package main

import (
	"errors"
	"math"
	"strings"
)

func decodeBase64(s string) ([]byte, error) {
	bits, err := base64ToBits(s)
	if err != nil {
		return nil, err
	}

	return bitsToBytes(bits)
}

func bitsToBytes(bits []bool) ([]byte, error) {
	if len(bits)%8 != 0 {
		return nil, errors.New("base64Decode error: length of string invalid, byte array not a multiple of 8")
	}

	var b []byte
	for i := 0; i < 128; i += 8 {
		b = append(b, byte(bitsToDez(bits[i:i+8])))
	}

	return b, nil
}

func base64ToBits(s string) ([]bool, error) {
	var bits []bool
	for i, val := range s {
		num, err := base64Decoding(string(val))
		if err != nil {
			return nil, err
		}

		if i == 0 {
			bits = append(bits, numToBits(num, 2)...)
			continue
		}
		bits = append(bits, numToBits(num, 6)...)
	}
	return bits, nil
}

func numToBits(n, s int) []bool {
	bits := byteToBits(byte(n))
	return bits[8-s : 8]
}

func base64Decoding(s string) (int, error) {
	const codeSet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	num := 0
	n := float64(len(s)) - 1
	for i := 0; i < len(s); i++ {
		index := strings.IndexByte(codeSet, s[i])
		if index == -1 {
			return num, errors.New(`base64Decode error: semantic error: 
			string was invalid, char not found in codeset`)
		}
		num += index * int(math.Pow(64, n))
	}
	return num, nil
}
