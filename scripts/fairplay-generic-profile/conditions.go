package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var genericSelect = regexp.MustCompile(`^(csel|csinc|csinv|csneg|cset|csetm|cinc|cinv|cneg) ([wx](?:[0-9]+|zr)), (.+), ([a-z]+)$`)
var genericSet = regexp.MustCompile(`^(cset|csetm) ([wx](?:[0-9]+|zr)), ([a-z]+)$`)
var genericCompare = regexp.MustCompile(`^(ccmp|ccmn) ([wx](?:[0-9]+|zr)), ([wx](?:[0-9]+|zr)|#(?:0x[0-9a-f]+|[0-9]+)), #(0x[0-9a-f]+|[0-9]+), ([a-z]+)$`)

func conditionCode(name string) (byte, error) {
	names := map[string]byte{"eq": 0, "ne": 1, "hs": 2, "cs": 2, "lo": 3, "cc": 3, "mi": 4, "pl": 5, "vs": 6, "vc": 7, "hi": 8, "ls": 9, "ge": 10, "lt": 11, "gt": 12, "le": 13, "al": 14, "nv": 15}
	code, ok := names[name]
	if !ok {
		return 0, fmt.Errorf("unknown condition")
	}
	return code, nil
}
func genericRegister(text string, width byte) (byte, error) {
	prefix := byte('w')
	if width == 64 {
		prefix = 'x'
	}
	if len(text) < 2 || text[0] != prefix {
		return 0, fmt.Errorf("register width mismatch")
	}
	if text[1:] == "zr" {
		return 31, nil
	}
	value, err := strconv.ParseUint(text[1:], 10, 5)
	return byte(value), err
}
func genericConditionalRecord(text string) (record [recordBytes]byte, found bool, err error) {
	if m := genericCompare.FindStringSubmatch(text); m != nil {
		record[0] = 4
		if m[1] == "ccmn" {
			record[1] = 1
		}
		record[2] = 32
		if m[2][0] == 'x' {
			record[2] = 64
		}
		record[3], err = genericRegister(m[2], record[2])
		if err != nil {
			return record, true, err
		}
		if strings.HasPrefix(m[3], "#") {
			record[5] = 1
			value, e := strconv.ParseUint(m[3][1:], 0, 5)
			if e != nil {
				return record, true, e
			}
			record[4] = byte(value)
		} else {
			record[4], err = genericRegister(m[3], record[2])
			if err != nil {
				return record, true, err
			}
		}
		nzcv, e := strconv.ParseUint(m[4], 0, 4)
		if e != nil {
			return record, true, e
		}
		record[7] = byte(nzcv)
		record[6], err = conditionCode(m[5])
		return record, true, err
	}
	set := genericSet.FindStringSubmatch(text)
	match := genericSelect.FindStringSubmatch(text)
	if set == nil && match == nil {
		return record, false, nil
	}
	var operation, destination, operands, condition string
	if set != nil {
		operation, destination, condition = set[1], set[2], set[3]
	} else {
		operation, destination, operands, condition = match[1], match[2], match[3], match[4]
	}
	record[0] = 3
	record[2] = 32
	if destination[0] == 'x' {
		record[2] = 64
	}
	record[3], err = genericRegister(destination, record[2])
	if err != nil {
		return record, true, err
	}
	record[6], err = conditionCode(condition)
	if err != nil {
		return record, true, err
	}
	switch operation {
	case "cset", "csetm":
		record[1] = 1
		if operation == "csetm" {
			record[1] = 2
		}
		record[4], record[5] = 31, 31
		record[6] ^= 1
	case "cinc", "cinv", "cneg":
		record[1] = map[string]byte{"cinc": 1, "cinv": 2, "cneg": 3}[operation]
		record[4], err = genericRegister(operands, record[2])
		if err != nil {
			return record, true, err
		}
		record[5] = record[4]
		record[6] ^= 1
	default:
		record[1] = map[string]byte{"csel": 0, "csinc": 1, "csinv": 2, "csneg": 3}[operation]
		parts := strings.Split(operands, ", ")
		if len(parts) != 2 {
			return record, true, fmt.Errorf("conditional operands invalid")
		}
		record[4], err = genericRegister(parts[0], record[2])
		if err != nil {
			return record, true, err
		}
		record[5], err = genericRegister(parts[1], record[2])
		if err != nil {
			return record, true, err
		}
	}
	return record, true, nil
}

var genericVectorLogic = regexp.MustCompile(`^(and|orr|eor|bic|orn|eon|bsl|bit|bif) v([0-9]+)\.(8|16)b, v([0-9]+)\.(8|16)b, v([0-9]+)\.(8|16)b$`)

func genericVectorRecord(text string) (record [recordBytes]byte, found bool, err error) {
	m := genericVectorLogic.FindStringSubmatch(text)
	if m == nil {
		return record, false, nil
	}
	if m[3] != m[5] || m[3] != m[7] {
		return record, true, fmt.Errorf("vector shape mismatch")
	}
	record[0] = 5
	record[1] = map[string]byte{"and": 0, "orr": 1, "eor": 2, "bic": 3, "orn": 4, "eon": 5, "bsl": 6, "bit": 7, "bif": 8}[m[1]]
	for i, index := range []int{2, 4, 6} {
		value, e := strconv.ParseUint(m[index], 10, 5)
		if e != nil {
			return record, true, e
		}
		record[3+i] = byte(value)
	}
	width, e := strconv.ParseUint(m[3], 10, 8)
	record[2] = byte(width)
	return record, true, e
}
