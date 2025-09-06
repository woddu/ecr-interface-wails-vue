<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { studentsStore } from '../../store/studentsStore';
import { EventsOn } from '../../../wailsjs/runtime/runtime';

const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

const studentStore = studentsStore();
 
const searchTerm = ref<string>('')

const filteredMales = computed(() => {
  const term = searchTerm.value.toUpperCase()
  return studentStore.males.filter(student =>
    student.toUpperCase().includes(term)
  )
})

const filteredFemales = computed(() => {
  const term = searchTerm.value.toUpperCase()
  return studentStore.females.filter(student =>
    student.toUpperCase().includes(term)
  )
})

async function studentClicked(index: number) {
  const { GetStudent } = await import('../../../wailsjs/go/main/App');
  await GetStudent(index)
}

onMounted(async () => {
    EventsOn('excel:student_scores', (studentData: any) => {
        console.log(studentData);
    })
})
</script>

<template>
    <div class="h-full w-full flex flex-col gap-4 p-4">
        <UInput v-model="searchTerm" class="my-4 max-w-2xs" icon="i-lucide-search" size="xl" variant="outline" placeholder="Search..." />
        <div class="h-4/10 w-full flex flex-col gap-4">
            <h2 class="text-2xl font-bold">Male Students</h2>
            <ul class=" overflow-y-scroll">
                <li v-for="(student, index) in filteredMales" :key="index"
                    class="mt-2 text-xl cursor-pointer hover:text-(--ui-primary) hover:font-bold"
                    @click="studentClicked(index)">
                    {{ student }}
                </li>
            </ul>
        </div>

        <div class="h-4/10 w-full flex flex-col gap-4">
            <h2 class="text-2xl font-bold">Female Students</h2>
            <ul class=" overflow-y-scroll">
                <li v-for="(student, index) in filteredFemales" :key="index"
                    class="mt-2 text-xl cursor-pointer hover:text-(--ui-primary) hover:font-bold"
                    @click="studentClicked(index)">
                    {{ student }}
                </li>
            </ul>
        </div>
    </div>
</template>