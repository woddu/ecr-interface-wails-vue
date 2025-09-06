import { defineStore } from "pinia";
import { ref } from 'vue'

export const fileViewStore = defineStore("fileView", () => {
    const fileName = ref<string>("");
    function setFileName(name: string) {
        fileName.value = name;
    }
    const tracks = ref<string[]>([]);
    function setTracks(newTracks: string[]) {
        tracks.value = newTracks;
    }
    const track = ref<string>("");
    function setTrack(newTrack: string) {
        track.value = newTrack;
    }
    const firstSem = ref<boolean>(true);
    function setFirstSem(isFirstSem: boolean) {
        firstSem.value = isFirstSem;
    }
    const doneReading = ref<boolean>(false);
    function setDoneReading(isDone: boolean) {
        doneReading.value = isDone;
    }
    return { fileName, track, firstSem, doneReading, tracks, setFileName, setTrack, setFirstSem, setDoneReading, setTracks };
})