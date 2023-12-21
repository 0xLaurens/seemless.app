export interface User {
    id: string
    username: string
    device: string
}

export enum UserTransfer {
    Response = "Awaiting response...",
    Transfer = "Transferring files...",
    Sent = "Sent"
}
