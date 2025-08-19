package protocol

import "strings"

type arrayRESP []RESP

func (a arrayRESP) IsArray() (string, bool) {
	return "|", true
}

func (a arrayRESP) IsMap() bool {
	return false
}

func (a arrayRESP) String() string {
	var resp []string
	for _, v := range a {
		resp = append(resp, v.String())
	}
	return strings.Join(resp, "|")
}
