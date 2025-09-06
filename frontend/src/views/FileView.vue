<script setup lang="ts">
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { onMounted } from 'vue';
import { fileViewStore } from '../store/fileViewStore';
import { ref } from 'vue';
import { studentsStore } from '../store/studentsStore';

const props = defineProps<{
  paddingTop: string
}>()

const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

const toast = useToast();

const fileStore = fileViewStore();

const studentStore = studentsStore();

const selectedTrack = ref<string>(fileStore.track);

const chooseFile = async () => {
  emit('loading', true);
  fileStore.setDoneReading(false);
  const { OpenFileDialog } = await import('../../wailsjs/go/main/App');
  await OpenFileDialog();
};

async function handleSelect(newTrack: string) {
  const { ChangeTrack } = await import('../../wailsjs/go/main/App');
  await ChangeTrack(newTrack, fileStore.tracks.indexOf(newTrack));
  fileStore.setTrack(newTrack);
}

onMounted(async () => {

  const { Tracks } = await import('../../wailsjs/go/main/App');
  fileStore.setTracks(await Tracks());
  if (fileStore.tracks.length > 0) {
    fileStore.setTrack(fileStore.tracks[0]);
  }

  EventsOn('excel:is_ecr', (isECR: boolean) => {
    if (!isECR) {
      toast.add({ title: 'File Might not be ECR', description: "Missing Sheets", color: 'error' });
    }
  });


  EventsOn('excel:students_male', (male: string[]) => {
    studentStore.setMales(male);
    console.log(studentStore.males);
  });
  console.log(studentStore.males);

  EventsOn('excel:students_female', (female: string[]) => {
    studentStore.setFemales(female);
    console.log(studentStore.females);
  });
  console.log(studentStore.females);

  EventsOn('excel:done_reading', async () => {
    const { FileName } = await import('../../wailsjs/go/main/App');
    fileStore.setFileName(await FileName())
    fileStore.setDoneReading(true);
    emit('loading', false);
  })


});
</script>

<template>
  <div class="h-full w-full relative flex items-center justify-center gap-2">

    <header v-if="fileStore.doneReading" class="absolute top-0 right-0 left-0 flex flex-col" :style="{ paddingTop: props.paddingTop }">
      <h1 class="mt-14 mx-auto font-bold text-6xl text-(--ui-primary)">
        {{ fileStore.fileName }}
      </h1>
      <div class="mx-auto mt-14">
        <USelect
          v-model="selectedTrack"
          :items="fileStore.tracks"
          class="w-auto min-w-70"
          @update:modelValue="handleSelect"
        />
      </div>
    </header>
    
    <UButton
      label="Choose File"
      icon="i-lucide-file"
      target="_blank"
      size="xl"
      @click="chooseFile" />

  </div>
</template>