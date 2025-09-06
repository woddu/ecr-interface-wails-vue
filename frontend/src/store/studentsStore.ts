import { defineStore } from "pinia";
import { ref } from "vue"

export const studentsStore = defineStore("students", () => {
    const males = ref<string[]>([]);
    function setMales(newMales: string[]) {
        males.value = newMales;
    }

    const females = ref<string[]>([]);
    function setFemales(newFemales: string[]) {
        females.value = newFemales;
    }

    return {
        males,
        setMales,
        females,
        setFemales
    }
})