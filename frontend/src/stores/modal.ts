import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useModalStore = defineStore('modal', () => {
  const title = ref('')
  const message = ref('')
  const isActive = ref(false)
  const btnText = ref('')

  function getTitle(): string {
    return title.value
  }

  function getMessage(): string {
    return message.value
  }

  function getBtnText(): string {
    return btnText.value
  }

  function getActiveState(): boolean {
    return isActive.value
  }

  function trigger(settings: { title: string; message: string; btnText?: string }) {
    title.value = settings.title
    message.value = settings.message

    btnText.value = settings.btnText != null ? settings.btnText : 'Accept'

    isActive.value = true
  }

  function close() {
    isActive.value = false
  }

  return {
    getTitle,
    getMessage,
    getActiveState,
    getBtnText,
    trigger,
    close
  }
})
