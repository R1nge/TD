import {getPlayerEntities} from "./engine.js";

const canvas = document.getElementById("myCanvas");
export const ctx = canvas.getContext("2d");

export function render(deltaTime) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    //requestAnimationFrame(renderPlayers);
    renderPlayers();
}

function renderPlayers() {
    if (getPlayerEntities().size === 0) {
        return;
    }

    for (const playerEntity of getPlayerEntities().values()) {
        if (playerEntity === undefined || playerEntity.id === undefined) {
            continue;
        }
        
        playerEntity.draw(ctx);
    }
}