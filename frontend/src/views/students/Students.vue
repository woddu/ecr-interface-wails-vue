<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useStudentsStore } from '../../store/studentsStore';
import { EventsOn } from '../../../wailsjs/runtime/runtime';
import { useStudentStore } from '../../store/studentStore';
import { useRouter } from 'vue-router';

const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

const router = useRouter()

const studentsStore = useStudentsStore();

const studentStore = useStudentStore()

const searchTerm = ref<string>('')

const newStudent = ref<string>('')

const genders = ref(['Male', 'Female'])
const gender = ref('Male')

const filteredMales = computed(() => {
  const term = searchTerm.value.toUpperCase()
  return studentsStore.males.filter(student =>
    student.toUpperCase().includes(term)
  )
})

const filteredFemales = computed(() => {
  const term = searchTerm.value.toUpperCase()
  return studentsStore.females.filter(student =>
    student.toUpperCase().includes(term)
  )
})

async function studentClicked(index: number) {
  emit('loading', true)
  const { GetStudent } = await import('../../../wailsjs/go/main/App');
  await GetStudent(index)
}

async function addStudent() {
  if (newStudent.value.trim() === '') {
    return
  }
  emit('loading', true)
  const { AddStudent } = await import('../../../wailsjs/go/main/App');
  await AddStudent(newStudent.value.trim(), gender.value === 'Male')
  newStudent.value = ''
}

interface StudentData {
  row: number,
  name: string,
  writtenWorks: number[],
  performance: number[],
  exam: number,
}

onMounted(async () => {
  EventsOn('excel:done_getting_student', (studentData: StudentData) => {
    studentStore.setRow(studentData.row),
      studentStore.setName(studentData.name),
      studentStore.setWrittenWorks(studentData.writtenWorks),
      studentStore.setPerformanceTasks(studentData.performance),
      studentStore.setExam(studentData.exam)
    emit('loading', false)
    router.push({ name: 'StudentScores' })
  })
})
</script>

<template>
  <div class="h-full w-full flex flex-row">
    <div class="h-full w-1/2 flex flex-col gap-4 p-10">
      <UInput v-model="searchTerm" class="my-4 max-w-2xs" icon="i-lucide-search" size="xl" variant="outline"
        placeholder="Search..." />
      <div class="h-4/10 w-full flex flex-col gap-4">
        <h2 class="text-2xl font-bold">Male Students</h2>
        <ul class=" overflow-y-scroll">
          <li v-for="(student, index) in filteredMales" :key="index"
            class="mt-2 text-xl cursor-pointer hover:text-(--ui-primary) hover:font-bold" @click="studentClicked(index)">
            {{ student }}
          </li>
        </ul>
      </div>
      <div class="h-4/10 w-full flex flex-col gap-4">
        <h2 class="text-2xl font-bold">Female Students</h2>
        <ul class=" overflow-y-scroll">
          <li v-for="(student, index) in filteredFemales" :key="index"
            class="mt-2 text-xl cursor-pointer hover:text-(--ui-primary) hover:font-bold" @click="studentClicked(index)">
            {{ student }}
          </li>
        </ul>
      </div>
    </div>
    <div class="h-full w-1/2 flex flex-col gap-4 items-center justify-center p-10">
        <h2 class="text-3xl font-bold">Add Student</h2>
        <UFormField label="Name" size="xl">
          <UInput v-model="newStudent" class="max-w-md" variant="outline" placeholder="Add Student" />
        </UFormField>
        <UFormField label="Gender" size="xl">
          <USelect v-model="gender" :items="genders" />
        </UFormField>
        <UButton class="mt-4" color="primary" size="lg" @click="addStudent" :disabled="newStudent.trim() === ''">
          Add Student
        </UButton>
    </div>
  </div>
</template>