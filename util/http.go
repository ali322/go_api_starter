package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Reply(w http.ResponseWriter, v any) {
	jsonStr, err := json.Marshal(map[string]any{
		"code": 0, "data": v,
	})
	if err != nil {
		Reject(w, -1, err)
	}
	fmt.Fprint(w, string(jsonStr))
}

func Reject(w http.ResponseWriter, code int, v any) {
	jsonStr, _ := json.Marshal(map[string]any{
		"code": code, "data": v,
	})
	fmt.Fprint(w, string(jsonStr))
}
