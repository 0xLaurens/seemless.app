export interface ToastSettings {
  message: string
  type: ToastType
}

export interface Toast extends ToastSettings {
  id: string
  timeoutId?: ReturnType<typeof setTimeout>
}

export enum ToastType {
  Info = 'alert-info',
  Success = 'alert-success',
  Warning = 'alert-warning',
  Error = 'alert-error'
}
