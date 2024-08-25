export class Utils{
    
    static lerp(a, b, t){
        return a + (b - a) * t
    }
    
    static randomInt(min, max){
        return Math.floor(Math.random() * (max - min + 1) + min)
    }
}