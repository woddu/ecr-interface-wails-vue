import { defineStore } from "pinia";
import { ref } from "vue";

export const useStudentStore = defineStore('student', () => {
    const row = ref<number>(0)
    const name = ref<string>('')
    const exam = ref<number>(0)
    const writtenWorks = ref<number[]>([])
    const performanceTasks = ref<number[]>([])

    function setRow(newRow: number){
        row.value = newRow
    }

    function setName(newName: string){
        name.value = newName
    }

    function setExam(newExam: number){
        exam.value = newExam
    }

    function setWrittenWorks(newWrittenWorks: number[]){
        writtenWorks.value = newWrittenWorks
    }

    function setPerformanceTasks(newPerformanceTasks: number[]){
        performanceTasks.value = newPerformanceTasks
    }

    return {
        row,
        name,
        exam,
        writtenWorks,
        performanceTasks,
        setRow,
        setName,
        setExam,
        setWrittenWorks,
        setPerformanceTasks
    }
})