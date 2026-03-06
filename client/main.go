package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	js.Global.Set("Backend", map[string]interface{}{
		"FetchRows": FetchRows,
	})
}

// FetchRows is called by React AG Grid Datasource
// params: { start: int, end: int, successCallback: func(rows, lastRow) }
func FetchRows(start, end int, successCallback *js.Object) {
	go func() {
		u := fmt.Sprintf("/api/rows?start=%d&end=%d", start, end)
		resp, err := http.Get(u)
		if err != nil {
			println("Error fetching rows:", err.Error())
			return
		}
		defer resp.Body.Close()

		var result struct {
			Rows    []map[string]interface{} `json:"rows"`
			LastRow int                      `json:"lastRow"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			println("Error decoding JSON:", err.Error())
			return
		}

		// Call the JS callback
		successCallback.Invoke(result.Rows, result.LastRow)
	}()
}
