// Game state management
let statsUpdateInProgress = false;

// Utility function to safely create text nodes and elements
function createTextElement(tag, text, className = null) {
    const element = document.createElement(tag);
    if (className) element.className = className;
    element.textContent = text; // Always use textContent for user data
    return element;
}

// Utility function to sanitize and validate data
function sanitizePlayerName(name) {
    if (typeof name !== 'string') return 'Unknown';
    return name.trim().substring(0, 50); // Limit length and trim
}

// Update statistics from API
function updateStats() {
    if (statsUpdateInProgress) return;
    statsUpdateInProgress = true;
    
    fetch('/api/stats')
        .then(response => response.json())
        .then(stats => {
            updateStatCards(stats);
            updateLeaderboard(stats.leaderboard);
            updateRecentGames(stats.recent_games);
        })
        .catch(error => {
            console.error('Error updating stats:', error);
        })
        .finally(() => {
            statsUpdateInProgress = false;
        });
}

// Update stat cards with new data
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

// Update leaderboard with new data
function updateLeaderboard(leaderboard) {
    const leaderboardContainer = document.querySelector('.leaderboard');
    if (!leaderboardContainer) return;
    
    // Keep the header, replace the content
    const header = leaderboardContainer.querySelector('.section-title');
    leaderboardContainer.innerHTML = '';
    leaderboardContainer.appendChild(header);
    
    leaderboard.forEach((player, index) => {
        const winRate = player.total > 0 ? (player.wins / player.total * 100).toFixed(1) : '0.0';
        
        const playerDiv = document.createElement('div');
        playerDiv.className = 'leaderboard-item';
        
        // Create rank element
        const rankDiv = document.createElement('div');
        rankDiv.className = 'rank';
        rankDiv.textContent = '#' + (index + 1);
        
        // Create player stats container
        const playerStatsDiv = document.createElement('div');
        playerStatsDiv.className = 'player-stats';
        
        // Create player name element (sanitized)
        const playerNameDiv = createTextElement('div', sanitizePlayerName(player.name), 'player-name');
        
        // Create player details element
        const playerDetailsDiv = document.createElement('div');
        playerDetailsDiv.className = 'player-details';
        playerDetailsDiv.textContent = player.total + ' games • ' + player.wins + 'W ' + player.losses + 'L ' + player.draws + 'D';
        
        // Create win rate element
        const winRateDiv = document.createElement('div');
        winRateDiv.className = 'win-rate';
        winRateDiv.textContent = winRate + '%';
        
        // Assemble the elements
        playerStatsDiv.appendChild(playerNameDiv);
        playerStatsDiv.appendChild(playerDetailsDiv);
        playerDiv.appendChild(rankDiv);
        playerDiv.appendChild(playerStatsDiv);
        playerDiv.appendChild(winRateDiv);
        
        leaderboardContainer.appendChild(playerDiv);
    });
}

// Update recent games with new data
function updateRecentGames(recentGames) {
    const recentContainer = document.querySelector('.recent-games');
    if (!recentContainer) return;
    
    // Keep the header, replace the content
    const header = recentContainer.querySelector('.section-title');
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
        
        // Create game header container
        const gameHeaderDiv = document.createElement('div');
        gameHeaderDiv.className = 'game-header';
        
        // Create game moves element (sanitized)
        const gameMovesDiv = document.createElement('div');
        gameMovesDiv.className = 'game-moves';
        
        const playerNameStrong = createTextElement('strong', sanitizePlayerName(game.player_name));
        
        gameMovesDiv.appendChild(playerNameStrong);
        gameMovesDiv.appendChild(document.createTextNode(': ' + game.player_move + ' vs ' + game.computer_move));
        
        // Create result element
        const gameResultSpan = document.createElement('span');
        gameResultSpan.className = 'game-result ' + game.result;
        gameResultSpan.textContent = game.result.toUpperCase();
        
        // Create game meta element
        const gameMetaDiv = document.createElement('div');
        gameMetaDiv.className = 'game-meta';
        gameMetaDiv.textContent = timestamp + ' • ' + game.player_ip;
        
        // Assemble the elements
        gameHeaderDiv.appendChild(gameMovesDiv);
        gameHeaderDiv.appendChild(gameResultSpan);
        gameDiv.appendChild(gameHeaderDiv);
        gameDiv.appendChild(gameMetaDiv);
        
        recentContainer.appendChild(gameDiv);
    });
}

// Play a game move
function playGame(playerMove) {
    const playerName = document.getElementById('playerName').value.trim();
    if (!playerName) {
        alert('Please enter your name first!');
        return;
    }
    
    const resultDiv = document.getElementById('gameResult');
    resultDiv.classList.remove('hidden');
    
    // Clear previous content and add loading message safely
    resultDiv.innerHTML = '';
    const loadingDiv = document.createElement('div');
    loadingDiv.style.color = '#667eea';
    loadingDiv.textContent = '🎲 Playing...';
    resultDiv.appendChild(loadingDiv);
    
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
        displayGameResult(data, resultDiv);
        
        // Update stats dynamically after game result
        setTimeout(() => {
            updateStats();
        }, 500);
        
        // Hide result after showing for a few seconds
        setTimeout(() => {
            resultDiv.classList.add('hidden');
        }, 3000);
    })
    .catch(error => {
        console.error('Error:', error);
        resultDiv.innerHTML = '';
        const errorDiv = document.createElement('div');
        errorDiv.style.color = 'red';
        errorDiv.textContent = 'Error playing game. Please try again.';
        resultDiv.appendChild(errorDiv);
    });
}

// Display game result with proper styling
function displayGameResult(data, resultDiv) {
    const resultClass = data.result;
    const resultIcon = data.result === 'win' ? '🎉' : data.result === 'loss' ? '😢' : '🤝';
    const moveIcons = {
        'rock': '🪨',
        'paper': '📄', 
        'scissors': '✂️'
    };
    
    // Clear previous content
    resultDiv.innerHTML = '';
    
    // Create main result container
    const resultContainer = document.createElement('div');
    resultContainer.className = resultClass;
    resultContainer.style.padding = '20px';
    resultContainer.style.borderRadius = '10px';
    
    // Create icon element
    const iconDiv = document.createElement('div');
    iconDiv.style.fontSize = '2em';
    iconDiv.style.marginBottom = '10px';
    iconDiv.textContent = resultIcon;
    
    // Create result text element
    const resultTextDiv = document.createElement('div');
    resultTextDiv.style.fontSize = '1.5em';
    resultTextDiv.style.marginBottom = '10px';
    resultTextDiv.textContent = data.result.toUpperCase() + '!';
    
    // Create moves element
    const movesDiv = document.createElement('div');
    movesDiv.textContent = 'You: ' + moveIcons[data.player_move] + ' ' + data.player_move + 
                          ' | Computer: ' + moveIcons[data.computer_move] + ' ' + data.computer_move;
    
    // Assemble the result
    resultContainer.appendChild(iconDiv);
    resultContainer.appendChild(resultTextDiv);
    resultContainer.appendChild(movesDiv);
    resultDiv.appendChild(resultContainer);
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    // Initial stats load
    setTimeout(updateStats, 1000);
});

// Auto-update stats every 15 seconds
setInterval(updateStats, 15000);