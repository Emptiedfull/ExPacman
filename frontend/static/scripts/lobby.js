const boardTemp = [
    "############################",
    "#............##............#",
    "#.####.#####.##.#####.####.#",
    "#.####.#####.##.#####.####.#",
    "#.####.#####.##.#####.####.#",
    "#..........................#",
    "#.####.##.########.##.####.#",
    "#.####.##.########.##.####.#",
    "#......##....##....##......#",
    "######.##### ## #####.######",
    "######.##### ## #####.######",
    "######.##          ##.######",
    "######.## ###  ### ##.######",
    "######.## #  b   # ##.######",
    " p     ##   c  d   ##       ",
    "######.## #   a  # ##.######",
    "######.## ###  ### ##.######",
    "######.##          ##.######",
    "######.## ######## ##.######",
    "######.## ######## ##.######",
    "#............##............#",
    "#.####.#####.##.#####.####.#",
    "#.####.#####.##.#####.####.#",
    "#o..##................##..o#",
    "###.##.##.########.##.##.###",
    "###.##.##.########.##.##.###",
    "#......##....##....##......#",
    "#.##########.##.##########.#",
    "#.##########.##.##########.#",
    "#..........................#",
    "############################",
]

var IsHost = false;
var PixelHeight
var PixelWidth
var Canvas;
var Role;
var ws

colors = ["#ffb751", "#ff0000", "#00ffff", "#de9751", "#ffb751"];

var pacman = { x: 1, y: 14, targetX: 1, targetY: 14, dir: "right" }
var Ghosts = [
    { x: 1, y: 14, targetX: 1, targetY: 14, dir: "right", name: "1", color: "#ffb751" },
    { x: 1, y: 14, targetX: 1, targetY: 14, dir: "right", name: "2", color: "#ff0000" },
    { x: 1, y: 14, targetX: 1, targetY: 14, dir: "right", name: "3", color: "#00ffff" },
    { x: 1, y: 14, targetX: 1, targetY: 14, dir: "right", name: "4", color: "#de9751" },
]



const ToggleHost = () => {
    if (IsHost === false) {
        GuestMode = document.getElementById("GuestMode");
        HostMode = document.getElementById("HostMode");
        GuestMode.style.display = "flex";
        HostMode.style.display = "None";
    }
    if (IsHost === true) {
        GuestMode = document.getElementById("GuestMode");
        HostMode = document.getElementById("HostMode");
        GuestMode.style.display = "None";
        HostMode.style.display = "flex";
    }
}

const Directions = {
    "ArrowUp": 0,
    "ArrowDown": 1,
    "ArrowLeft": 2,
    "ArrowRight": 3,
}

document.addEventListener("DOMContentLoaded", () => {

    SettingsContainer = document.getElementById("SettingsContainer");
    ToggleHost();

    const GameContainer = document.getElementById("gameContainer");
    GameContainer.style.display = "None"

    args = window.location.href.split("/");

    const Name = args.pop();
    const LobbyID = args.pop();
    console.log("Lobby ID:", LobbyID);

    const WsURL = `ws://localhost:8080/ws/${LobbyID}`;
    console.log("WebSocket URL:", WsURL);


    ws = new WebSocket(WsURL);


    document.addEventListener("keydown", (event) => {


        if (event.key in Directions) {
            event.preventDefault();
            const direction = Directions[event.key];


            const Message = {
                type: "MoveState",
                direction: direction,
            };
            ws.send(JSON.stringify(Message));
            console.log("Move state sent:", Message, ws.readyState);
        }
    })

    ws.onopen = () => {
        const NameResponse = {
            "id": LobbyID,
            "name": Name,
        }
        JSON.stringify(NameResponse);
        ws.send(JSON.stringify(NameResponse));
    }

    ws.onmessage = (event) => {
        event.preventDefault();
        HandleWsMessage(event.data);
    }

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        window.location.href = "/";
    };

    startButton = document.getElementById("startGameButton");
    startButton.addEventListener("click", () => {
        Message = {
            type: "StartGame",
            lobbyID: LobbyID,
        }
        ws.send(JSON.stringify(Message));
    })


});


