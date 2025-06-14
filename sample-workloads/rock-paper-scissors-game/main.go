package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Game struct {
	ID          int       `json:"id"`
	PlayerName  string    `json:"player_name"`
	PlayerMove  string    `json:"player_move"`
	ComputerMove string   `json:"computer_move"`
	Result      string    `json:"result"`
	Timestamp   time.Time `json:"timestamp"`
	PlayerIP    string    `json:"player_ip"`
}

type Player struct {
	Name       string    `json:"name"`
	Wins       int       `json:"wins"`
	Losses     int       `json:"losses"`
	Draws      int       `json:"draws"`
	Total      int       `json:"total"`
	LastActive time.Time `json:"last_active"`
}

type GameStats struct {
	TotalGames    int                `json:"total_games"`
	TotalPlayers  int                `json:"total_players"`
	RecentGames   []Game             `json:"recent_games"`
	Leaderboard   []Player           `json:"leaderboard"`
	MoveStats     map[string]int     `json:"move_stats"`
	WinStats      map[string]int     `json:"win_stats"`
}

type GameStatsCache struct {
	TotalGames   int
	MoveStats    map[string]int
	WinStats     map[string]int
	LastUpdated  time.Time
	mu           sync.RWMutex
}

const (
	maxGamesHistory           = 1000         // Keep last 1000 games to prevent memory growth
	maxPlayers               = 500          // Maximum number of players to keep in memory
	playerInactivityDuration = 24 * time.Hour // Remove players inactive for 24 hours
)

var (
	games     []Game
	players   map[string]*Player
	gameID    int
	mu        sync.RWMutex
	moves     = []string{"rock", "paper", "scissors"}
	statsCache *GameStatsCache
	homeTemplate *template.Template // Add global template variable
)


