package elestio

import (
	"encoding/json"
	"strconv"
)

// A FlexString is an string that can be unmarshalled from a JSON field
// that has either a number or a string value.
// E.g. if the json field contains an number 42, the
// FlexString value will be "42".
type FlexString string

func (fi *FlexString) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		return json.Unmarshal(b, (*string)(fi))
	}

	var i int
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	*fi = FlexString(strconv.Itoa(i))
	return nil
}
