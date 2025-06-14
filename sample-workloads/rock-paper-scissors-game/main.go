package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"
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
	Name   string `json:"name"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Draws  int    `json:"draws"`
	Total  int    `json:"total"`
}

type GameStats struct {
	TotalGames    int                `json:"total_games"`
	TotalPlayers  int                `json:"total_players"`
	RecentGames   []Game             `json:"recent_games"`
	Leaderboard   []Player           `json:"leaderboard"`
	MoveStats     map[string]int     `json:"move_stats"`
	WinStats      map[string]int     `json:"win_stats"`
}

var (
	games    []Game
	players  map[string]*Player
	gameID   int
	mu       sync.RWMutex
	moves    = []string{"rock", "paper", "scissors"}
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>🎮 Rock Paper Scissors Arena</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            color: #333;
        }
        .container { 
            max-width: 1200px; 
            margin: 0 auto; 
            padding: 20px;
        }
        .header {
            text-align: center;
            color: white;
            margin-bottom: 30px;
        }
        .header h1 {
            font-size: 3em;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }
        .game-section {
            background: white;
            border-radius: 15px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.3);
        }
        .move-buttons {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin: 20px 0;
        }
        .move-btn {
            font-size: 4em;
            width: 120px;
            height: 120px;
            border: none;
            border-radius: 50%;
            cursor: pointer;
            transition: all 0.3s ease;
            display: flex;
            align-items: center;
            justify-content: center;
            box-shadow: 0 5px 15px rgba(0,0,0,0.2);
        }
        .move-btn:hover {
            transform: scale(1.1);
            box-shadow: 0 8px 25px rgba(0,0,0,0.3);
        }
        .rock { background: linear-gradient(145deg, #ff6b6b, #ee5a5a); }
        .paper { background: linear-gradient(145deg, #4ecdc4, #45b7b8); }
        .scissors { background: linear-gradient(145deg, #ffe66d, #ffcc02); }
        
        .player-input {
            margin: 20px 0;
            text-align: center;
        }
        .player-input input {
            padding: 12px 20px;
            font-size: 16px;
            border: 2px solid #ddd;
            border-radius: 25px;
            width: 250px;
            text-align: center;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }
        .stat-card {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 10px;
            border-left: 5px solid #667eea;
        }
        .stat-value {
            font-size: 2em;
            font-weight: bold;
            color: #667eea;
        }
        
        .leaderboard, .recent-games {
            background: white;
            border-radius: 15px;
            padding: 20px;
            margin-bottom: 20px;
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        .leaderboard h3, .recent-games h3 {
            color: #667eea;
            margin-bottom: 15px;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        .leaderboard-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px;
            margin: 5px 0;
            background: #f8f9fa;
            border-radius: 8px;
        }
        .rank {
            font-weight: bold;
            color: #667eea;
            width: 30px;
        }
        .player-stats {
            flex: 1;
            margin-left: 15px;
        }
        .win-rate {
            font-weight: bold;
            color: #28a745;
        }
        
        .game-item {
            background: #f8f9fa;
            padding: 15px;
            margin: 10px 0;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        .game-result {
            font-weight: bold;
            padding: 5px 10px;
            border-radius: 15px;
            color: white;
        }
        .win { background: #28a745; }
        .loss { background: #dc3545; }
        .draw { background: #ffc107; color: #333; }
        
        .result-display {
            text-align: center;
            margin: 20px 0;
            padding: 20px;
            border-radius: 10px;
            font-size: 1.2em;
        }
        
        @media (max-width: 768px) {
            .move-buttons { flex-direction: column; align-items: center; }
            .move-btn { width: 100px; height: 100px; font-size: 3em; }
            .stats-grid { grid-template-columns: 1fr; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎮 Rock Paper Scissors Arena</h1>
            <p>Challenge the computer and climb the leaderboard!</p>
        </div>
        
        <div class="game-section">
            <h2 style="text-align: center; margin-bottom: 20px;">Play Game</h2>
            <div class="player-input">
                <input type="text" id="playerName" placeholder="Enter your name" maxlength="20">
            </div>
            <div class="move-buttons">
                <button class="move-btn rock" onclick="playGame('rock')">🪨</button>
                <button class="move-btn paper" onclick="playGame('paper')">📄</button>
                <button class="move-btn scissors" onclick="playGame('scissors')">✂️</button>
            </div>
            <div id="gameResult" class="result-display" style="display: none;"></div>
        </div>
        
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-value">{{.TotalGames}}</div>
                <div>Total Games Played</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{.TotalPlayers}}</div>
                <div>Active Players</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{index .MoveStats "rock"}}</div>
                <div>🪨 Rock Played</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{index .MoveStats "paper"}}</div>
                <div>📄 Paper Played</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{index .MoveStats "scissors"}}</div>
                <div>✂️ Scissors Played</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{index .WinStats "win"}}</div>
                <div>🏆 Player Wins</div>
            </div>
        </div>
        
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px;">
            <div class="leaderboard">
                <h3>🏆 Leaderboard</h3>
                {{range $index, $player := .Leaderboard}}
                <div class="leaderboard-item">
                    <div class="rank">#{{add $index 1}}</div>
                    <div class="player-stats">
                        <div style="font-weight: bold;">{{$player.Name}}</div>
                        <div style="font-size: 0.9em; color: #666;">
                            {{$player.Total}} games • {{$player.Wins}}W {{$player.Losses}}L {{$player.Draws}}D
                        </div>
                    </div>
                    <div class="win-rate">{{printf "%.1f" (winRate $player)}}%</div>
                </div>
                {{end}}
            </div>
            
            <div class="recent-games">
                <h3>📋 Recent Games</h3>
                {{range .RecentGames}}
                <div class="game-item">
                    <div style="display: flex; justify-content: space-between; align-items: center;">
                        <div>
                            <strong>{{.PlayerName}}</strong>: {{.PlayerMove}} vs {{.ComputerMove}}
                        </div>
                        <span class="game-result {{.Result}}">{{upper .Result}}</span>
                    </div>
                    <div style="font-size: 0.8em; color: #666; margin-top: 5px;">
                        {{.Timestamp.Format "15:04:05"}} • {{.PlayerIP}}
                    </div>
                </div>
                {{end}}
            </div>
        </div>
    </div>
    
    <script>
        let statsUpdateInProgress = false;

        function updateStats() {
            if (statsUpdateInProgress) return;
            statsUpdateInProgress = true;
            
            fetch('/api/stats')
                .then(response => response.json())
                .then(stats => {
                    // Update stat cards
                    updateStatCards(stats);
                    
                    // Update leaderboard
                    updateLeaderboard(stats.leaderboard);
                    
                    // Update recent games
                    updateRecentGames(stats.recent_games);
                })
                .catch(error => {
                    console.error('Error updating stats:', error);
                })
                .finally(() => {
                    statsUpdateInProgress = false;
                });
        }

        function updateStatCards(stats) {
            const statValues = document.querySelectorAll('.stat-value');
            if (statValues.length >= 6) {
                statValues[0].textContent = stats.total_games;
                statValues[1].textContent = stats.total_players;
                statValues[2].textContent = stats.move_stats.rock || 0;
                statValues[3].textContent = stats.move_stats.paper || 0;
                statValues[4].textContent = stats.move_stats.scissors || 0;
                statValues[5].textContent = stats.win_stats.win || 0;
            }
        }

        function updateLeaderboard(leaderboard) {
            const leaderboardContainer = document.querySelector('.leaderboard');
            if (!leaderboardContainer) return;
            
            // Keep the header, replace the content
            const header = leaderboardContainer.querySelector('h3');
            leaderboardContainer.innerHTML = '';
            leaderboardContainer.appendChild(header);
            
            leaderboard.forEach((player, index) => {
                const winRate = player.total > 0 ? (player.wins / player.total * 100).toFixed(1) : '0.0';
                
                const playerDiv = document.createElement('div');
                playerDiv.className = 'leaderboard-item';
                playerDiv.innerHTML = 
                    '<div class="rank">#' + (index + 1) + '</div>' +
                    '<div class="player-stats">' +
                        '<div style="font-weight: bold;">' + player.name + '</div>' +
                        '<div style="font-size: 0.9em; color: #666;">' +
                            player.total + ' games • ' + player.wins + 'W ' + player.losses + 'L ' + player.draws + 'D' +
                        '</div>' +
                    '</div>' +
                    '<div class="win-rate">' + winRate + '%</div>';
                
                leaderboardContainer.appendChild(playerDiv);
            });
        }

        function updateRecentGames(recentGames) {
            const recentContainer = document.querySelector('.recent-games');
            if (!recentContainer) return;
            
            // Keep the header, replace the content
            const header = recentContainer.querySelector('h3');
            recentContainer.innerHTML = '';
            recentContainer.appendChild(header);
            
            recentGames.forEach(game => {
                const gameDiv = document.createElement('div');
                gameDiv.className = 'game-item';
                
                const timestamp = new Date(game.timestamp).toLocaleTimeString('en-US', {
                    hour12: false,
                    hour: '2-digit',
                    minute: '2-digit',
                    second: '2-digit'
                });
                
                gameDiv.innerHTML = 
                    '<div style="display: flex; justify-content: space-between; align-items: center;">' +
                        '<div>' +
                            '<strong>' + game.player_name + '</strong>: ' + game.player_move + ' vs ' + game.computer_move +
                        '</div>' +
                        '<span class="game-result ' + game.result + '">' + game.result.toUpperCase() + '</span>' +
                    '</div>' +
                    '<div style="font-size: 0.8em; color: #666; margin-top: 5px;">' +
                        timestamp + ' • ' + game.player_ip +
                    '</div>';
                
                recentContainer.appendChild(gameDiv);
            });
        }

        function playGame(playerMove) {
            const playerName = document.getElementById('playerName').value.trim();
            if (!playerName) {
                alert('Please enter your name first!');
                return;
            }
            
            const resultDiv = document.getElementById('gameResult');
            resultDiv.style.display = 'block';
            resultDiv.innerHTML = '<div style="color: #667eea;">🎲 Playing...</div>';
            
            fetch('/play', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    player_name: playerName,
                    player_move: playerMove
                })
            })
            .then(response => response.json())
            .then(data => {
                let resultClass = data.result;
                let resultIcon = data.result === 'win' ? '🎉' : data.result === 'loss' ? '😢' : '🤝';
                let moveIcons = {
                    'rock': '🪨',
                    'paper': '📄', 
                    'scissors': '✂️'
                };
                
                resultDiv.innerHTML = 
                    '<div class="' + resultClass + '" style="padding: 20px; border-radius: 10px;">' +
                    '<div style="font-size: 2em; margin-bottom: 10px;">' + resultIcon + '</div>' +
                    '<div style="font-size: 1.5em; margin-bottom: 10px;">' + data.result.toUpperCase() + '!</div>' +
                    '<div>You: ' + moveIcons[data.player_move] + ' ' + data.player_move + ' | Computer: ' + moveIcons[data.computer_move] + ' ' + data.computer_move + '</div>' +
                    '</div>';
                
                // Update stats dynamically after game result
                setTimeout(() => {
                    updateStats();
                }, 500);
                
                // Hide result after showing for a few seconds
                setTimeout(() => {
                    resultDiv.style.display = 'none';
                }, 3000);
            })
            .catch(error => {
                console.error('Error:', error);
                resultDiv.innerHTML = '<div style="color: red;">Error playing game. Please try again.</div>';
            });
        }
        
        // Auto-update stats every 15 seconds (reduced from 30)
        setInterval(updateStats, 15000);
        
        // Initial stats load
        document.addEventListener('DOMContentLoaded', () => {
            setTimeout(updateStats, 1000);
        });
    </script>
</body>
</html>
`

