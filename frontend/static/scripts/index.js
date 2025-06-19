const getLobbies = async () =>{
    const response = await fetch("http://localhost:8080/lobbies");
    if (!response.ok) {
        throw new Error("Network response was not ok");
    }
    return await response.json();
}



document.addEventListener("DOMContentLoaded", async () => {

   

    const lobbies = await getLobbies()
    console.log(lobbies);
    
    const LobbyList = document.getElementById("lobby-list");
    // lobbies.map(e => {
    //     const lobbyItem = document.createElement("div");
    //     lobbyItem.className = "LobbyItem";
    //     lobbyItem.innerHTML = `
    //         <h2>${e.name}</h2>
    //         <p>Players: ${e.players.length}/5</p>
    //         <button class="join-button" data-lobby-id="${e.id}">Join</button>`
    //     LobbyList.appendChild(lobbyItem);
    // });

    if (lobbies.length === 0) {
        const noLobbiesMessage = document.createElement("div");
        noLobbiesMessage.className = "NoLobbies";
        noLobbiesMessage.innerHTML = "<p>No lobbies available. Create one to start playing!</p>";
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
}
  
)