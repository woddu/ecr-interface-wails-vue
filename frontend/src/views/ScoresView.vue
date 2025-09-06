<script setup lang="ts">
import ScoresField from '../components/ScoresField.vue'
import { computed, onMounted, ref } from 'vue'
import { useHighestScoresStore } from '../store/highestScoresStore';
import { EventsOn } from '../../wailsjs/runtime/runtime';

const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

const highestScoresStore = useHighestScoresStore();

const wwHighestScores = ref<number[]>([])
const ptHighestScores = ref<number[]>([])
const examHighestScore = ref<number>(0)

const isFetching = ref<boolean>(false)

const hasChangesWW = computed(() => {
  return highestScoresStore.wwHighestScores.some(
    (score, i) => score !== wwHighestScores.value[i]
  )
})

const hasChangesPT = computed(() => {
    return highestScoresStore.ptHighestScores.some(
      (score, i) => score !== ptHighestScores.value[i]
    )
})

async function saveChangesWW() {
    emit('loading', true);
    const { EditHighestScores } = await import('../../wailsjs/go/main/App');
    await EditHighestScores(wwHighestScores.value, true);
}

async function saveChangesPT() {
    emit('loading', true);
    const { EditHighestScores } = await import('../../wailsjs/go/main/App');
    await EditHighestScores(ptHighestScores.value, false);
}

async function saveChangesExam(){
    emit('loading', true);
    const { EditExamHighestScore } = await import('../../wailsjs/go/main/App');
    await EditExamHighestScore(examHighestScore.value);
}

onMounted(() => {
    isFetching.value = true
    wwHighestScores.value = [...highestScoresStore.wwHighestScores]
    ptHighestScores.value = [...highestScoresStore.ptHighestScores]
    examHighestScore.value = highestScoresStore.examHighestScore
    isFetching.value = false

    EventsOn('excel:done_editing_highest_scores', (writtenWorks: boolean, scores: number[]) => {
        if (writtenWorks){
            highestScoresStore.setWwHighestScores(scores)
        } else {
            highestScoresStore.setPtHighestScores(scores)
        }
        emit('loading', false)
    })

    EventsOn('done_editing_exam_highest_score', (score: number) => {
        highestScoresStore.setExamHighestScore(score)
        emit('loading', false)
    })
})

</script>

<template>
    <div class="h-full w-full flex flex-col items-center justify-center p-10">
        <div class="grid grid-cols-6 my-auto">
            

            <div class="flex flex-col gap-8 col-span-5">
                <div class="grid grid-cols-8">
                    <div v-if="!isFetching" class="col-span-7">
                        <h1 class="text-xl font-medium">Written Works - {{ highestScoresStore.weightedScores[0] * 100 }}%</h1>
                        <div class="grid grid-cols-5 gap-4">
                            <ScoresField v-for="(, index) in wwHighestScores" :key="index"
                                v-model="wwHighestScores[index]" :label="`#${index + 1}`" 
                                size="lg"/>
                        </div>
                    </div>
                    <div v-else class="w-[635px] h-[164px] col-span-7">
                        <USkeleton class="h-6 w-45 mb-4" />                        
                        <div class="grid grid-cols-5 gap-4">    
                            <div v-for="i in 10" :key="i" class="flex flex-col">
                                <USkeleton class="h-5 w-10 mb-2" />
                                <USkeleton class="h-5 w-30" />
                            </div>
                        </div>
                    </div>
                    <div class="h-full flex">
                        <UButton
                          class="mt-auto mx-auto"
                          label="Save"
                          size="xl"
                          :variant="!hasChangesWW ? 'ghost' : 'solid'"
                          :disabled="!hasChangesWW" 
                          @click="saveChangesWW"/>
                    </div>
                </div>
                <div class="grid grid-cols-8">
                    <div v-if="!isFetching" class="col-span-7">
                        <h1 class="text-xl font-medium">Performance Tasks - {{ highestScoresStore.weightedScores[1] * 100 }}%</h1>
                        <div class="grid grid-cols-5 gap-4">
                            <ScoresField v-for="(, index) in ptHighestScores" :key="index"
                                v-model="ptHighestScores[index]" :label="`#${index + 1}`" 
                                size="lg"/>
                        </div>
                    </div>
                    <div v-else class="w-[635px] h-[164px]  col-span-7">
                        <USkeleton class="h-6 w-45 mb-4" />
                        <div class="grid grid-cols-5 gap-4">                    
                            <div v-for="i in 10" :key="i" class="flex flex-col">
                                <USkeleton class="h-5 w-10 mb-2" />
                                <USkeleton class="h-5 w-30" />
                            </div>
                        </div>
                    </div>
                    <div class="h-full flex">
                        <UButton
                          class="mt-auto mx-auto"
                          label="Save"
                          size="xl"
                          :variant="!hasChangesPT ? 'ghost' : 'solid'"
                          :disabled="!hasChangesPT" 
                          @click="saveChangesPT"/>
                    </div>
                </div>
            </div>

            <div class="flex p-8">
                <div v-if="!isFetching" class="m-auto">
                    <h1 class="text-xl font-medium mb-4">Exam - {{ highestScoresStore.weightedScores[2] * 100 }}%</h1>
                    <ScoresField v-model="examHighestScore" label="" size="xl"/>
                    <UButton
                      class="mt-4"
                      label="Save"
                      size="xl"
                      :variant="examHighestScore === highestScoresStore.examHighestScore ? 'ghost' : 'solid'"
                      :disabled="examHighestScore === highestScoresStore.examHighestScore" 
                      @click="saveChangesExam"/>
                </div>
                <div v-else class="w-[90px] h-[168px] m-auto">
                    <USkeleton class="w-full h-[24px] mb-2" />
                    <USkeleton class="w-[70%] h-[24px] mb-6" />
                    <USkeleton class="w-full h-[30px] mb-4" />
                    <USkeleton class="w-[70%] h-[24px]" />
                </div>
            </div>
        </div>
    </div>
</template>