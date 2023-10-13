import { defineStore } from 'pinia'
import type { Ref } from 'vue'
import { ref } from 'vue'
import type { Toast, ToastSettings } from '@/models/toast'

function randomUUID(): string {
  const random = Math.random()
  return Number(random).toString(32)
}

export const useToastStore = defineStore('toast', () => {
  const toasts: Ref<Toast[]> = ref([])

  function close(id: string) {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }

  function notify(settings: ToastSettings) {
    const id: string = randomUUID()
    const toast: Toast = { ...settings, id }
    setTimeout(close, 5000, id)
    toasts.value.push(toast)
  }

  return {
    toasts,
    notify,
    close
  }
})
