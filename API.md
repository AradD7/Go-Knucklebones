# Knucklebones Game API Documentation

## Table of Contents
- [Environment Variables](#environment-variables)
- [Health & Admin](#health--admin)
- [Authentication](#authentication)
- [Players](#players)
- [Tokens](#tokens)
- [Games](#games)
- [WebSocket](#websocket)

---

## Environment Variables

The following environment variables must be configured before running the API:

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `TOKEN_SECRET` | Yes | Secret key for JWT token signing | `your-super-secret-key-here` |
| `GOOGLE_CLIENT_ID` | Yes | Google OAuth 2.0 Client ID | `123456789-abc.apps.googleusercontent.com` |
| `DB_URL` | Yes | PostgreSQL database connection string | `postgres://user:pass@localhost:5432/dbname?sslmode=disable` |
| `FRONTEND_URL` | Yes | Frontend application URL for CORS | `http://localhost:3000` |
| `PLATFORM` | No | Environment mode (`dev` or `prod`) | `dev` |

### Example .env file
```env
TOKEN_SECRET=your-super-secret-key-here
GOOGLE_CLIENT_ID=123456789-abc.apps.googleusercontent.com
DB_URL=postgres://user:password@localhost:5432/knucklebones?sslmode=disable
FRONTEND_URL=http://localhost:3000
PLATFORM=dev
```

**Notes:**
- `PLATFORM=dev` enables the `/admin/reset` endpoint
- `FRONTEND_URL` is used for CORS and WebSocket origin validation
- Keep `TOKEN_SECRET` secure and never commit it to version control

---

## Health & Admin

### Health Check

<details>
<summary><b>GET</b> <code>/api/health</code> - Check API status</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | No |
| **Description** | Check if the API is running |

**Response:**
```json
{
  "status": "ok"
}
```

</details>

---

### Reset Database

<details>
<summary><b>POST</b> <code>/admin/reset</code> - Reset database (Dev only)</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | No (Dev mode only) |
| **Description** | Resets the entire database (development only) |

**Response:** `204 No Content`

**Notes:**
- Only available when `PLATFORM=dev`

</details>

---

## Authentication

### Google Sign In

<details>
<summary><b>POST</b> <code>/api/auth/google</code> - Authenticate with Google OAuth</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | No |
| **Description** | Authenticate using Google OAuth |

**Request Body:**
```json
{
  "id_token": "google_oauth_token_here"
}
```

**Response:**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "username": "username",
  "refresh_token": "token",
  "token": "jwt_token",
  "avatar": "avatar_url",
  "display_name": "Display Name"
}
```

</details>

---

## Players

### Create New Player

<details>
<summary><b>POST</b> <code>/api/players/new</code> - Register a new player</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | No |
| **Description** | Register a new player account |

**Request Body:**
```json
{
  "username": "player123",
  "email": "player@example.com",
  "password": "securepassword"
}
```

**Response:** `201 Created` with `null` body (verification email sent)

</details>

---

### Player Login

<details>
<summary><b>POST</b> <code>/api/players/login</code> - Login with credentials</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | No |
| **Description** | Login with username and password |

**Request Body:**
```json
{
  "username": "player123",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "username": "player123",
  "refresh_token": "refresh_token_here",
  "token": "jwt_token_here",
  "avatar": "avatar_url",
  "display_name": "Display Name",
  "email_verified": true
}
```

**Notes:**
- Returns `403 Forbidden` if email is not verified
- Requires email verification before login

</details>

---

### Get Player Info

<details>
<summary><b>GET</b> <code>/api/players/getplayer</code> - Get current player info</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | Yes (Bearer Token) |
| **Description** | Get current authenticated player's information |

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "id": "uuid",
  "username": "player123",
  "avatar": "avatar_url",
  "display_name": "Display Name"
}
```

</details>

---

### Update Player Profile

<details>
<summary><b>POST</b> <code>/api/players/update</code> - Update profile</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | Yes (Bearer Token) |
| **Description** | Update player's display name and avatar |

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "display_name": "New Display Name",
  "avatar": "new_avatar_url"
}
```

**Response:** `200 OK` with `null` body

</details>

---

### Verify Email

<details>
<summary><b>POST</b> <code>/api/players/verify</code> - Verify email address</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | No |
| **Description** | Verify email address using token from email |

**Request Body:**
```json
{
  "token": "verification_token_from_email"
}
```

**Response:**
```json
{
  "id": "uuid",
  "email": "player@example.com",
  "display_name": "Display Name",
  "avatar": "avatar_url",
  "refresh_token": "refresh_token",
  "token": "jwt_token"
}
```

**Notes:**
- Token expires after a certain time
- Automatically logs in the user after verification

</details>

---

### Resend Verification Email

<details>
<summary><b>POST</b> <code>/api/players/resendverification</code> - Resend verification</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | No |
| **Description** | Resend verification email |

**Request Body:**
```json
{
  "email": "player@example.com",
  "username": "player123"
}
```

**Response:**
```json
{
  "message": "Verification email sent"
}
```

**Notes:**
- Rate limited to once per 30 minutes
- Returns `429 Too Many Requests` if requested too frequently

</details>

---

## Tokens

### Refresh JWT Token

<details>
<summary><b>GET</b> <code>/api/tokens/refresh</code> - Get new JWT token</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | Yes (Refresh Token) |
| **Description** | Get a new JWT token using refresh token |

**Headers:**
```
Authorization: Bearer <refresh_token>
```

**Response:**
```json
{
  "token": "new_jwt_token"
}
```

</details>

---

### Revoke Refresh Token

<details>
<summary><b>GET</b> <code>/api/tokens/revoke</code> - Logout/revoke token</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | Yes (Refresh Token) |
| **Description** | Revoke/logout the current refresh token |

**Headers:**
```
Authorization: Bearer <refresh_token>
```

**Response:** `200 OK` with `null` body

</details>

---

## Games

### Create New Game

<details>
<summary><b>GET</b> <code>/api/games/new</code> - Create new multiplayer game</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | Yes (Bearer Token) |
| **Description** | Create a new game and wait for another player to join |

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "id": "game_uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "board1": [[0,0,0], [0,0,0], [0,0,0]],
  "board2": null
}
```

</details>

---

### Get All Player Games

<details>
<summary><b>GET</b> <code>/api/games</code> - List all player's games</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | Yes (Bearer Token) |
| **Description** | Get list of all game IDs the player is participating in |

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "ids": ["game_uuid_1", "game_uuid_2"]
}
```

</details>

---

### Get Specific Game

<details>
<summary><b>GET</b> <code>/api/games/{game_id}</code> - Get game details</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | Yes (Bearer Token) |
| **Description** | Get details of a specific game |

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**URL Parameters:**
- `game_id`: UUID of the game

**Response:**
```json
{
  "id": "game_uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "board1": [[1,2,3], [4,5,6], [1,2,3]],
  "board2": [[6,5,4], [3,2,1], [6,5,4]],
  "score1": 42,
  "score2": 38,
  "is_turn": true,
  "is_over": false
}
```

**Notes:**
- `board1` is always the current player's board
- `board2` is always the opponent's board
- `is_turn` indicates if it's the current player's turn

</details>

---

### Join Game

<details>
<summary><b>GET</b> <code>/api/games/{game_id}/join</code> - Join existing game</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | Yes (Bearer Token) |
| **Description** | Join an existing game as the second player |

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**URL Parameters:**
- `game_id`: UUID of the game to join

**Response:**
```json
{
  "id": "game_uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "board1": [[0,0,0], [0,0,0], [0,0,0]],
  "board2": [[0,0,0], [0,0,0], [0,0,0]],
  "is_turn": false,
  "opp_name": "Opponent Name",
  "opp_avatar": "opponent_avatar_url"
}
```

**Notes:**
- Randomly assigns who goes first
- Broadcasts join event to connected WebSocket clients

</details>

---

### Make Move

<details>
<summary><b>POST</b> <code>/api/games/move/{game_id}</code> - Make a move</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | Yes (Bearer Token) |
| **Description** | Make a move in an online game |

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**URL Parameters:**
- `game_id`: UUID of the game

**Request Body:**
```json
{
  "dice": 5,
  "row": 2,
  "col": 1
}
```

**Response:**
```json
{
  "board1": [[0,5,0], [0,0,0], [0,0,0]],
  "board2": [[0,0,0], [0,0,0], [0,0,0]],
  "score1": 25,
  "score2": 18,
  "is_over": false
}
```

**Notes:**
- Board indices are 0-based (0, 1, 2)
- Dice values are 1-6
- Automatically updates opponent's board (removes matching dice in same column)
- Broadcasts move to opponent via WebSocket
- Determines winner when board is full

</details>

---

### Local Game (Pass and Play)

<details>
<summary><b>POST</b> <code>/api/games/localgame</code> - Local multiplayer move</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | No |
| **Description** | Process a move in a local pass-and-play game |

**Request Body:**
```json
{
  "board1": [[1,0,0], [2,0,0], [3,0,0]],
  "board2": [[0,0,4], [0,0,5], [0,0,6]],
  "turn": "player1",
  "dice": 3,
  "row": 1,
  "col": 1
}
```

**Response:**
```json
{
  "board1": [[1,0,0], [2,3,0], [3,0,0]],
  "board2": [[0,0,4], [0,0,5], [0,0,6]],
  "score1": 14,
  "score2": 75,
  "is_over": false
}
```

**Notes:**
- `turn` must be either "player1" or "player2"
- No authentication required for local games
- Client manages game state

</details>

---

### Computer Game (vs AI)

<details>
<summary><b>POST</b> <code>/api/games/computergame</code> - Play vs computer</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | No |
| **Description** | Process player's move and get computer's response |

**Request Body:**
```json
{
  "board1": [[0,0,0], [0,0,0], [0,0,0]],
  "board2": [[0,0,0], [0,0,0], [0,0,0]],
  "dice": 4,
  "row": 2,
  "col": 1,
  "difficulty": "medium"
}
```

**Response:**
```json
{
  "board1": [[0,0,0], [0,0,0], [0,4,0]],
  "board2": [[0,0,0], [0,0,0], [0,0,0]],
  "next_board1": [[0,0,0], [0,0,0], [0,0,0]],
  "next_board2": [[0,0,0], [0,0,0], [3,0,0]],
  "score1": 16,
  "score2": 9,
  "next_score1": 16,
  "next_score2": 18,
  "next_dice": 3,
  "is_over": false,
  "is_over_next": false
}
```

**Notes:**
- `difficulty` can be: "easy", "medium", or "hard"
- Returns both the state after player's move and after computer's move
- `next_*` fields contain the state after computer's move
- Computer difficulty affects which move it selects from best to worst

</details>

---

### Roll Dice

<details>
<summary><b>GET</b> <code>/api/games/roll</code> - Roll dice</summary>

| Property | Value |
|----------|-------|
| **Auth Required** | Yes (Bearer Token) |
| **Description** | Roll a dice and broadcast to all players in the game |

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `game_id`: UUID of the game

**Example:**
```
GET /api/games/roll?game_id=550e8400-e29b-41d4-a716-446655440000
```

**Response:**
```json
{
  "dice": 4
}
```

**Notes:**
- Returns random number 1-6
- Broadcasts dice roll to all WebSocket connections for that game

</details>

---

## WebSocket

### Game WebSocket Connection

<details>
<summary><b>WebSocket</b> <code>/ws/games/{game_id}</code> - Real-time game updates</summary>

| Property | Value |
|----------|-------|
| **Protocol** | WebSocket |
| **Auth Required** | Yes (JWT via initial message) |
| **Description** | Real-time game updates via WebSocket |

**URL Parameters:**
- `game_id`: UUID of the game

**Connection Flow:**
1. Connect to WebSocket endpoint
2. Send authentication message immediately:
```json
{
  "type": "auth",
  "token": "jwt_token_here"
}
```

**Message Types Received:**

#### Refresh Event
```json
{
  "type": "refresh"
}
```
Sent when opponent makes a move. Client should fetch latest game state.

#### Joined Event
```json
{
  "type": "joined",
  "display_name": "Player Name",
  "avatar": "avatar_url"
}
```
Sent when a second player joins the game.

#### Roll Event
```json
{
  "type": "roll",
  "dice": 4
}
```
Sent when dice is rolled (via `/api/games/roll` endpoint).

</details>

---

## Game Board Format

The game board is represented as a 3x3 2D array:

```json
[
  [0, 0, 0],  // Row 0 (top)
  [0, 0, 0],  // Row 1 (middle)
  [0, 0, 0]   // Row 2 (bottom)
]
```

**Rules:**
- Each cell can contain a dice value (1-6) or 0 (empty)
- Dice must be placed in the lowest available row of a column
- When a player places a dice, matching dice in opponent's same column are removed
- Score is calculated per column: sum of (value × count × count) for each unique value

**Example Score Calculation:**
Column with [3, 3, 5]:
- 3 appears 2 times: 3 × 2 × 2 = 12
- 5 appears 1 time: 5 × 1 × 1 = 5
- Column score: 17

---

## Error Responses

All endpoints return errors in this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

**Common HTTP Status Codes:**
- `400 Bad Request` - Invalid input data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Action not allowed (e.g., unverified email)
- `404 Not Found` - Resource doesn't exist
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server-side error

---

## Authentication Flow

### Standard Registration & Login
1. **Register**: `POST /api/players/new` - Creates account, sends verification email
2. **Verify**: `POST /api/players/verify` - Verifies email, returns tokens
3. **Use JWT**: Include JWT in `Authorization: Bearer <token>` header for protected endpoints
4. **Refresh**: Use `GET /api/tokens/refresh` when JWT expires (60 min)

### Google OAuth Flow
1. **Authenticate with Google**: Get ID token from Google OAuth
2. **Sign in**: `POST /api/auth/google` with ID token
3. **Use tokens**: Same as standard flow

### Token Lifetimes
- **JWT Token**: 60 minutes
- **Refresh Token**: 60 days
- **Verification Token**: 24 hours
