import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Ref } from 'vue'

export const useRetryStore = defineStore('retry', () => {
  const DELAY: number = 1000
  const MAX_TIMEOUT: number = 8000 //MAX MS between reload

  const attempts = ref(0)
  const timer: Ref<number | undefined> = ref()

  function start(fn: () => void) {
    if (isActive()) clearTimeout(timer.value)

    let backOffDuration = MAX_TIMEOUT
    if (attempts.value < 3) {
      backOffDuration = Math.pow(2, attempts.value) * DELAY
    }
    console.log('BACKOFF', backOffDuration, 'MS')
    timer.value = setTimeout(() => {
      fn()
      attempts.value++
      start(fn)
    }, backOffDuration)
  }

  function stop() {
    if (timer.value) {
      clearTimeout(timer.value)
    }
    attempts.value = 0
  }

  function isActive() {
    return timer.value != undefined
  }

  return { start, stop, isActive }
})
