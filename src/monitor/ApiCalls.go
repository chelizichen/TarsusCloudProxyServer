package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiData struct {
	request_day string
	api_calls   int32
}

func ApiCalls(w http.ResponseWriter, r *http.Request) {
	port := r.URL.Query().Get("port")
	fmt.Println("getPortApiCallsChart  %s", port)
	rows, err := DB.Query("\nSELECT \n    DATE(request_start_time) AS request_day,\n    COUNT(*) AS api_calls\nFROM performance_monitoring\nWHERE request_start_time >= NOW() - INTERVAL 7 DAY \nAND performance_monitoring.`port` = ?\nGROUP BY DATE(request_start_time)\nORDER BY request_day;\n", port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	var requestDay string
	var apiCalls int
	for rows.Next() {
		err := rows.Scan(&requestDay, &apiCalls)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		item := map[string]interface{}{
			"request_day": requestDay,
			"api_calls":   apiCalls,
		}
		data = append(data, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
