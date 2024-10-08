export interface FileMessage {
    status: FileStatus
    name: string
    mime: string
    progress: number
    size: number
}

export interface FileOffer {
    id: string
    status: FileSetup
    from: string
    target: string
    files: FileMessage[]
    current: number
}

export enum FileSetup {
    Offer = 'Offer', // status when files are being offered
    AcceptOffer = 'AcceptOffer', // status when offer is accepted
    DenyOffer = 'DenyOffer', // status when offer is denied
    DownloadProgress = 'DownloadProgress', // status for updating download status
    LatestOffer = "LatestOffer", // status for setting the latest offer
    RequestNext = 'RequestNext', // status receiver is ready for the next file is done
    Complete = "Complete", // status when the offer is finished
}

export enum FileStatus {
    Init = 'Init',
    Complete = 'Complete',
    Busy = 'Busy'
}
