package monitor

import (
	"encoding/json"
	"net/http"
)

var (
	id                 int
	port               int
	requestStartTime   string
	responseTimeMs     float64
	requestParameters  string
	responseBodyLength int
	responseStatusCode int
	requestUrl         string
	otherField         string
)

func PerformanceAnalyse(w http.ResponseWriter, r *http.Request) {
	port := r.URL.Query().Get("port")

	rows, err := DB.Query("SELECT * from performance_monitoring WHERE PORT = ?", port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(&id, &port, &requestStartTime, &responseTimeMs, &requestParameters, &responseBodyLength, &responseStatusCode, &requestUrl, &otherField)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		item := map[string]interface{}{
			"id":                   id,
			"port":                 port,
			"request_start_time":   requestStartTime,
			"response_time_ms":     responseTimeMs,
			"request_parameters":   requestParameters,
			"response_body_length": responseBodyLength,
			"response_status_code": responseStatusCode,
			"request_url":          requestUrl,
			"other_field":          otherField,
		}
		data = append(data, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
