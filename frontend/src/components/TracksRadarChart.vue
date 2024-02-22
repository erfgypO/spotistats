<script setup lang="ts">
import {
  Chart as ChartJS,
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend
} from 'chart.js';

import { Radar } from 'vue-chartjs'
import {useStatsStore} from "@/store/stats";
import { DefaultRadarChartOptions } from "@/types/ChartData";

ChartJS.register(
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend,
)


const statsStore = useStatsStore()
</script>

<template>
<Radar :options="DefaultRadarChartOptions" :data="statsStore.tracksChartData" />
</template>

<style scoped>

</style>
