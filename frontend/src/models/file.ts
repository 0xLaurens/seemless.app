export interface FileMessage {
  status: FileStatus
  name: string
  mime: string
}

export enum FileStatus {
  Complete = 'Complete',
  Busy = 'Busy'
}
