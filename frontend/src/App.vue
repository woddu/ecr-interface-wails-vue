<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui';
import { nextTick, onMounted, ref } from 'vue';
import { EventsOn } from '../wailsjs/runtime/runtime';

const isLoading = ref<boolean>(false);

const navItems = ref<NavigationMenuItem[]>([
  { label: 'File', to: '/' , disabled: false},
  { label: 'Scores', to: '/scores', disabled: true},
  { label: 'Students', to: '/students', disabled: true},
]);

const headerRef = ref<HTMLElement | null>(null)
const mainPaddingTop = ref('0px')

onMounted( async () => {
  await nextTick() // wait until DOM is rendered
  if (headerRef.value) {
    const height = headerRef.value.offsetHeight
    mainPaddingTop.value = `${height}px`
  }
  console.log('App mounted');
  EventsOn('excel:progress', (message: string) => {
    console.log('Progress:', message);
  })
  EventsOn('excel:done', () => {
    isLoading.value = false;
  })
  EventsOn('excel:done_reading', () => {
    isLoading.value = false;
    navItems.value[1].disabled = false;
    navItems.value[2].disabled = false;
  })
})

</script>

<template>
  <UApp class="h-screen w-screen">
    <UProgress class="fixed top-0 left-0 right-0 z-50" v-if="isLoading" />
    <div ref="headerRef" class="fixed top-0 left-0 right-0 z-40">
      <nav class="w-screen flex items-center justify-between p-2 bg-(--ui-bg) border-b border-(--ui-border)">

        <UNavigationMenu
          class="w-full justify-center"
          orientation="horizontal"
          :items="navItems"
          :ui="{
            item: 'mx-2',   // horizontal margin between items
            link: 'px-4 py-2 text-md' // padding inside each link
          }" />
      </nav>
    </div>

    <router-view
      class="flex-1"
      :style="{ paddingTop: mainPaddingTop }" v-slot="{ Component }">
      <component
        :is="Component"
        :padding-top="mainPaddingTop"
        class="flex-1"
        @loading="(value: boolean) => isLoading = value"
      />
    </router-view>
  </UApp>
</template>
