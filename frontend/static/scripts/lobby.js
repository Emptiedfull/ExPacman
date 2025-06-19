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

console.log(boardTemp.length,boardTemp[0].length);

document.addEventListener("DOMContentLoaded", () => {
    const GameContainer = document.getElementById("gameContainer");
    GameContainer.style.display = "None"

    const LobbyID = window.location.href.split("/").pop();
    console.log("Lobby ID:", LobbyID);

    const WsURL = `ws://localhost:8080/ws/${LobbyID}`;
    console.log("WebSocket URL:", WsURL);

    const ws = new WebSocket(WsURL);

    ws.onopen = () =>{
        const NameResponse = {
            "id": LobbyID,
            "name": "ExamplePlayer",
        }
        JSON.stringify(NameResponse);
        ws.send(JSON.stringify(NameResponse));
    }

    ws.onmessage = (event) => {
        event.preventDefault();
        HandleWsMessage(event.data);
    }



});


const HandleWsMessage = (message) =>{

    const LobbyID = window.location.href.split("/").pop();

    const data = JSON.parse(message)
    console.log("Handling WebSocket message:", data);


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
        "p": "#0f0",
        "a": "#000",
        "b": "#000",
        "c": "#000",
        "d": "#000"
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


    // console.log("Canvas loaded with dimensions:", Width, Height,PixelHeight,PixelWidth);

}