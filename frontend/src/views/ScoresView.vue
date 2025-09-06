<script setup lang="ts">
import ScoresField from '../components/ScoresField.vue'
import { computed, onMounted, ref } from 'vue'

// Turn each into a ref
const wwHighestScores = ref<number[]>([])
const originalWwHighestScores = ref<number[]>([])
const ptHighestScores = ref<number[]>([])
const originalPtHighestScores = ref<number[]>([])
const examHighestScore = ref<number>(0)
const originalExamHighestScore = ref<number>(0)
const weightedScores = ref<number[]>([])

const isFetching = ref<boolean>(false)

const hasChangesWW = computed(() => {
  return wwHighestScores.value.some(
    (score, i) => score !== originalWwHighestScores.value[i]
  )
})

const hasChangesPT = computed(() => {
    return ptHighestScores.value.some(
      (score, i) => score !== originalPtHighestScores.value[i]
    )
})

onMounted(async () => {
    isFetching.value = true
    const { Scores } = await import('../../wailsjs/go/main/App');
    await Scores().then(res => {
        console.log('Scores:', res);
        wwHighestScores.value = res.wwHighestScores;
        originalWwHighestScores.value = [...res.wwHighestScores];
        ptHighestScores.value = res.ptHighestScores;
        originalPtHighestScores.value = [...res.ptHighestScores];
        examHighestScore.value = res.examHighestScore;
        originalExamHighestScore.value = res.examHighestScore;
        weightedScores.value = res.weightedScores;
    }).finally(() => isFetching.value = false);
})

</script>

<template>
    <div class="h-full w-full flex flex-col items-center justify-center p-10">
        <div class="grid grid-cols-6 my-auto">
            

            <div class="flex flex-col gap-8 col-span-5">
                <div class="grid grid-cols-8">
                    <div v-if="!isFetching" class="col-span-7">
                        <h1 class="text-xl font-medium">Written Works - {{ weightedScores[0] * 100 }}%</h1>
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
                          :disabled="!hasChangesWW" />
                    </div>
                </div>
                <div class="grid grid-cols-8">
                    <div v-if="!isFetching" class="col-span-7">
                        <h1 class="text-xl font-medium">Performance Tasks - {{ weightedScores[1] * 100 }}%</h1>
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
                          :disabled="!hasChangesPT" />
                    </div>
                </div>
            </div>

            <div class="flex p-8">
                <div v-if="!isFetching" class="m-auto">
                    <h1 class="text-xl font-medium mb-4">Exam - {{ weightedScores[2] * 100 }}%</h1>
                    <ScoresField v-model="examHighestScore" label="" size="xl"/>
                    <UButton
                      class="mt-4"
                      label="Save"
                      size="xl"
                      :variant="examHighestScore === originalExamHighestScore ? 'ghost' : 'solid'"
                      :disabled="examHighestScore === originalExamHighestScore" />
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