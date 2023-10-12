package structutil

import (
	"bytes"
	"encoding/json"
)

func ToString(v any) string {
	bs, _ := json.Marshal(v)
	buf := new(bytes.Buffer)
	_ = json.Indent(buf, bs, "", "    ")
	return buf.String()
}
