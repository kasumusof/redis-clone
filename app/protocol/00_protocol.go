package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const (
	newline    = '\n'
	returnLine = '\r'
	crlf       = "\r\n"
)
const (
	simpleString = '+'
	errorString  = '-'
	integer      = ':'
	bulkString   = '$'
	array        = '*'
	nulls        = '_'
	booleans     = '#'
	doubles      = ','
	bigNumbers   = '('
	bulkErrors   = '!'
	verbatim     = '='
	maps         = '%'
	attributes   = '|'
	sets         = '~'
	pushes       = '>'
)

type RESP interface {
	String() string
	IsArray() (string, bool)
	IsMap() bool
}

var (
	_ RESP = (*arrayRESP)(nil)
	_ RESP = (*bulkStringRESP)(nil)
	_ RESP = (*simpleStringRESP)(nil)
	_ RESP = (*errorStringRESP)(nil)
	_ RESP = (*integerRESP)(nil)
)

func ParseRequest(buff *bufio.Reader) (RESP, error) {
	prefix, err := buff.ReadByte()
	if err != nil {
		return nil, err
	}

	switch prefix {
	case simpleString:
		arg, err := buff.ReadString(newline)
		if err != nil {
			return nil, err
		}

		arg = strings.TrimSpace(arg)
		return simpleStringRESP(arg), nil
	case errorString:
		arg, err := buff.ReadString(newline)
		if err != nil {
			return nil, err
		}

		arg = strings.TrimSpace(arg)
		return errorStringRESP(arg), nil
	case integer:
		arg, err := buff.ReadString(newline)
		if err != nil {
			return nil, err
		}

		arg = strings.TrimSpace(arg)
		n, err := strconv.Atoi(arg)
		if err != nil {
			return nil, errors.Join(errors.New(" failed to parse integer"), err)
		}

		return integerRESP(n), nil
	case bulkString:
		nStr, err := buff.ReadString(newline)
		if err != nil {
			return nil, err
		}

		nStr = strings.TrimSpace(nStr)
		n, err := strconv.Atoi(nStr)
		if err != nil {
			return nil, errors.Join(errors.New(" failed to parse bulk string size"), err)
		}
		_ = n

		str, err := buff.ReadString(newline)
		if err != nil {
			return nil, err
		}
		str = strings.TrimSpace(str)

		return bulkStringRESP(str), nil
	case array:
		nStr, err := buff.ReadString(newline)
		if err != nil {
			return nil, err
		}

		nStr = strings.TrimSpace(nStr)
		n, err := strconv.Atoi(nStr)
		if err != nil {
			return nil, errors.Join(errors.New(" failed to parse bulk string size"), err)
		}

		if n == -1 {
			return nil, nil
		}

		args := make(arrayRESP, n)
		for i := 0; i < n; i++ {
			args[i], err = ParseRequest(buff)
			if err != nil {
				return nil, err
			}
		}
		return args, nil
	//case nulls:
	//case booleans:
	//case doubles:
	//case bigNumbers:
	//case bulkErrors:
	//case verbatim:
	//case maps:
	//case attributes:
	//case sets:
	//case pushes:
	default:
		return errorStringRESP("ERR not implemented"), nil

	}
}

func SimpleString(s string) string {
	return fmt.Sprintf("%c%s\r\n", simpleString, s)
}

func ErrorString(s string) string {
	return fmt.Sprintf("%c%s\r\n", errorString, s)
}

func Integer(i int) string {
	return fmt.Sprintf("%c%d\r\n", integer, i)
}

func BulkString(s string) string {
	l := len(s)
	if l == 0 {
		l = -1
	}

	resp := fmt.Sprintf("%c%d\r\n%s\r\n", bulkString, l, s)
	return resp
}

func Array(a []any) string {
	i := len(a)
	if i == 0 { // null array
		i = -1
	}

	resp := fmt.Sprintf("%c%d\r\n", array, i)
	for _, v := range a { // TODO: add more types
		switch v.(type) {
		case string:
			resp += BulkString(v.(string))
		case int:
			resp += Integer(v.(int))
		case bool:
			resp += Booleans(v.(bool))
		case error:
			resp += ErrorString(v.(error).Error())
		default:
			resp += Nulls()
		}
	}
	return resp
}

func Nulls() string {
	return fmt.Sprintf("%c\r\n", nulls)
}

func Booleans(val bool) string {
	str := "f"
	if val {
		str = "t"
	}
	return fmt.Sprintf("%c%s\r\n", booleans, str)
}

func Doubles(val float64) string {
	return fmt.Sprintf("%c%f\r\n", doubles, val)
}

func BigNumbers(val big.Int) string {
	return fmt.Sprintf("%c%s\r\n", bigNumbers, val.String())
}

func BulkErrors(val []error) string {
	resp := fmt.Sprintf("%c%d\r\n", bulkErrors, len(val))
	for _, v := range val {
		resp += ErrorString(v.Error())
	}

	return resp
}

func Verbatim(val, encoding string) string {
	defaultEncoding := "txt"
	if encoding == "" {
		encoding = defaultEncoding
	}

	resp := fmt.Sprintf("%c%d\r\n", verbatim, len(val))
	return resp + BulkString(encoding+":"+val)
}

func Maps(val map[any]any) string {
	resp := fmt.Sprintf("%c%d\r\n", maps, len(val))
	for k, v := range val {
		keyEncoded := encodeRandomString(k)
		valueEncoded := encodeRandomString(v)
		resp += fmt.Sprintf("%v:%v\r\n", keyEncoded, valueEncoded)
	}
	return resp
}

func encodeRandomString(val any) string {
	switch val.(type) {
	case string:
		return SimpleString(val.(string))
	case int:
		return Integer(val.(int))
	case bool:
		return Booleans(val.(bool))
	case error:
		return ErrorString(val.(error).Error())
	case float64:
		return Doubles(val.(float64))
	case big.Int:
		return BigNumbers(val.(big.Int))
	case []error:
		return BulkErrors(val.([]error))
	case map[any]any:
		return Maps(val.(map[any]any))
	case []any:
		return Array(val.([]any))
	default:
		return Nulls()
	}
}

func Attributes(val map[string]interface{}) string {
	return fmt.Sprintf("%c%d\r\n", attributes, len(val))
}

func Sets(val map[string]interface{}) string {
	return fmt.Sprintf("%c%d\r\n", sets, len(val))
}

func Pushes(val map[string]interface{}) string {
	return fmt.Sprintf("%c%d\r\n", pushes, len(val))
}
