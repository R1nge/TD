import {Constants} from "./constants.js";
import {localId, sendToServer} from "./server.js";
import {PlayerEntity} from "./playerEntity.js";
import {ctx, render} from "./renderer.js";
import {Utils} from "./utils.js";

window.addEventListener('click', function (event) {

    const player = {
        id: localId,
        mousePositionX: event.clientX,
        mousePositionY: event.clientY
    }
    sendToServer(player, Constants.commands.shoot);
});

function gameLoop() {
    render(Constants.deltaTime);
}

setInterval(gameLoop, Constants.deltaTime);