func init() {
	players = make(map[string]*Player)
	statsCache = &GameStatsCache{
		TotalGames: 0,
		MoveStats: map[string]int{
			"rock":     0,
			"paper":    0,
			"scissors": 0,
		},
		WinStats: map[string]int{
			"win":  0,
			"loss": 0,
			"draw": 0,
		},
		LastUpdated: time.Now(),
	}
	
	// Parse template once at startup with custom functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"winRate": func(p Player) float64 {
			if p.Total == 0 { return 0 }
			return float64(p.Wins) / float64(p.Total) * 100
		},
		"upper": func(s string) string {
			if s == "win" { return "WIN" }
			if s == "loss" { return "LOSS" }
			return "DRAW"
		},
	}
	
	var err error
	homeTemplate, err = template.New("index.html").Funcs(funcMap).ParseFiles("./web/templates/index.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/api/stats", statsHandler)
	http.HandleFunc("/health", healthHandler)

	// Middleware
	var handler http.Handler = http.DefaultServeMux
	handler = requestLogger(handler)
	handler = securityHeaders(handler)
	handler = gzipHandler(handler)

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		log.Printf("🎮 Rock Paper Scissors Arena starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Println("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	stats := getGameStats()
	
	w.Header().Set("Content-Type", "text/html")
	err := homeTemplate.Execute(w, stats)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request struct {
		PlayerName string `json:"player_name"`
		PlayerMove string `json:"player_move"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if request.PlayerName == "" || request.PlayerMove == "" {
		http.Error(w, "Player name and move are required", http.StatusBadRequest)
		return
	}
	
	// Validate move
	if !slices.Contains(moves, request.PlayerMove) {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}
	
	// Generate computer move
	computerMove := moves[rand.Intn(len(moves))]
	
	// Determine result
	result := determineWinner(request.PlayerMove, computerMove)
	
	// Get client IP
	clientIP := getClientIP(r)
	
	// Create game record
	mu.Lock()
	gameID++
	game := Game{
		ID:           gameID,
		PlayerName:   request.PlayerName,
		PlayerMove:   request.PlayerMove,
		ComputerMove: computerMove,
		Result:       result,
		Timestamp:    time.Now(),
		PlayerIP:     clientIP,
	}
	games = append(games, game)
	
	// Cap games history to prevent indefinite memory growth
	if len(games) > maxGamesHistory {
		// Keep only the last maxGamesHistory games
		games = games[len(games)-maxGamesHistory:]
	}
	
	// Update player stats
	player, exists := players[request.PlayerName]
	if !exists {
		player = &Player{Name: request.PlayerName}
		players[request.PlayerName] = player
	}
	
	player.Total++
	player.LastActive = time.Now()
	switch result {
	case "win":
		player.Wins++
	case "loss":
		player.Losses++
	case "draw":
		player.Draws++
	}
	
	// Update stats cache incrementally
	updateStatsCache(game)
	
	// Evict inactive players if we have too many
	evictInactivePlayers()
	mu.Unlock()
	
	// Return game result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	stats := getGameStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func determineWinner(playerMove, computerMove string) string {
	if playerMove == computerMove {
		return "draw"
	}
	
	winConditions := map[string]string{
		"rock":     "scissors",
		"paper":    "rock",
		"scissors": "paper",
	}
	
	if winConditions[playerMove] == computerMove {
		return "win"
	}
	return "loss"
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for load balancers/proxies)
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For can contain a comma-separated list of IPs
		parts := strings.Split(xForwardedFor, ",")
		if len(parts) > 0 {
			clientIP := strings.TrimSpace(parts[0])
			if clientIP != "" {
				return clientIP
			}
		}
	}
	
	// Check X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}
	
	// Fall back to RemoteAddr
	return r.RemoteAddr
}

func getGameStats() GameStats {
	mu.RLock()
	defer mu.RUnlock()
	
	// Use cached stats for performance
	statsCache.mu.RLock()
	moveStats := make(map[string]int)
	for k, v := range statsCache.MoveStats {
		moveStats[k] = v
	}
	winStats := make(map[string]int)
	for k, v := range statsCache.WinStats {
		winStats[k] = v
	}
	totalGames := statsCache.TotalGames
	statsCache.mu.RUnlock()
	
	// Get recent games (last 10)
	recentGames := make([]Game, 0)
	start := max(0, len(games)-10)
	for i := len(games) - 1; i >= start; i-- {
		recentGames = append(recentGames, games[i])
	}
	
	// Create leaderboard
	leaderboard := make([]Player, 0, len(players))
	for _, player := range players {
		leaderboard = append(leaderboard, *player)
	}
	
	// Sort by win rate, then by total games
	sort.Slice(leaderboard, func(i, j int) bool {
		var winRateI, winRateJ float64
		if leaderboard[i].Total > 0 {
			winRateI = float64(leaderboard[i].Wins) / float64(leaderboard[i].Total)
		}
		if leaderboard[j].Total > 0 {
			winRateJ = float64(leaderboard[j].Wins) / float64(leaderboard[j].Total)
		}
		
		if winRateI == winRateJ {
			return leaderboard[i].Total > leaderboard[j].Total
		}
		return winRateI > winRateJ
	})
	
	// Limit leaderboard to top 10
	if len(leaderboard) > 10 {
		leaderboard = leaderboard[:10]
	}
	
	return GameStats{
		TotalGames:   totalGames,
		TotalPlayers: len(players),
		RecentGames:  recentGames,
		Leaderboard:  leaderboard,
		MoveStats:    moveStats,
		WinStats:     winStats,
	}
}

// updateStatsCache incrementally updates the stats cache when a game is played
// Note: This function assumes the caller holds the write lock (mu.Lock())
func updateStatsCache(game Game) {
	statsCache.mu.Lock()
	defer statsCache.mu.Unlock()
	
	statsCache.TotalGames++
	statsCache.MoveStats[game.PlayerMove]++
	statsCache.WinStats[game.Result]++
	statsCache.LastUpdated = time.Now()
}

// evictInactivePlayers removes inactive players to prevent memory growth
// Note: This function assumes the caller holds the write lock (mu.Lock())
func evictInactivePlayers() {
	if len(players) <= maxPlayers {
		return
	}
	
	now := time.Now()
	cutoff := now.Add(-playerInactivityDuration)
	
	// First, remove players inactive beyond the cutoff
	for name, player := range players {
		if player.LastActive.Before(cutoff) {
			delete(players, name)
		}
	}
	
	// If still too many players, remove the least recently active ones
	if len(players) > maxPlayers {
		type playerActivity struct {
			name       string
			lastActive time.Time
		}
		
		// Collect all players with their last activity
		activities := make([]playerActivity, 0, len(players))
		for name, player := range players {
			activities = append(activities, playerActivity{
				name:       name,
				lastActive: player.LastActive,
			})
		}
		
		// Sort by last activity (oldest first)
		sort.Slice(activities, func(i, j int) bool {
			return activities[i].lastActive.Before(activities[j].lastActive)
		})
		
		// Remove the oldest players until we're under the limit
		playersToRemove := len(players) - maxPlayers
		for i := 0; i < playersToRemove; i++ {
			delete(players, activities[i].name)
		}
	}
}

// Middleware

// requestLogger logs incoming requests
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// securityHeaders adds security-related headers to the response
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "no-referrer-when-downgrade")
		w.Header().Set("Permissions-Policy", "geolocation=(self), microphone=(), camera=()")
		
		// Content Security Policy (CSP) - tightened for security
		csp := "default-src 'self'; img-src 'self' data:; script-src 'self'; style-src 'self';"
		w.Header().Set("Content-Security-Policy", csp)
		
		next.ServeHTTP(w, r)
	})
}

// gzipHandler compresses response bodies with gzip
func gzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Bypass compression if the client doesn't support it
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		
		// Set headers before writing starts
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		
		// Create gzip writer that writes directly to response
		gzw := gzip.NewWriter(w)
		defer gzw.Close()
		
		// Wrap response writer with gzip writer
		gzipResponseWriter := &gzipResponseWriter{
			ResponseWriter: w,
			gzipWriter:     gzw,
		}
		
		// Serve the request with wrapped writer
		next.ServeHTTP(gzipResponseWriter, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.gzipWriter.Write(b)
}

func (grw *gzipResponseWriter) WriteHeader(statusCode int) {
	grw.ResponseWriter.WriteHeader(statusCode)
}

