package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type IPInfo struct {
	ClientIP    string `json:"client_ip"`
	ServerIP    string `json:"server_ip"`
	Hostname    string `json:"hostname"`
	UserAgent   string `json:"user_agent"`
	Headers     map[string]string `json:"headers"`
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>IP Information Service</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        .ip-info { background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .header { color: #333; border-bottom: 2px solid #007acc; padding-bottom: 10px; }
        .info-row { margin: 10px 0; }
        .label { font-weight: bold; color: #555; }
        .value { font-family: monospace; background: #fff; padding: 5px; border-radius: 3px; }
        .json-output { background: #2d3748; color: #e2e8f0; padding: 15px; border-radius: 5px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="header">🌐 IP Information Service</h1>
        
        <div class="ip-info">
            <div class="info-row">
                <span class="label">Your IP Address:</span> 
                <span class="value">{{.ClientIP}}</span>
            </div>
            <div class="info-row">
                <span class="label">Server IP:</span> 
                <span class="value">{{.ServerIP}}</span>
            </div>
            <div class="info-row">
                <span class="label">Hostname:</span> 
                <span class="value">{{.Hostname}}</span>
            </div>
            <div class="info-row">
                <span class="label">User Agent:</span> 
                <span class="value">{{.UserAgent}}</span>
            </div>
        </div>

        <h3>Request Headers:</h3>
        <div class="ip-info">
            {{range $key, $value := .Headers}}
            <div class="info-row">
                <span class="label">{{$key}}:</span> 
                <span class="value">{{$value}}</span>
            </div>
            {{end}}
        </div>

        <h3>JSON Output:</h3>
        <div class="json-output">
            <pre>{{.}}</pre>
        </div>

        <p><a href="/api">Get JSON API response</a></p>
    </div>
</body>
</html>
`

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for load balancers/proxies)
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func getServerIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func getIPInfo(r *http.Request) IPInfo {
	hostname, _ := os.Hostname()
	
	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = strings.Join(values, ", ")
	}

	return IPInfo{
		ClientIP:  getClientIP(r),
		ServerIP:  getServerIP(),
		Hostname:  hostname,
		UserAgent: r.UserAgent(),
		Headers:   headers,
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo := getIPInfo(r)
	
	tmpl, err := template.New("ip").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, ipInfo)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo := getIPInfo(r)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ipInfo)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}