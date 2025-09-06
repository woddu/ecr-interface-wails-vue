<script setup lang="ts">
import { ref, computed } from 'vue'
import { studentsStore } from '../store/studentsStore';

const studentStore = studentsStore();

</script>

<template>
    <div class="flex flex-col gap-4 m-4">
        <h2 class="text-xl font-bold">Male Students</h2>
        <ul>
            <li v-for="student in studentStore.males" :key="student">{{ student }}</li>
        </ul>

        <h2 class="text-xl font-bold">Female Students</h2>
        <ul>
            <li v-for="student in studentStore.females" :key="student">{{ student }}</li>
        </ul>
    </div>
</template>