export class Constants {
    static get frameRate() {
        return 60;
    }

    static get deltaTime() {
        return 1 / Constants.frameRate;
    }

    static get commands() {
        return {
            join: "Join",
            leave: "Leave",
            sync: "Sync",
            summon: "Summon"
        }
    }
}