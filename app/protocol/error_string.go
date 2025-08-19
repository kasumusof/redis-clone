package protocol

type errorStringRESP string

func (e errorStringRESP) String() string {
	return string(e)
}

func (e errorStringRESP) IsArray() (string, bool) {
	return "", false
}

func (e errorStringRESP) IsMap() bool {
	return false
}