const UpdateCanvas = (board) => {
    if (!Canvas) {
        console.error("Canvas not loaded yet");
        return;
    }

    const ctx = Canvas.getContext("2d");
    if (!ctx) {
        console.error("Canvas context not found");
        return;
    }

    ctx.clearRect(0, 0, Canvas.width, Canvas.height);

    const mapping = {
        "#": "#000",
        ".": "#2121ff",
        "o": "#ff0",
    };

    for (let i = 0; i < board.length; i++) {
        for (let j = 0; j < board[i].length; j++) {
            const char = board[i][j];
            ctx.fillStyle = mapping[char] || "#2121ff";



            ctx.fillRect(j * PixelWidth, i * PixelHeight, PixelWidth + 0.4, PixelHeight + 0.4);
            if (char === ".") {
                ctx.fillStyle = "#ff0";
                ctx.beginPath();
                ctx.arc(
                    j * PixelWidth + PixelWidth / 2,
                    i * PixelHeight + PixelHeight / 2,
                    PixelWidth / 4,
                    0,
                    2 * Math.PI
                );
                ctx.fill();

            }
        }
    }



    // Ghosts.forEach(g => {
    //     ctx.fillStyle = g.color;
    //     ctx.beginPath() 
    //     ctx.arc(
    //         g.x * PixelWidth + PixelWidth / 2,
    //         g.y * PixelHeight + PixelHeight / 2,
    //         PixelWidth / 2 - 2,
    //         0,
    //         2 * Math.PI
    //     )
    //     ctx.fill()
    // })

}


const HandleWsMessage = (message) => {

    const LobbyID = window.location.href.split("/").pop();

    const data = JSON.parse(message)

    if (data.type === "UserInfoUpdate") {
        UpdateUsers(data.users);
    }

    if (data.type === "HostUpdate") {
        console.log("Host update received:", data.host);
        IsHost = true;
        ToggleHost();
    }

    if (data.type === "StartAlert") {
        console.log("Start alert received:", data);
        Role = data.role;
        console.log("Role assigned:", Role);

        SettingsContainer = document.getElementById("SettingsContainer");
        SettingsContainer.style.display = "None";

        GameContainer = document.getElementById("gameContainer");
        GameContainer.style.display = "flex";

        Canvas = document.getElementById("gameCanvas");
        Canvas.style.display = "block";

        pacman.targetX = data.pacman.target_y;
        pacman.targetY = data.pacman.target_x;
        pacman.x = data.pacman.pos_y;
        pacman.y = data.pacman.pos_x;
        pacman.dir = data.pacman.dir;


        Ghosts = data.enemy.map(ghost => ({
            x: ghost.pos_y,
            y: ghost.pos_x,
            targetX: ghost.target_y,
            targetY: ghost.target_x,
            dir: ghost.dir,
            name: ghost.name,
            color: colors[parseInt(ghost.name) - 1] || "#000"

        }))

        loadCanvas();
    }

    if (data.type === "BoardUpdate") {

        data.scores.forEach((score, index) => {
            const scoreElement = document.getElementById(index.toString());
            if (scoreElement) {
                scoreElement.textContent = score;
            } else {
                // console.warn(`Score element with ID ${index} not found.`);
            }
        })

        pacman.targetX = data.pacman.target_y;
        pacman.targetY = data.pacman.target_x;
        pacman.dir = data.pacman.dir;


        Ghosts = data.enemy.map(ghost => ({
            targetX: ghost.target_y,
            targetY: ghost.target_x,
            x: Ghosts[Ghosts.findIndex(g => g.name === ghost.name)].x,
            y: Ghosts[Ghosts.findIndex(g => g.name === ghost.name)].y,
            dir: ghost.dir,
            name: ghost.name,
            color: colors[parseInt(ghost.name) - 1] || "#000"

        }))
        // console.log("Board update received:", data.board);
        //  UpdateCanvas(data.board); 
        animateEntities(data.board);
    }
}

