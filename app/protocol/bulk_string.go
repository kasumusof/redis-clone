package protocol

type bulkStringRESP string

func (b bulkStringRESP) IsArray() (string, bool) {
	return "", false
}

func (b bulkStringRESP) IsMap() bool {
	return false
}

func (b bulkStringRESP) String() string {
	return string(b)
}
