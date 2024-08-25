import {createPlayer, getPlayerEntities, moveCallback, removePlayer, shoot} from "./engine.js";
import {Constants} from "./constants.js";

let socket = new WebSocket("ws://localhost:8080");

export let localId = 0;

socket.onopen = function (e) {
    console.log(`Connected to server`);

    const loginData = {
        name: "R1nge"
    }

    console.log(`Login data sent: ${loginData.id} ${loginData.x} ${loginData.y} ${loginData.name}`);

    sendToServer(loginData, Constants.commands.join);
};

socket.onmessage = function (event) {
    //trim before space
    const data = event.data.substring(event.data.indexOf(" "));
    console.log(`Message received: ${data}`);
    const parsedData = JSON.parse(data);

    if (event.data.startsWith(Constants.commands.join)) {
        console.log("Join message received: " + parsedData.id);
        
        if (localId === 0) {
            localId = parsedData.id;
        }
        
        createPlayer(parsedData.id);
        return;
    }

    if (event.data.startsWith(Constants.commands.create)) {
        console.log(`Create message received: ${parsedData.id}`);
        return;
    }

    if (event.data.startsWith(Constants.commands.sync)) {
        for (let i = 0; i < parsedData.length; i++) {
            console.log(`Sync message received: ${parsedData[i]}`);
            console.log(`Player ${parsedData[i].id} ${parsedData[i].x} ${parsedData[i].y}`);
            createPlayer(parsedData[i].id);
            const player = getPlayerEntities().get(parsedData[i].id);
            if (!player) {
                console.log(`Player ${parsedData[i].id} not found`);
                return;
            }
            player.x = parsedData[i].x;
            player.y = parsedData[i].y;
            player.rotationAngle = parsedData[i].rotationAngle;
        }

        return;
    }

    if (event.data.startsWith(Constants.commands.shoot)) {
        console.log(`Shoot message received: ${parsedData.id}`);
        shoot();
        return;
    }

    if (event.data.startsWith(Constants.commands.leave)) {
        console.log(`Leave message received: ${parsedData.id}`);
        removePlayer(parsedData.id);
        return;
    }

    console.log(`received a message: ${event.data}`);
}

export function sendToServer(dataStruct, messageType) {
    const json = JSON.stringify(dataStruct)
    console.log(`Send to server json sent: ${json}`);
    socket.send(messageType + json);
}