function animateEntities(boardTemp) {
    if (
        Math.abs(pacman.x - pacman.targetX) <= 0.1 &&
        Ghosts.some(ghost => Math.abs(ghost.x - ghost.targetX) <= 0.1) &&
        Math.abs(pacman.y - pacman.targetY) <= 0.1 &&
        Ghosts.some(ghost => Math.abs(ghost.y - ghost.targetY) <= 0.1)
    ) {
        console.log("No animation")
        return
    }





    ctx = Canvas.getContext("2d");

    //1 cell/1 sec -> 0.004 animation speed
    //1 cell/0.5 sec -> 0.008 animation speed
    const speed = 0.016



    function moveEntity(entity) {

        if (entity.targetY === entity.y && entity.targetX === entity.x) {
            return;
        }

        if (
            Math.abs(entity.x - entity.targetX) < 0.1 && Math.abs(entity.x - entity.targetX) !== 0
        ) {
            entity.x = entity.targetX;
            console.log("clipped")
            return;
        }

        if (
            Math.abs(entity.y - entity.targetY) < 0.1 && Math.abs(entity.y - entity.targetY) !== 0
        ) {
            entity.y = entity.targetY;
            console.log("clipped")
            return;
        }

        DirX = 0
        DirY = 0

        if (entity.dir === 0) {
            if (entity.targetX != entity.x) {
                entity.x = entity.targetX;
                console.log("directioin clipped")
            }


            DirX = 0
            DirY = -1;
        } else if (entity.dir === 1) {
            if (entity.targetX != entity.x) {
                entity.x = entity.targetX
                console.log("directioin clipped")
            }
            DirX = 0
            DirY = 1;
        } else if (entity.dir === 2) {
            if (entity.targetY != entity.y) {
                entity.y = entity.targetY
                console.log("directioin clipped")
            }
            DirX = -1
            DirY = 0;
        } else if (entity.dir === 3) {
            if (entity.targetY != entity.y) {
                entity.y = entity.targetY;
                console.log("directioin clipped")
            }
            DirX = 1
            DirY = 0;
        }
        entity.x += DirX * speed;
        entity.y += DirY * speed;
    }

    moveEntity(pacman);
    UpdateCanvas(boardTemp);
    ctx.fillStyle = "#ffff00";
    ctx.beginPath()
    ctx.arc(
        pacman.x * PixelWidth + PixelWidth / 2,
        pacman.y * PixelHeight + PixelHeight / 2,
        PixelWidth / 2 - 2,
        0,
        2 * Math.PI
    )
    ctx.fill()
    Ghosts.forEach(moveEntity);

    Ghosts.forEach(g => {
        ctx.fillStyle = g.color;
        ctx.beginPath()
        ctx.arc(
            g.x * PixelWidth + PixelWidth / 2,
            g.y * PixelHeight + PixelHeight / 2,
            PixelWidth / 2 - 2,
            0,
            2 * Math.PI
        )
        ctx.fill()
    })

    





    if (
        Math.abs(pacman.x - pacman.targetX) > 0.1 ||
        Ghosts.some(ghost => Math.abs(ghost.x - ghost.targetX) > 0.1) ||
        Math.abs(pacman.y - pacman.targetY) > 0.1 ||
        Ghosts.some(ghost => Math.abs(ghost.y - ghost.targetY) > 0.1)
    ) {
        // UpdateCanvas(boardTemp);

        setTimeout(() => animateEntities(boardTemp), 1000 / 60);


        // animateEntities(boardTemp)
    }
}

const UpdateUsers = (users) => {
    const PlayerContainer = document.getElementById("playersContainer");
    PlayerContainer.innerHTML = "";



    users.forEach(e => {
        var id
        if (e.pacman) {
            id = 0
        }
        else {
            id = e.enemy
        }
        const div = document.createElement("div");
        div.className = "PlayerItem";
        const Title = document.createElement("h1");
        Title.className = "PlayerTitle"
        host = e.host ? "*" : "";
        Title.innerHTML = `${host}${e.name}<img src="/static/images/ghost.png" alt="Ghost Icon" class="ghost-icon">`
        const Score = document.createElement("h2")
        Score.className = "PlayerScore";
        Score.innerHTML = ` Score: <span class="yellow" id="${id}">0</span>`;
        div.appendChild(Title);
        div.appendChild(Score);
        PlayerContainer.appendChild(div);
    });
}


const loadCanvas = () => {
    const GameContainer = document.getElementById("gameContainer");

    Canvas = document.getElementById("gameCanvas");
    const ctx = Canvas.getContext("2d");
    if (!Canvas) {
        console.error("Canvas element not found");
        return;
    }



    const mapping = {
        "#": "#000",
        ".": "#fff",
        "o": "#ff0",

    };


    Canvas.height = GameContainer.clientHeight;
    Canvas.width = Canvas.height * (boardTemp[0].length / boardTemp.length);

    PixelHeight = Canvas.height / boardTemp.length;
    PixelWidth = Canvas.width / boardTemp[0].length;

    // console.log("Pixel dimensions:", PixelWidth, PixelHeight);



    for (let i = 0; i < boardTemp.length; i++) {
        for (let j = 0; j < boardTemp[i].length; j++) {
            const char = boardTemp[i][j];
            ctx.fillStyle = mapping[char] || "#fff"; // Default to white if char not found



            ctx.fillRect(j * PixelWidth, i * PixelHeight, PixelWidth + 1, PixelHeight + 1);
        }
    }
    ctx.imageSmoothingEnabled = false;


    ctx.fillStyle = "#ffff00";
    ctx.beginPath()
    ctx.arc(
        pacman.x * PixelWidth + PixelWidth / 2,
        pacman.y * PixelHeight + PixelHeight / 2,
        PixelWidth / 2 - 2,
        0,
        2 * Math.PI
    )
    ctx.fill()


    Ghosts.forEach(g => {
        ctx.fillStyle = g.color;
        ctx.beginPath()
        ctx.arc(
            g.x * PixelWidth + PixelWidth / 2,
            g.y * PixelHeight + PixelHeight / 2,
            PixelWidth / 2 - 2,
            0,
            2 * Math.PI
        )
        ctx.fill()
    })



    // console.log("Canvas loaded with dimensions:", Width, Height,PixelHeight,PixelWidth);

}