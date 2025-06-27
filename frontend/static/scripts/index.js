const getLobbies = async () =>{
    const response = await fetch("/lobbies");
    if (!response.ok) {
        throw new Error("Network response was not ok");
    }
    return await response.json();
}

var mode
var JoiningId


document.addEventListener("DOMContentLoaded", async () => {


    const input = document.getElementById('playerName');
    const overlay = document.getElementById('retroInputOverlay');

    const createLobbyButton = document.getElementById("createLobby");
    const enterGameButton = document.getElementById("enterGame");
    const refreshButton = document.getElementById("refreshLobby");

    refreshButton.addEventListener("click", async () => {
        await UpdateLobbyList();
    })


    function updateOverlay() {
    const value = input.value;

    if (value.length === 0) {
        button = document.getElementById("enterGame")
        button.disabled = true;
    }else {
        button = document.getElementById("enterGame")
        button.disabled = false; 
    }

    if (value.length > 5){
        input.value = value.slice(0, 5); 
        overlay.innerHTML = input.value;  
    }
    const cursor = '<span class="blink-cursor">_</span>';
    overlay.innerHTML = value + cursor;
    }

    input.addEventListener('input', updateOverlay);
    input.addEventListener('focus', updateOverlay);
    input.addEventListener('blur', () => {
        overlay.innerHTML = input.value; 
    });
    

    const joinButtons = await UpdateLobbyList();
    if (joinButtons === undefined) {
        console.error("No join buttons found. Lobby list might be empty.");
    }else{
         joinButtons.forEach(button => {
        button.addEventListener("click",async (event) => {
            LobbyContainer = document.getElementById("lobbyContainer");
            LobbyContainer.style.display = "None";
            
            JoiningId = event.target.getAttribute("data-lobby-id");
            mode = "join";

            Namefield = document.getElementById("nameField");
            Namefield.style.display = "flex";
            document.getElementById("playerName").focus()

            updateOverlay();
            console.log("updating overlay")

            
        })
    })
    }
   

    createLobbyButton.addEventListener("click",async ()=>{
        LobbyContainer = document.getElementById("lobbyContainer");
        LobbyContainer.style.display = "None";

        mode = "create";
        Namefield = document.getElementById("nameField");
        Namefield.style.display = "flex";

        updateOverlay();
    })

    enterGameButton.addEventListener("click", async () => {
        console.log("Entering game with mode:", mode);
        if (mode === "create") {
            await createLobby(input.value);
        } 
        if (mode === "join"){
            window.location.href = `/lobby/${JoiningId}/${input.value}`;
        }
    })
    document.addEventListener("keydown", (event) => {
        if (event.key === "Enter") {
            enterGameButton.click();
        }
    })
}
)

const createLobby = async (name) =>{
     const response = await fetch("/create")
        if (!response.ok) {
            throw new Error("Network response was not ok");
        }
        const data = await response.json();
        console.log("Created Lobby:", data);
        window.location.href = `/lobby/${data.ID}/${name}`;   
}

const UpdateLobbyList = async () => {
    const lobbies = await getLobbies();
    const LobbyList = document.getElementById("lobby-list");
    LobbyList.innerHTML = ""; 

    if (lobbies.length === 0) {
        const noLobbiesMessage = document.createElement("div");
        noLobbiesMessage.className = "NoLobbies";
        noLobbiesMessage.innerHTML = "<p>No lobbies available. Create one to start playing!</p>";
        noLobbiesMessage.style.textAlign = "center";
        LobbyList.appendChild(noLobbiesMessage);
        return;
    }

    lobbies.forEach(e => {
        const lobbyItem = document.createElement("div");
        lobbyItem.className = "LobbyItem";
        lobbyItem.innerHTML = `
            <h2 class="LobbyName">${e.name}</h2>
            <p>Players: ${e.players}/5</p>
            <button class="join-button" data-lobby-id="${e.id}">Join</button>`;
        LobbyList.appendChild(lobbyItem);
    });

    const joinButtons = document.querySelectorAll(".join-button");
    return joinButtons;
    
}
    
