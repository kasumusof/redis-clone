package protocol

import "strconv"

type integerRESP int

func (i integerRESP) String() string {
	return strconv.Itoa(int(i))
}

func (i integerRESP) IsArray() (string, bool) {
	return "", false
}

func (i integerRESP) IsMap() bool {
	return false
}
