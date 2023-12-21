import {defineStore} from "pinia";
import {ref} from "vue";
import type {Setting} from "@/models/settings";

export const useSettingsStore = defineStore('settings', () => {
    const defaultSettings: Setting[] = [

        {
            name: "saveUsername",
            description: "Remember your custom username for next time.",
            value: true,
        },
        {
            name: "autoTransfer",
            description: "Auto accept all incoming transfer requests",
            value: false,
        }
    ]
    const settings = ref(defaultSettings)
    loadSettings()

    function saveSettings(settings: Setting[]) {
        localStorage.removeItem("settings")
        localStorage.setItem("settings", JSON.stringify(settings))
    }

    function loadSettings() {
        const storedSettings = localStorage.getItem("settings")
        console.log(storedSettings)
        if (storedSettings) {
            settings.value = JSON.parse(storedSettings)
        }
    }

    function getSettings(): Setting[] {
        return settings.value
    }

    return {
        saveSettings,
        loadSettings,
        getSettings,
    }
})
