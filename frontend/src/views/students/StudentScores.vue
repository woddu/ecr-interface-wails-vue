<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useHighestScoresStore } from '../../store/highestScoresStore';
import { useStudentStore } from '../../store/studentStore';
import { EventsOn } from '../../../wailsjs/runtime/runtime';

interface GradeRange {
  min: number,
  max: number,
  transmuted: number,
}

const gradeList: GradeRange[] = [
  { min: 100.0, max: 100.0, transmuted: 100 },
  { min: 98.40, max: 99.99, transmuted: 99 }, 
  { min: 96.80, max: 98.39, transmuted: 98 }, 
  { min: 95.20, max: 96.79, transmuted: 97 }, 
  { min: 93.60, max: 95.19, transmuted: 96 }, 
  { min: 92.00, max: 93.59, transmuted: 95 }, 
  { min: 90.40, max: 91.99, transmuted: 94 }, 
  { min: 88.80, max: 90.39, transmuted: 93 }, 
  { min: 87.20, max: 88.79, transmuted: 92 }, 
  { min: 85.60, max: 87.19, transmuted: 91 }, 
  { min: 84.00, max: 85.59, transmuted: 90 }, 
  { min: 82.40, max: 83.99, transmuted: 89 }, 
  { min: 80.80, max: 82.39, transmuted: 88 }, 
  { min: 79.20, max: 80.79, transmuted: 87 }, 
  { min: 77.60, max: 79.19, transmuted: 86 }, 
  { min: 76.00, max: 77.59, transmuted: 85 }, 
  { min: 74.40, max: 75.99, transmuted: 84 }, 
  { min: 72.80, max: 74.39, transmuted: 83 }, 
  { min: 71.20, max: 72.79, transmuted: 82 }, 
  { min: 69.60, max: 71.19, transmuted: 81 }, 
  { min: 68.00, max: 69.59, transmuted: 80 }, 
  { min: 66.40, max: 67.99, transmuted: 79 }, 
  { min: 64.80, max: 66.39, transmuted: 78 }, 
  { min: 63.20, max: 64.79, transmuted: 77 }, 
  { min: 61.60, max: 63.19, transmuted: 76 }, 
  { min: 60.00, max: 61.59, transmuted: 75 }, 
  { min: 56.00, max: 59.99, transmuted: 74 }, 
  { min: 52.00, max: 55.99, transmuted: 73 }, 
  { min: 48.00, max: 51.99, transmuted: 72 }, 
  { min: 44.00, max: 47.99, transmuted: 71 }, 
  { min: 40.00, max: 43.99, transmuted: 70 }, 
  { min: 36.00, max: 39.99, transmuted: 69 }, 
  { min: 32.00, max: 35.99, transmuted: 68 }, 
  { min: 28.00, max: 31.99, transmuted: 67 }, 
  { min: 24.00, max: 27.99, transmuted: 66 }, 
  { min: 20.00, max: 23.99, transmuted: 65 }, 
  { min: 16.00, max: 19.99, transmuted: 64 }, 
  { min: 12.00, max: 15.99, transmuted: 63 }, 
  { min: 8.00, max:  11.99, transmuted: 62 }, 
  { min: 4.00, max:  7.99, transmuted:  61 }, 
  { min: 0.00, max:  3.99, transmuted:  60 }, 
]

const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

const highestScoresStore = useHighestScoresStore()

const studentStore = useStudentStore()

const writtenWorks = ref<number[]>([])
const performanceTasks = ref<number[]>([])
const examScore = ref<number>(0)

const grade = computed<number>(() => {
  var wwHighestTotal: number = 0
  var wwTotal: number = 0
  highestScoresStore.wwHighestScores.forEach((score, index)  => {
      if (score !== 0){
        wwHighestTotal += score
        wwTotal += studentStore.writtenWorks[index]
      }
  });
  var ptHighestTotal: number = 0
  var ptTotal: number = 0
  highestScoresStore.ptHighestScores.forEach((score, index) => {
      if (score !== 0){
        ptHighestTotal += score
        ptTotal += studentStore.performanceTasks[index]
      }
  });
  
  const wwWeightedScore: number = ( wwTotal / wwHighestTotal ) * highestScoresStore.weightedScores[0] * 100
  const ptWeightedScore: number = ( ptTotal / ptHighestTotal ) * highestScoresStore.weightedScores[1] * 100
  const examWeightedScore: number = ( studentStore.exam / highestScoresStore.examHighestScore ) * highestScoresStore.weightedScores[2] * 100
  const initalGrade: number = wwWeightedScore + ptWeightedScore + examWeightedScore
  var transmuted: number = 0
  gradeList.filter((range: GradeRange) => {
    if (initalGrade >= range.min, initalGrade <= range.max){
      transmuted = range.transmuted
    }
  })
  
  return transmuted
})

