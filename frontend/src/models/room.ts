export type RoomCode = string

export interface Room {
    id: string,
    roomType: RoomType
    code: RoomCode,
}


export enum RoomType {
    PublicRoom = "publicRoom",
    PrivateRoom = "localRoom"
}