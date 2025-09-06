import { defineStore } from "pinia";
import { ref } from "vue";

export const useHighestScoresStore = defineStore('highestScores', () => {
    const wwHighestScores = ref<number[]>([]);
    const ptHighestScores = ref<number[]>([]);
    const examHighestScore = ref<number>(0);
    const weightedScores = ref<number[]>([]);
    
    function setWwHighestScores(scores: number[]) {
        wwHighestScores.value = scores;
    }
    function setPtHighestScores(scores: number[]) {
        ptHighestScores.value = scores;
    }
    function setExamHighestScore(score: number) {
        examHighestScore.value = score;
    }
    function setWeightedScores(scores: number[]) {
        weightedScores.value = scores;
    }

    return {
        wwHighestScores,
        ptHighestScores,
        examHighestScore,
        weightedScores,
        setWwHighestScores,
        setPtHighestScores,
        setExamHighestScore,
        setWeightedScores
    };
});