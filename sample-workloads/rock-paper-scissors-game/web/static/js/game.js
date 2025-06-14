// Game state management
let statsUpdateInProgress = false;

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
        playerDiv.innerHTML = 
            '<div class="rank">#' + (index + 1) + '</div>' +
            '<div class="player-stats">' +
                '<div class="player-name">' + player.name + '</div>' +
                '<div class="player-details">' +
                    player.total + ' games • ' + player.wins + 'W ' + player.losses + 'L ' + player.draws + 'D' +
                '</div>' +
            '</div>' +
            '<div class="win-rate">' + winRate + '%</div>';
        
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
        
        gameDiv.innerHTML = 
            '<div class="game-header">' +
                '<div class="game-moves">' +
                    '<strong>' + game.player_name + '</strong>: ' + game.player_move + ' vs ' + game.computer_move +
                '</div>' +
                '<span class="game-result ' + game.result + '">' + game.result.toUpperCase() + '</span>' +
            '</div>' +
            '<div class="game-meta">' +
                timestamp + ' • ' + game.player_ip +
            '</div>';
        
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
        resultDiv.innerHTML = '<div style="color: red;">Error playing game. Please try again.</div>';
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
    
    resultDiv.innerHTML = 
        '<div class="' + resultClass + '" style="padding: 20px; border-radius: 10px;">' +
        '<div style="font-size: 2em; margin-bottom: 10px;">' + resultIcon + '</div>' +
        '<div style="font-size: 1.5em; margin-bottom: 10px;">' + data.result.toUpperCase() + '!</div>' +
        '<div>You: ' + moveIcons[data.player_move] + ' ' + data.player_move + ' | Computer: ' + moveIcons[data.computer_move] + ' ' + data.computer_move + '</div>' +
        '</div>';
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    // Initial stats load
    setTimeout(updateStats, 1000);
});

// Auto-update stats every 15 seconds
setInterval(updateStats, 15000);