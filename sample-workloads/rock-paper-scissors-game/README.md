# 🎮 Rock Paper Scissors Arena

An interactive, real-time Rock Paper Scissors game with statistics tracking and leaderboard functionality, deployed on Google Kubernetes Engine.

## Features

### 🎯 Game Features
- **Interactive Gameplay**: Play Rock Paper Scissors against the computer
- **Real-time Statistics**: Live tracking of game statistics and player performance
- **Leaderboard**: Competitive ranking system based on win rate and total games
- **Player Profiles**: Individual player statistics with wins, losses, and draws
- **Recent Games**: Display of the latest game results
- **Move Analytics**: Statistics showing most popular moves

### 🎨 User Experience
- **Responsive Design**: Works on desktop and mobile devices
- **Modern UI**: Beautiful gradient design with intuitive controls
- **Real-time Updates**: Automatic page refresh to show latest statistics
- **Interactive Animations**: Smooth hover effects and transitions
- **Emoji-based Moves**: Fun visual representation (🪨📄✂️)

### 📊 Statistics Tracked
- Total games played across all players
- Number of active players
- Individual player win/loss/draw records
- Win rate percentages
- Move frequency analysis
- Recent game history with timestamps and IP tracking

## Endpoints

- `/` - Main game interface (HTML)
- `/play` - POST endpoint for submitting moves (JSON API)
- `/api/stats` - GET statistics in JSON format
- `/health` - Health check endpoint

## API Usage

### Play a Game
```bash
curl -X POST http://EXTERNAL_IP/play \
  -H "Content-Type: application/json" \
  -d '{
    "player_name": "YourName",
    "player_move": "rock"
  }'
```

Valid moves: `rock`, `paper`, `scissors`

### Get Statistics
```bash
curl http://EXTERNAL_IP/api/stats
```

## Deployment

### Prerequisites
- GKE cluster running
- Docker configured for GCR
- kubectl configured for your cluster

### Quick Deploy
```bash
# Build and push Docker image
docker buildx build --platform linux/amd64 \
  -t gcr.io/PROJECT_ID/rock-paper-scissors-game:v2.0 . --push

# Deploy to Kubernetes
kubectl apply -f deployment.yaml

# Get external IP
kubectl get svc rock-paper-scissors-service
```

### Access the Game
Once deployed, access the game at: `http://EXTERNAL_IP`

## Game Rules

1. **Rock** (🪨) beats **Scissors** (✂️)
2. **Scissors** (✂️) beats **Paper** (📄)
3. **Paper** (📄) beats **Rock** (🪨)
4. Same moves result in a **Draw** (🤝)

## Architecture

- **Backend**: Go HTTP server with in-memory storage
- **Frontend**: HTML/CSS/JavaScript with responsive design
- **Deployment**: Kubernetes with 2 replicas for high availability
- **Load Balancer**: GCP LoadBalancer for external access
- **Health Checks**: Kubernetes liveness and readiness probes

## Configuration

- **Port**: 8080 (configurable via PORT environment variable)
- **Replicas**: 2 (for high availability)
- **Resources**: 128Mi-256Mi memory, 100m-200m CPU
- **Storage**: In-memory (resets on pod restart)

## Development

### Local Testing
```bash
# Run locally
go run main.go

# Access at http://localhost:8080
```

### Docker Testing
```bash
# Build image
docker build -t rock-paper-scissors-game .

# Run container
docker run -p 8080:8080 rock-paper-scissors-game
```

## Future Enhancements

- 🗄️ Persistent storage with database
- 👥 Multiplayer functionality
- 🎨 Custom themes and avatars
- 📱 Progressive Web App (PWA) support
- 🔔 Real-time notifications
- 🏆 Tournaments and competitions
- 📈 Advanced analytics and charts
