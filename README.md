# Silly Mini Games - Knucklebones
The backend for [Silly Mini Game](https://www.sillyminigames.com), built with Go and featuring real-time gameplay via WebSockets, Google OAuth authentication, and AI opponents. Silly Mini Games is planned as a game hub to play fun, simple, and lesser-known mini games with your friends! Currently only the game of Knucklebones is available, but more games are coming.
## Motivation
I wanted to play Knucklebones with my friends! The video game Cult of the Lamb has no online option, and I couldn't find a way to play this game with my friends, so I decided to make a web app for it. Although Silly Mini Games currently supports Knucklebones only, all of the non-game-logic code can be reused to build different games, which is what I intend to do in the future.
## Features

- 🎲 **Multiple Game Modes**
  - Online multiplayer (real-time via WebSocket)
  - Local pass-and-play
  - Computer opponent with 3 difficulty levels

- 🔐 **Authentication**
  - Email/password registration with verification
  - Google OAuth integration
  - JWT-based authentication with refresh tokens

- 🎮 **Real-time Gameplay**
  - WebSocket connections for live game updates
  - Instant move broadcasting
  - Player join notifications

- 🤖 **Smart AI Opponent**
  - Easy, medium, and hard difficulty levels
  - Strategic move selection based on score optimization

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Authentication**: JWT, Google OAuth 2.0
- **Real-time**: WebSocket (Gorilla WebSocket)
- **Email**: SMTP verification emails

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Google OAuth 2.0 credentials (for Google sign-in)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/AradD7/Go-Knuclebones.git
cd Go-Knuclebones
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL database:
```bash
createdb knucklebones
```

4. Create a `.env` file in the root directory:
```env
TOKEN_SECRET=your-super-secret-jwt-key-here
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
DB_URL=postgres://user:password@localhost:5432/knucklebones?sslmode=disable
FRONTEND_URL=http://localhost:3000
PLATFORM=dev
```

5. Run database migrations with goose:
```bash
# Install goose if you haven't already
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
cd sql/schema
goose postgres "your-connection-string" up
cd ../..
```

6. Run the server:
```bash
go run .
```

The API will be available at `http://localhost:8080`

## Quick Start

### 1. Register a New Player

```bash
curl -X POST http://localhost:8080/api/players/new \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "email": "player1@example.com",
    "password": "securepassword123"
  }'
```

### 2. Verify Email

Check your email for the verification token, then:

```bash
curl -X POST http://localhost:8080/api/players/verify \
  -H "Content-Type: application/json" \
  -d '{
    "token": "your-verification-token"
  }'
```

This will return your JWT token and refresh token.

### 3. Create a Game

```bash
curl -X GET http://localhost:8080/api/games/new \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Play vs Computer

```bash
curl -X POST http://localhost:8080/api/games/computergame \
  -H "Content-Type: application/json" \
  -d '{
    "board1": [[0,0,0],[0,0,0],[0,0,0]],
    "board2": [[0,0,0],[0,0,0],[0,0,0]],
    "dice": 4,
    "row": 2,
    "col": 1,
    "difficulty": "medium"
  }'
