const ws = new WebSocket("ws://192.168.178.40:1337/codedeck");

ws.addEventListener("open", (event) => {
    ws.send("auth abc");

    console.log("connected to server");
});

ws.addEventListener("message", (event) => {
    console.log("message from server: " + event.data);
});


function startProgram(program){
    ws.send("run " + program)
}