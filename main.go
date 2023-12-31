package main

import (
	"TarsusCloudProxyServer/src/monitor"
	"TarsusCloudProxyServer/src/request"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	monitor.InitDB()
	defer monitor.DB.Close()

	http.HandleFunc("/proxy/", handler)
	http.HandleFunc("/performance", monitor.PerformanceAnalyse)
	http.HandleFunc("/apicalls", monitor.ApiCalls)
	http.ListenAndServe(":3402", nil)
}

var limiter = rate.NewLimiter(100, 200)

func handler(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}
	targetPort := r.Header.Get("x-target-port")
	if targetPort == "" {
		http.Error(w, "Missing target port", http.StatusBadRequest)
		return
	}
	targetPath := strings.Replace(r.URL.String(), "/proxy/", "/", 1)
	port, err := strconv.Atoi(targetPort)
	if err != nil {
		http.Error(w, "target port is not a number", 400)
		return
	}

	url := fmt.Sprintf("http://127.0.0.1:%s%s", targetPort, targetPath)
	body, _ := io.ReadAll(r.Body)
	log, respBody, secs, startTime, err := request.Fetch(url, body)
	requestBody := string(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("request body is %s", requestBody)
	_, err = w.Write(respBody)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
	fmt.Println(log)

	monitor.LoggerPerformance(port, startTime, secs, requestBody, targetPath, len(respBody), 200)

}
