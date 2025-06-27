# <span style="color: #2977F5;">Ex</span><span style="color: #ffff00;">Pacman</span>

ExPacman is a multiplayer reimagination of the classic Pacman game. Each player controls either the pacman or one of the ghosts, and the goal is to collect all the dots while avoiding being caught by the ghosts. Written entirely in golang and vanilla html/js, it is high optimized for performance on any device.

Check it out at [https://expacman.com](https://expacman.com).


## Features
- Real-time multiplayer with up to 5 players
- Optimized for performance on any device
- Simple and intuitive controls
- High fault tolerance
- Open source and free to use

## Running it locally
To run ExPacman locally, you need to have Go installed on your machine. Follow these steps:
1. Clone the repository:
   ```bash
   git clone
2. Move into the project directory:
   ```bash
   cd expacman/backend
   ```
3. Install the dependencies:
   ```bash
    go mod tidy
    ```
4. Start the server:
   ```bash
    go run *.go
    ```
5. Open your web browser and go to `http://localhost:8080` to play the game.