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

var pacman = { x: 1, y: 14, targetX: 1, targetY: 14, dir: "right",frame:0 }
var Ghosts = [
    { x: 1, y: 14, targetX: 1, targetY: 14, dir: "right", name: "1", color: "#ffb751" },
    { x: 1, y: 14, targetX: 1, targetY: 14, dir: "right", name: "2", color: "#ff0000" },
    { x: 1, y: 14, targetX: 1, targetY: 14, dir: "right", name: "3", color: "#00ffff" },
    { x: 1, y: 14, targetX: 1, targetY: 14, dir: "right", name: "4", color: "#de9751" },
]
var started = false
var powerUped = false
const images = {}

var users = 0

const flipframe = () =>{
    if (pacman.frame === 0) {
        pacman.frame = 1;
    } else {
        pacman.frame = 0;
    }
}

const modifyUsers = (number)=>{
    users = number
    startButton = document.getElementById("startGameButton");
    buttonOverlay = document.getElementById("buttonOverlay");
    if (users <= 1 ){
        
        console.log("disabling",startButton)
        if (startButton){
            startButton.classList.add("disabled");
            buttonOverlay.style.display = "flex";
        }

    }else{
        buttonOverlay.style.display = "None";
        startButton.classList.remove("disabled");
    }
}



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

var timer;

document.addEventListener("DOMContentLoaded", () => {


    timer = document.getElementById("timer");
    timer.display = "None";


    gameSpeedSlider = document.getElementById("gameSpeedSlider");
    gameSpeedValue = document.getElementById("gameSpeedValue");
    gameSpeedValue.textContent = gameSpeedSlider.value
    gameSpeedSlider.addEventListener("input", (event) => {
        gameSpeedValue.textContent = event.target.value;
    })

    DurationSlider = document.getElementById("DurationSlider");
    DurationValue = document.getElementById("DurationValue");
    DurationValue.textContent = DurationSlider.value
    DurationSlider.addEventListener("input", (event) => {
        DurationValue.textContent = event.target.value;
    })

    loadImages();

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

    inviteLink = document.getElementById("LobbyId");
    inviteLink.textContent = "#"+ LobbyID;


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
        if (users <= 1){
            return
        }
        Message = {
            type: "StartGame",
            options :{
                game_speed:1000- parseInt(gameSpeedSlider.value),
                duration: parseInt(DurationSlider.value),
            },
            lobbyID: LobbyID,
        }
        console.log("Start game message:", Message);
        ws.send(JSON.stringify(Message));
    })

    if (startButton){
        document.addEventListener("keydown", (event) => {
            if (event.key === "Enter" && !startButton.disabled) {
                startButton.click();
            }
        })
    }


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
        ".": "#9797ae",
    };

    for (let i = 0; i < board.length; i++) {
        for (let j = 0; j < board[i].length; j++) {
            const char = board[i][j];
            ctx.fillStyle = mapping[char] || "#9797ae";



            ctx.fillRect(j * PixelWidth, i * PixelHeight, PixelWidth + 0.4, PixelHeight + 0.4);
            if (char === ".") {
                ctx.fillStyle = "#ff0";
                ctx.beginPath();
                ctx.arc(
                    j * PixelWidth + PixelWidth / 2,
                    i * PixelHeight + PixelHeight / 2,
                    PixelWidth / 6,
                    0,
                    2 * Math.PI
                );
                ctx.fill();

            }
            if (char === "0"){
                  ctx.fillStyle = "#ff0";
                ctx.beginPath();
                ctx.arc(
                    j * PixelWidth + PixelWidth / 2,
                    i * PixelHeight + PixelHeight / 2,
                    PixelWidth / 3,
                    0,
                    2 * Math.PI
                );
                ctx.fill();
                
            }
        }
    }

    // ctx.fillStyle = "#ffff00";
    // ctx.beginPath()
    // ctx.arc(
    //     pacman.targetX * PixelWidth + PixelWidth / 2,
    //     pacman.targetY * PixelHeight + PixelHeight / 2,
    //     PixelWidth / 2 - 2,
    //     0,
    //     2 * Math.PI
    // )
    // ctx.fill()


    ctx.drawImage(
        images[`0-${pacman.dir || 0}-${pacman.frame || 0}`],
        pacman.targetX * PixelWidth,
        pacman.targetY * PixelHeight,
        PixelWidth,
        PixelHeight,
    )
    flipframe()



    Ghosts.forEach(g => {

        var name = g.name
        var dir = g.dir || 0
        
        if (powerUped ){
            name = "5"
            dir = Math.random() > 0.5 ? 2 : 1; 
        }

        ctx.drawImage(
            images[`${name}-${dir}`],
            g.targetX * PixelWidth,
            g.targetY * PixelHeight,
            PixelWidth,
            PixelHeight,
        )

        // ctx.fillStyle = g.color;
        // ctx.beginPath() 
        // ctx.arc(
        //     g.targetX * PixelWidth + PixelWidth / 2,
        //     g.targetY * PixelHeight + PixelHeight / 2,
        //     PixelWidth / 2 - 2,
        //     0,
        //     2 * Math.PI
        // )
        // ctx.fill()
    })

}


