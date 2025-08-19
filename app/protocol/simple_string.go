package protocol

type simpleStringRESP string

func (s simpleStringRESP) String() string {
	return string(s)
}

func (s simpleStringRESP) IsArray() (string, bool) {
	return "", false
}

func (s simpleStringRESP) IsMap() bool {
	return false
}
