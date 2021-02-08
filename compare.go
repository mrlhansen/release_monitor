package main

import (
	"strings"
	"regexp"
	"strconv"
)

func convert_prerelease(str string) int {
	var x int
	str = strings.ToLower(str)
	switch str {
		case "alpha":
			x = 0
		case "beta":
			x = 1
		case "rc":
			x = 2
		default:
			x = -1
	}
	return x
}

func version_compare(str1, str2 string) int {
	pa := strings.Split(str1, "-")
	pb := strings.Split(str2, "-")
	va := strings.Split(pa[0], ".")
	vb := strings.Split(pb[0], ".")

	// Check version string
	x := len(va)
	y := len(vb)

	for x < y {
		va = append(va, "0")
		x++
	}

	for y < x {
		vb = append(vb, "0")
		y++
	}

	for i := 0; i < x; i++ {
		a,_ := strconv.Atoi(va[i])
		b,_ := strconv.Atoi(vb[i])

		if a < b {
			return -1
		}

		if a > b {
			return +1
		}
	}

	// Equal so far, check pre-release
	x = len(pa)
	y = len(pb)

	if x > y {
		return -1
	}

	if x < y {
		return +1
	}

	if x == 1 {
		return 0
	}

	// Both have a pre-release string
	re := regexp.MustCompile(`([a-zA-Z]+)\.?(\d+)?`)
	ma := re.FindStringSubmatch(pa[1])
	mb := re.FindStringSubmatch(pb[1])

	x = convert_prerelease(ma[1])
	y = convert_prerelease(mb[1])

	if x < y {
		return -1
	}

	if x > y {
		return +1
	}

	// Still equal, check pre-release number
	if len(ma) == 3 {
		x,_ = strconv.Atoi(ma[2])
	} else {
		x = 1
	}

	if len(mb) == 3 {
		y,_ = strconv.Atoi(mb[2])
	} else {
		y = 1
	}

	if x < y {
		return -1
	}

	if x > y {
		return +1
	}

	// They are equal
	return 0
}