const HandleWsMessage = (message) => {

    const LobbyID = window.location.href.split("/").pop();

    const data = JSON.parse(message)
    console.log(data.type)

    if (data.type === "UserInfoUpdate") {
        console.log("User info update received:", data.users);
        modifyUsers(data.users.length)
        UpdateUsers(data.users);
    }

    if (data.type === "powerUp") {
        powerUped = data.status;
        console.log("Power-up received:", data);
    }

    if (data.type === "sound"){
        playSound(data.sound);

    }

    if (data.type === "HostUpdate") {
        console.log("Host update received:", data.host);
        IsHost = true;
        ToggleHost();
    }

    if (data.type === "GameEnd"){
        console.log(data);
        OverOverlay = document.getElementById("GameOverOverlay");
        OverOverlay.style.display = "flex";
        const Winner = document.getElementById("winner");
        Winner.textContent = `Winner: ${data.winner}`;
        const scores = document.getElementById("Scores");
        var i = 1
        // data.scores.map((user, score) => {
           
        // })

        Object.keys(data.scores).forEach(user => {
            const score = data.scores[user];
             const scoreElement = document.createElement("div");
            scoreElement.className = "ScoreItem";
            const playerName = document.createElement("h2");
            const playerNameInner = document.createElement("span");
            i = Object.keys(data.scores).indexOf(user) + 1
            playerNameInner.textContent = `${i}. ${user}`
            playerName.appendChild(playerNameInner);
            const scoreh2 = document.createElement("h2");
            scoreh2.className = "score";
            scoreh2.textContent = score;
            scoreElement.appendChild(playerName);
            scoreElement.appendChild(scoreh2);
            scores.appendChild(scoreElement);

        })
    }

    if (data.type === "timer"){
        console.log(data)
        timer.style.display = "flex";
        timer.textContent = data.val;
    }

    if (data.type === "StartAlert") {
        started = true
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
        UpdateCanvas(data.board);
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

    console.log("Updating users:", users);

    users.forEach(e => {
        var id
        var imgsrc
        if (!started){
            id = e.enemy || 0;
            imgsrc = `/static/images/ghost.png`;
        }else if  (e.pacman) {
            id = 0
            imgsrc = "/static/images/sprites/0/3-0.png";
        }
        else {
            id = e.enemy
            imgsrc = `/static/images/sprites/${id}/1.png`;
        }
        const div = document.createElement("div");
        div.className = "PlayerItem";
        if (e.you === true){
            div.className = "PlayerItem You"
        }
        const Title = document.createElement("h1");
        Title.className = "PlayerTitle"
        host = e.host ? "*" : "";
        Title.innerHTML = `${host}${e.name}<img src="${imgsrc}" alt="Ghost Icon" class="ghost-icon">`
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


const loadImages = () => {
    const path = "/static/images/sprites/";
    for (let i = 0; i <= 5; i++) {
        if (i === 5){
            for (let j=0; j <= 1; j++) {
                const img = new Image();
                img.src = `${path}${i}/${j+1}.png?v=2`;
                images[`${i}-${j}`] = img;
            }
        }
        for (let j = 0; j <= 3; j++) {

            if (i === 0) {
                for (let k = 0; k <= 1; k++) {
                    const img = new Image();
                    img.src = `${path}${i}/${j}-${k}.png?v=1`;
                    images[`${i}-${j}-${k}`] = img;
                }
            } else {


                const img = new Image();
                img.src = `${path}${i}/${j}.png`;
                images[`${i}-${j}`] = img;
            }
        }
    }
}

var eatSoundStatus = 0 

const playSound = (sound) => {
    var audio = new Audio(`/static/sounds/${sound}`);
    audio.volume = 0.5;
    audio.play().catch(error => {
        console.error("Error playing sound:", error);
    }); 


}