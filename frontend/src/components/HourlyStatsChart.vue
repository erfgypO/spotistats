<script setup lang="ts">
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale);

import {useStatsStore} from "@/store/stats";
import {computed} from "vue";

const statsStore = useStatsStore()


const config = computed(() => {
  return {
    type: 'bar',
    data: statsStore.hourlyChartData,
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        x: {
          stacked: true,
        },
        y: {
          stacked: true
        }
      }
    }
  }
});
</script>

<template>
<Bar :data="config.data" :options="config.options" />
</template>

<style scoped>

</style>
