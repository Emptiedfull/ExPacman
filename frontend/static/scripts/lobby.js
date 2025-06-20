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

colors = ["ffb7ff", "ff0000", "00ffff", "de9751","ffb7ff"];

var pacman = {x:1,y:14,targetX:1,targetY:14,dir:"right"}
var Ghosts = [
    {x:1,y:14,targetX:1,targetY:14,dir:"right",name:"1",color:"#ffb7ff"},
    {x:1,y:14,targetX:1,targetY:14,dir:"right",name:"2",color:"#ff0000"},
    {x:1,y:14,targetX:1,targetY:14,dir:"right",name:"3",color:"#00ffff"},
    {x:1,y:14,targetX:1,targetY:14,dir:"right",name:"4",color:"#de9751"},
]

console.log(boardTemp.length,boardTemp[0].length);

const ToggleHost = () => {
    if (IsHost === false){
        GuestMode = document.getElementById("GuestMode");
        HostMode = document.getElementById("HostMode");
        GuestMode.style.display = "flex";
        HostMode.style.display = "None";
    }
    if (IsHost === true){
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
    console.log("Settings Container:", SettingsContainer);
    console.log("Host status:", IsHost);
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
       
        console.log("Key pressed:", event.key);
        if (event.key in Directions) {
             event.preventDefault();
            const direction = Directions[event.key];
            

            const Message = {
                type: "MoveState",
                direction: direction,
            };
            ws.send(JSON.stringify(Message));
            console.log("Move state sent:", Message,ws.readyState);
        } else {
            console.log("Key pressed is not a valid direction:", event.key);
        }

    })

    ws.onopen = () =>{
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


const UpdateCanvas = (board)=>{
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
        ".": "#fff",
        "o": "#ff0",
        // "P": "#0f0",
        // "a": "#000",
        // "b": "#000",
        // "c": "#000",
        // "d": "#000"
    };

    for (let i = 0; i < board.length; i++) {
        for (let j = 0; j < board[i].length; j++) {
            const char = board[i][j];
            ctx.fillStyle = mapping[char] || "#fff"; 
            ctx.fillRect(j * PixelWidth, i * PixelHeight, PixelWidth+0.4, PixelHeight+0.4);
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


    console.log("Canvas updated with new board state");
}


const HandleWsMessage = (message) =>{

    const LobbyID = window.location.href.split("/").pop();

    const data = JSON.parse(message)
    console.log("Handling WebSocket message:", data);

    if (data.type === "UserInfoUpdate"){
        UpdateUsers(data.users);
    }

    if (data.type === "HostUpdate"){
        console.log("Host update received:", data.host);
        IsHost = true;
        ToggleHost();
    }

    if (data.type === "StartAlert"){
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

    if (data.type === "BoardUpdate"){
        pacman.targetX = data.pacman.target_y;
        pacman.targetY = data.pacman.target_x;
        pacman.dir = data.pacman.dir;

          
        Ghosts = data.enemy.map(ghost => ({
            targetX: ghost.target_y,
            targetY: ghost.target_x,
            x: Ghosts[Ghosts.findIndex(g => g.name === ghost.name)].x,
            y: Ghosts[Ghosts.findIndex(g => g.name === ghost.name)].y,
            dir: ghost.dir,
            name: ghost.name
            
        }))
        // console.log("Board update received:", data.board);
         UpdateCanvas(boardTemp); 
        animateEntities(data.board);
    }
}

function animateEntities(boardTemp) {
        if (
            Math.abs(pacman.x - pacman.targetX) < 0.1 &&
            Ghosts.some(ghost => Math.abs(ghost.x - ghost.targetX) < 0.1) &&
            Math.abs(pacman.y - pacman.targetY) < 0.1 &&
            Ghosts.some(ghost => Math.abs(ghost.y - ghost.targetY) < 0.1)
        ) {return}



        UpdateCanvas(boardTemp);
        ctx = Canvas.getContext("2d");

        //1 cell/1 sec -> 0.004 animation speed
        //1 cell/0.5 sec -> 0.008 animation speed
        const speed = 0.008

        

        function moveEntity(entity){

            if (
                Math.abs(entity.x - entity.targetX) < 0.1 
            ) {
                entity.x = entity.targetX;
                return;
            }

            if (
                Math.abs(entity.y - entity.targetY) < 0.1
            ) {
                entity.y = entity.targetY; }

            DirX = 0
            DirY = 0

            if(entity.dir === 0){
                DirX = 0
                DirY = 1;
            }else   if(entity.dir === 1){
                DirX= 0
                DirY = -1;
            }else if(entity.dir === 2){
                DirX = -1
                DirY = 0;
            }else if (entity.dir === 3){
                DirX = 1
                DirY = 0;
            }
                
            entity.x += DirX * PixelWidth * speed;
            entity.y += DirY * PixelHeight * speed;
        }

        moveEntity(pacman);
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
       

        console.log("animating",pacman.x,pacman.targetX)

        if (
            Math.abs(pacman.x - pacman.targetX) > 0.1 ||
            Ghosts.some(ghost => Math.abs(ghost.x - ghost.targetX) > 0.1)
        ){
            // UpdateCanvas(boardTemp);
            
            setTimeout(() => animateEntities(boardTemp), 1000 / 20);

            
            // animateEntities(boardTemp)
        }else{
            console.log("Animation complete, redrawing canvas");
        }
}

const UpdateUsers = (users) =>{
    const PlayerContainer = document.getElementById("playersContainer");
    PlayerContainer.innerHTML = "";

    users.forEach(e => {
        const div = document.createElement("div");
        div.className = "PlayerItem";
        const Title = document.createElement("h1");
        Title.className = "PlayerTitle"
        host = e.host ? "*" : "";
        Title.innerHTML = `${host}${e.name}<img src="/static/images/ghost.png" alt="Ghost Icon" class="ghost-icon">`
        const Score = document.createElement("h2")
        Score.className = "PlayerScore";
        Score.innerHTML = ` Score: <span class="yellow">0</span>`;
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
            ctx.fillRect(j * PixelWidth, i * PixelHeight, PixelWidth+0.4, PixelHeight+0.4);
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

    console.log(Ghosts)
    Ghosts.forEach(g => {
        console.log("Drawing ghost:", g.name, "at position:", g.x, g.y);
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