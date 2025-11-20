const params = new URLSearchParams(window.location.search);
const room = params.get("room");
const username = params.get("username");

if (!room || !username) {
    alert("Missing room or username. Redirecting...");
    window.location.href = "/";
}

const socket = new WebSocket(`ws://${location.host}/room?room=${room}&username=${username}`);

socket.onmessage = (event) => {
    try {
        const data = JSON.parse(event.data);

        console.log("RAW:", event.data);

        if (data.type === "users") {
            updateUserList(data);
            return;
        }

        renderMessage(data);
    } catch (err) {
        console.error("Invalid JSON received:", event.data);
    }
};


function sendMessage() {
    const input = document.getElementById("msg");
    if (input.value.trim() !== "") {
        socket.send(input.value);
        input.value = "";
    }
}

document.getElementById("sendBtn").addEventListener("click", sendMessage);

document.getElementById("msg").addEventListener("keyup", function (event) {
    if (event.key === "Enter") {
        sendMessage();
    }
});

//===========================================================
function updateUserList(data) {
    const listDiv = document.getElementById("user-list");
    const countDiv = document.getElementById("user-count");

    listDiv.innerHTML = "";

    data.users.forEach(u => {
        const div = document.createElement("div");
        div.textContent = u;
        listDiv.appendChild(div);
    });

    countDiv.textContent = "Total: " + data.count;
}

function renderMessage(data) {
    const msgContainer = document.createElement("div");
    msgContainer.classList.add("message-container");

    const usernameDiv = document.createElement("div");
    usernameDiv.classList.add("username");
    usernameDiv.textContent = data.name;

    const messageDiv = document.createElement("div");
    messageDiv.classList.add("message");
    messageDiv.textContent = data.message;

    msgContainer.appendChild(usernameDiv);
    msgContainer.appendChild(messageDiv);
    document.getElementById("messages").appendChild(msgContainer);

    const messagesDiv = document.getElementById("messages");
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}
