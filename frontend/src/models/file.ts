export interface FileMessage {
    status: FileStatus
    name: string
    mime: string
    size: number
    from: string
}

export enum FileStatus {
    Offer = 'Offer',
    Complete = 'Complete',
    Busy = 'Busy'
}
