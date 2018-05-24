# Scoreboard-API

An API for the scoreboard app to handle MongoDB interactions

## Routes:
- `/` - home page, simply returns 200 and "Welcome!"
- `/health` - health page, returns 200 if application is running
- `/games` - returns all active games
- `/games/recent` - gets most recent active game
- `/games/code/{code}` - gets a game by its unique code
- `/games/name/{name}` - gets a game by its name
- `/games/new`* - creates a new game
    - `?winningScore={num}` sets what the winning score should be for the new game
- `/games/delete`* - deletes a game (most recent if no code or name is specified)
- `/set-red/{num}`* - sets red score to {num} (most recent game if no code or name is specified)
- `/set-blue/{num}`* - sets blue score to {num} (most recent game if no code or name is specified)
- `/increment-red/{num}`* - increments red score by {num} or 1 (most recent game if no code or name is specified)
- `/increment-blue/{num}`* - increments blue score by {num} or 1 (most recent game if no code or name is specified)
    - ``?gameName={name}` or `&gameCode={code}` (or both) should be specified for all routes with *.