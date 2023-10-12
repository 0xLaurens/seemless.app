export interface DatachannelMessage {
  status: DcStatus
  username: string
}

export enum DcStatus {
  ClientHello = 'ClientHello',
  ClientClose = 'ClientClose'
}
