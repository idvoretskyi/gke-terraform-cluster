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

### Option 1: Using Makefile (Simplest)

```bash
# Deploy with latest version
make deploy PROJECT_ID=your-gcp-project-id

# Deploy with specific version
make deploy PROJECT_ID=your-gcp-project-id VERSION=v4.1.0

# Deploy with auto-generated git version
make auto-deploy PROJECT_ID=your-gcp-project-id

# See all available commands
make help
```

### Option 2: Automated Deployment Script

```bash
# Simple deployment with latest version
./deploy.sh -p your-gcp-project-id

# Deploy with specific version
./deploy.sh -p your-gcp-project-id -v v4.1.0

# Auto-generate version from git commit
./deploy.sh -p your-gcp-project-id --auto-version

# Build only (no deployment)
./deploy.sh -p your-gcp-project-id --build-only

# See all options
./deploy.sh --help
```

### Option 3: Manual Deployment
```bash
# Set your configuration
export PROJECT_ID=your-gcp-project-id
export VERSION=${VERSION:-latest}  # Use latest by default, or set custom version

# Build and push Docker image
docker buildx build --platform linux/amd64 \
  -t gcr.io/$PROJECT_ID/rock-paper-scissors-game:$VERSION . --push

# Update deployment.yaml with your project ID and version
sed -i "s/PROJECT_ID/$PROJECT_ID/g" deployment.yaml
sed -i "s/VERSION_TAG/$VERSION/g" deployment.yaml

# Deploy to Kubernetes
kubectl apply -f deployment.yaml

# Get external IP
kubectl get svc rock-paper-scissors-service
```

#### Version Management Examples
```bash
# Use latest (default)
export VERSION=latest

# Use semantic versioning
export VERSION=v4.0.1

# Use git commit hash
export VERSION=$(git rev-parse --short HEAD)

# Use timestamp-based version
export VERSION=$(date +%Y%m%d-%H%M%S)
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

## 🚀 Production Deployment Enhancements

### Advanced Kubernetes Features

The enhanced deployment includes several production-ready features:

#### 🔐 Security Enhancements
- **Security Context**: Runs as non-root user with restricted capabilities
- **Pod Security**: Read-only root filesystem, dropped capabilities
- **Network Policies**: Ingress/egress traffic control
- **Security Headers**: CSP, HSTS, XSS protection in HTTP responses

#### 📈 Scalability & Reliability
- **Horizontal Pod Autoscaler (HPA)**: Auto-scales based on CPU/memory usage
- **Pod Disruption Budget (PDB)**: Ensures minimum availability during updates
- **Rolling Updates**: Zero-downtime deployments with controlled rollout
- **Pod Anti-Affinity**: Spreads pods across different nodes for high availability

#### 🔍 Monitoring & Health Checks
- **Startup Probe**: Handles slow-starting containers
- **Liveness Probe**: Detects and restarts unhealthy containers
- **Readiness Probe**: Controls traffic routing to healthy pods
- **Graceful Shutdown**: Proper signal handling for clean termination

#### ⚡ Performance Optimizations
- **HTTP Compression**: Gzip compression for response bodies
- **Template Caching**: Pre-parsed HTML templates for faster responses
- **Request Logging**: Structured logging with timing information
- **Resource Limits**: Defined CPU/memory limits for predictable performance

### Environment-Specific Configurations

#### Development Environment
```bash
# Use development settings
cp terraform.tfvars.dev terraform.tfvars
terraform plan
terraform apply
```

#### Production Environment
```bash
# Use production settings
cp terraform.tfvars.prod terraform.tfvars
terraform plan
terraform apply
```

### Deployment Monitoring

Check deployment status:
```bash
# View deployment status
kubectl get deployments,pods,services,hpa,pdb -o wide

# Check pod autoscaling
kubectl get hpa rock-paper-scissors-hpa -w

# View application logs
kubectl logs -f deployment/rock-paper-scissors-game

# Check resource usage
kubectl top pods -l app=rock-paper-scissors-game
```

### Troubleshooting

Common issues and solutions:

#### Pod Startup Issues
```bash
# Check pod events
kubectl describe pod <pod-name>

# Check startup probe
kubectl logs <pod-name> --previous
```

#### Network Connectivity
```bash
# Test network policy
kubectl exec -it <pod-name> -- wget -q -O- http://google.com

# Check service endpoints
kubectl get endpoints rock-paper-scissors-service
```

#### Resource Constraints
```bash
# Check resource usage
kubectl top nodes
kubectl top pods

# View resource limits
kubectl describe deployment rock-paper-scissors-game
```