```

## Game Rules

**Knucklebones** is a strategic dice game played on 3x3 grids:

1. Players take turns rolling a die (1-6) and placing it on their board
2. Dice must be placed in the lowest available row of a column
3. When you place a die, all matching dice in your opponent's same column are removed
4. Score is calculated per column: each unique value contributes `value × count × count`
5. Game ends when one player's board is full
6. Highest score wins

**Example Score Calculation:**
- Column with [3, 3, 5]: 
  - 3 appears twice: 3 × 2 × 2 = 12
  - 5 appears once: 5 × 1 × 1 = 5
  - Total: 17 points

## API Documentation

Full API documentation with all endpoints is available in [API.md](./API.md).

Quick links:
- [Authentication Endpoints](./API.md#authentication)
- [Player Management](./API.md#players)
- [Game Endpoints](./API.md#games)
- [WebSocket Connection](./API.md#websocket)

## Project Structure

```
.
├── README.md
├── go.mod
├── go.sum
├── main.go                          # Application entry point & server setup
├── handler_*.go                     # HTTP endpoint handlers
├── websocket.go                     # WebSocket implementation
├── json.go                          # JSON response helpers
├── reset.go                         # Database reset (dev only)
├── move_test.go                     # Game logic tests
├── index.html                       # Static file
├── internal/
│   ├── auth/                        # Authentication & JWT
│   │   ├── auth_test.go
│   │   ├── hash.go                  # Password hashing
│   │   ├── jwt.go                   # JWT token generation/validation
│   │   └── refresh_token.go         # Refresh token generation
│   ├── database/                    # Database queries (sqlc generated)
│   │   ├── *.sql.go                 # Generated query functions
│   │   ├── db.go
│   │   └── models.go
│   └── verification/
│       └── email_verification.go    # Email verification logic
├── sql/
│   ├── queries/                     # SQL query definitions (for sqlc)
│   │   ├── 001_players.sql
│   │   ├── 002_boards.sql
│   │   ├── 003_games.sql
│   │   ├── 004_refresh_token.sql
│   │   ├── 005_purge.sql
│   │   └── 006_verification_token.sql
│   └── schema/                      # Database migrations (goose)
│       ├── 001_players.sql
│       ├── 002_boards.sql
│       ├── 003_games.sql
│       ├── 004_link_boards_to_games.sql
│       ├── 005_refresh_tokens.sql
│       ├── 006_add_scores_to_boards.sql
│       ├── 007_make_board2_nullable_game.sql
│       ├── 008_change_board_type.sql
│       ├── 009_add_turn_to_games.sql
│       ├── 010_add_display_name_to_player.sql
│       ├── 011_add_default_value_to_avatar.sql
│       ├── 012_add_google_auth.sql
│       ├── 013_add_email_verified_to_players.sql
│       └── 014_verification_tokens.sql
└── sqlc.yaml                        # sqlc configuration
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `TOKEN_SECRET` | Yes | Secret key for JWT signing |
| `GOOGLE_CLIENT_ID` | Yes | Google OAuth client ID |
| `DB_URL` | Yes | PostgreSQL connection string |
| `FRONTEND_URL` | Yes | Frontend URL for CORS |
| `PLATFORM` | No | Set to `dev` to enable admin endpoints |

## Development

### Running Tests

```bash
go test ./...
```

### Database Migrations (Goose)

```bash
# Create a new migration
cd sql/schema
goose create add_new_feature sql

# Run pending migrations
goose postgres "$DB_URL" up

# Rollback last migration
goose postgres "$DB_URL" down

# Check migration status
goose postgres "$DB_URL" status
```

### Generate Database Code (sqlc)

After modifying queries in `sql/queries/*.sql`:

```bash
sqlc generate
```

This will regenerate the Go code in `internal/database/`.

### Reset Database (Dev Only)

```bash
curl -X POST http://localhost:8080/admin/reset
```

**Note:** Only works when `PLATFORM=dev`

## WebSocket Usage

Connect to a game's WebSocket for real-time updates:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/games/YOUR_GAME_ID');

// Authenticate immediately after connection
ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'auth',
    token: 'YOUR_JWT_TOKEN'
  }));
};

// Handle incoming messages
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  if (data.type === 'refresh') {
    // Opponent made a move - refresh game state
  } else if (data.type === 'joined') {
    // Second player joined
    console.log(`${data.display_name} joined the game`);
  } else if (data.type === 'roll') {
    // Dice was rolled
    console.log(`Dice rolled: ${data.dice}`);
  }
};
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -m 'Add some new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request


## Acknowledgments

- Game concept inspired by the dice game "Knucklebones" from Cult of the Lamb
- Built with [Gorilla WebSocket](https://github.com/gorilla/websocket)
- Database layer generated with [sqlc](https://github.com/kyleconroy/sqlc)
