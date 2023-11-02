export interface FileMessage {
    status: FileStatus
    name: string
    mime: string
    size: number
}

export interface FileOffer {
    id: string
    status: FileSetup
    from: string
    files: FileMessage[]
}

export enum FileSetup {
    Offer = 'Offer',
    Accept = 'Accept',
    Deny = 'Deny',
}

export enum FileStatus {
    Init = 'Init',
    Complete = 'Complete',
    Busy = 'Busy'
}
