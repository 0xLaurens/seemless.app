import {defineStore} from "pinia";
import type {Ref} from "vue";
import {ref} from "vue";
import type {Room, RoomCode} from "@/models/room";

export const useRoomStore = defineStore('room', () => {
    // const publicRoom: Ref<Room |undefined> = ref()
    const publicRoomCode: Ref<RoomCode |undefined> = ref()

    function setRoomCode(code: RoomCode) {
        publicRoomCode.value = code
    }

    function getRoomCode(): RoomCode | undefined {
        return publicRoomCode.value
    }


    return {
        setRoomCode,
        getRoomCode
    }
})