func init() {
	players = make(map[string]*Player)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/api/stats", statsHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("🎮 Rock Paper Scissors Arena starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	stats := getGameStats()
	
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
	
	tmpl, err := template.New("game").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, stats)
	if err != nil {
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
	
	// Update player stats
	player, exists := players[request.PlayerName]
	if !exists {
		player = &Player{Name: request.PlayerName}
		players[request.PlayerName] = player
	}
	
	player.Total++
	switch result {
	case "win":
		player.Wins++
	case "loss":
		player.Losses++
	case "draw":
		player.Draws++
	}
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
	
	// Calculate move statistics
	moveStats := map[string]int{
		"rock":     0,
		"paper":    0,
		"scissors": 0,
	}
	
	winStats := map[string]int{
		"win":  0,
		"loss": 0,
		"draw": 0,
	}
	
	for _, game := range games {
		moveStats[game.PlayerMove]++
		winStats[game.Result]++
	}
	
	// Get recent games (last 10)
	recentGames := make([]Game, 0)
start := len(games) - 10
if start < 0 {
	start = 0
}
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
		TotalGames:   len(games),
		TotalPlayers: len(players),
		RecentGames:  recentGames,
		Leaderboard:  leaderboard,
		MoveStats:    moveStats,
		WinStats:     winStats,
	}
}
