document.getElementById("connectForm").onsubmit = () => {
    pwd = document.getElementById("passwordField").value;
    ip = document.getElementById("ipField").value;

    console.log("pwd: " + pwd);
    console.log("ip: " + ip);

    connectToWS(ip, pwd);

    return false;
}

function startProgram(program){
    ws.send("run " + program)
}


let ws;

function connectToWS(ip, pwd){
    ws = new WebSocket("ws://"+ip+":1337/codedeck");

    ws.addEventListener("open", (event) => {
        ws.send("auth " + pwd);
    
        window.alert("Connected to server")
        console.log("connected to server");
    });
    
    ws.addEventListener("message", (event) => {
        console.log("message from server: " + event.data);
    });
}