// Provides BasicType(s) that is/are then handled by NamedType(s).
package basic

import (
	"encoding/json"
	"bytes"
)

// Type of our enum
type BasicType int

// Definition of the BasicType enum
const ( // NOTE: this will get out of hand, FIND ALTERNATIVE!!!
	Num BasicType = iota
	Str
)

var StrToBasicTypes = map[string]BasicType {
	"num": Num,
	"str": Str,
}

var basicTypesToStr = map[BasicType]string {
	Num: "num",
	Str: "str",
}

// Convert a BasicType to it's string counterpart
func (t BasicType) String() (str string) {
	str, ok := basicTypesToStr[t]
	if !ok {
		str = "unknown"
	}
	return
}

// Convert a string to it's BasicType counterpart
func StrToBasicType(str string) (BasicType, bool) {
	typ, ok := StrToBasicTypes[str]
	return typ, ok
}

// Convert a BasicType to JSON
func (typ BasicType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(typ.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// Convert JSON to a BasicType
func (typ *BasicType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	readTyp, ok := StrToBasicType(j)
	if ok {
		*typ = readTyp
	} 	
	return nil
}