const hasChangesWW = computed<boolean>(() => {
  return studentStore.writtenWorks.some(
    (score, i) => score !== writtenWorks.value[i]
  )
})

const hasChangesPT = computed<boolean>(() => {
    return studentStore.performanceTasks.some(
      (score, i) => score !== performanceTasks.value[i]
    )
})

async function saveChangesWW(){
  emit('loading', true)
  const { EditStudentScores } = await import('../../../wailsjs/go/main/App')
  await EditStudentScores(studentStore.row, writtenWorks.value, true)
}

async function saveChangesPT() {
  emit('loading', true)
  const { EditStudentScores } = await import('../../../wailsjs/go/main/App')
  await EditStudentScores(studentStore.row, performanceTasks.value, false)
}

async function saveChangesExam() {
  emit('loading', true)
  const { EditStudentExamScore } = await import('../../../wailsjs/go/main/App')
  await EditStudentExamScore(studentStore.row, examScore.value)
}

onMounted(() => {
  writtenWorks.value = [... studentStore.writtenWorks]
  performanceTasks.value = [... studentStore.performanceTasks]
  examScore.value = studentStore.exam

  EventsOn('excel:done_editing_student_scores', (isWrittenWorks: boolean, newScores: number[]) => {
    if (isWrittenWorks){
      studentStore.setWrittenWorks(newScores)
      writtenWorks.value = [... studentStore.writtenWorks]
    } else {
      studentStore.setPerformanceTasks(newScores)
      performanceTasks.value = [... studentStore.performanceTasks]
    }
    emit('loading', false)
  })
  EventsOn('excel:done_editing_student_exam_score', (newExamScore: number) => {
    studentStore.setExam(newExamScore)
    examScore.value = studentStore.exam
    emit('loading', false)
  })
})

</script>

<template>
  <div class="h-full w-full flex flex-col items-center justify-center p-10">
    <div class="flex flex-col">
      <h1 class="text-3xl font-bold mb-4">{{ studentStore.name }}    Grade: {{ grade }}</h1>
      <div class="grid grid-cols-6 my-auto">
        <div class="flex flex-col gap-8 col-span-5">
          <div class="grid grid-cols-8">
            <div  class="col-span-7">
              <h1 class="text-2xl font-medium">Written Works - {{ highestScoresStore.weightedScores[0] * 100 }}%</h1>
              <div class="grid grid-cols-5 gap-4">
                <ScoresField v-for="( ,index) in highestScoresStore.wwHighestScores.filter(e => e !== 0)" :key="index" v-model="writtenWorks[index]"
                  :label="`Total: ${highestScoresStore.wwHighestScores[index]}`" :max="highestScoresStore.wwHighestScores[index]" size="lg" />
              </div>
            </div>
            <div class="h-full flex">
              <UButton class="mt-auto mx-auto" label="Save" size="xl" :variant="!hasChangesWW ? 'ghost' : 'solid'"
                :disabled="!hasChangesWW" @click="saveChangesWW" />
            </div>
          </div>
          <div class="grid grid-cols-8">
            <div  class="col-span-7">
              <h1 class="text-2xl font-medium">Performance Tasks - {{ highestScoresStore.weightedScores[1] * 100 }}%</h1>
              <div class="grid grid-cols-5 gap-4">
                <ScoresField v-for="( ,index) in highestScoresStore.ptHighestScores.filter(e => e !== 0)" :key="index" v-model="performanceTasks[index]"
                  :label="`Total: ${highestScoresStore.ptHighestScores[index]}`" :max="highestScoresStore.ptHighestScores[index]" size="lg" />
              </div>
            </div>
            <div class="h-full flex">
              <UButton class="mt-auto mx-auto" label="Save" size="xl" :variant="!hasChangesPT ? 'ghost' : 'solid'"
                :disabled="!hasChangesPT" @click="saveChangesPT" />
            </div>
          </div>
        </div>
        <div class="flex p-8">
          <div  class="m-auto">
            <h1 class="text-2xl font-medium mb-4">Exam - {{ highestScoresStore.weightedScores[2] * 100 }}%</h1>
            <ScoresField v-model="examScore" :label="`Total: ${highestScoresStore.examHighestScore.toString()}`" :max="highestScoresStore.examHighestScore" size="xl" />
            <UButton class="mt-4" label="Save" size="xl"
              :variant="examScore === studentStore.exam ? 'ghost' : 'solid'"
              :disabled="examScore === studentStore.exam" @click="saveChangesExam" />
          </div>
        </div>
      </div>
        </div>
    </div>

</template>