<script setup lang="ts">
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { onMounted } from 'vue';
import { useFileViewStore } from '../store/fileViewStore';
import { ref } from 'vue';
import { useHighestScoresStore } from '../store/highestScoresStore';

const props = defineProps<{
  paddingTop: string
}>()

const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

const isLoading = ref<boolean>(false);

const toast = useToast();

const fileStore = useFileViewStore();

const highestScoresStore = useHighestScoresStore();

const selectedTrack = ref<string>(fileStore.track);

const isFirstSem = ref<boolean>(true);

const chooseFile = async () => {
  emit('loading', true);
  fileStore.setDoneReading(false);
  const { OpenFileDialog } = await import('../../wailsjs/go/main/App');
  await OpenFileDialog();
};

async function handleSelect(newTrack: string) {
  if (fileStore.track != newTrack && fileStore.doneReading){
    emit('loading', true);
    const { ChangeTrack } = await import('../../wailsjs/go/main/App');
    await ChangeTrack(newTrack, fileStore.tracks.indexOf(newTrack));
    fileStore.setTrack(newTrack);
  }
}

async function changeSem() {
  emit('loading', true);
  isLoading.value = true;
  const { ChangeSem } = await import('../../wailsjs/go/main/App');
  await ChangeSem(isFirstSem.value);
  
}

onMounted(async () => {

  if(!fileStore.doneReading) {
    const { Tracks } = await import('../../wailsjs/go/main/App');
    await Tracks().then((tracks: string[]) => {
      fileStore.setTracks(tracks);
      if (fileStore.tracks.length > 0) {
        fileStore.setTrack(fileStore.tracks[0]);
        selectedTrack.value = fileStore.tracks[0];
      }
    });
  }

  EventsOn('excel:is_ecr', (isECR: boolean) => {
    if (!isECR) {
      toast.add({ title: 'File Might not be ECR', description: "Missing Sheets", color: 'error' });
    }
  });


  EventsOn('excel:done_reading', async () => {
    const { FileName } = await import('../../wailsjs/go/main/App');
    fileStore.setFileName(await FileName())
    const { Scores } = await import('../../wailsjs/go/main/App');
    await Scores().then(res => {
        highestScoresStore.setWwHighestScores(res.wwHighestScores);
        highestScoresStore.setPtHighestScores(res.ptHighestScores);
        highestScoresStore.setExamHighestScore(res.examHighestScore);
        highestScoresStore.setWeightedScores(res.weightedScores);
    }).finally(() => {
      fileStore.setDoneReading(true);
      emit('loading', false);
      isLoading.value = false;
    });
  })

  EventsOn('excel:track_changed', (weightedScores: number[]) => {
    highestScoresStore.setWeightedScores(weightedScores);
    console.log(weightedScores);
    emit('loading', false);
  });

});
</script>

<template>
  <div class="h-full w-full relative flex items-center justify-center gap-2">

    <header v-if="fileStore.doneReading" class="absolute top-0 right-0 left-0 flex flex-col" :style="{ paddingTop: props.paddingTop }">
      <h1 class="mt-14 mx-auto font-bold text-6xl text-(--ui-primary)">
        {{ fileStore.fileName }}
      </h1>
      <USelect
        v-model="selectedTrack"
        :items="fileStore.tracks"
        class="w-auto min-w-70 mx-auto mt-12"
        @update:modelValue="handleSelect"
      />
      <UCheckbox
        v-model="isFirstSem"
        class="w-auto mx-auto mt-12"
        label="First Semester"
        @change="changeSem"
        :disabled="isLoading"
      />
    </header>
    
    <UButton
      label="Choose File"
      icon="i-lucide-file"
      target="_blank"
      size="xl"
      @click="chooseFile" 
      :disabled="isLoading && fileStore.doneReading"
    />

  </div>
</template>