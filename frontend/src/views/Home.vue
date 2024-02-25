<template>
  <v-container fluid>
    <v-row>
      <page-title-col title="Dashboard" />
      <v-col cols="12" md="6" class="d-flex align-center">
        <div class="d-flex flex-row" :class="$vuetify.display.mdAndUp ? 'ml-auto' : 'mx-auto'">
        <v-btn-toggle divided variant="outlined" v-model="btnGroupModel" color="primary"
                      @update:modelValue="onTimeRangeChange">
          <v-btn v-for="range in statsStore.timeRange" :key="range" :text="range" />
        </v-btn-toggle>
        </div>
      </v-col>
      <v-col cols="12">
        <v-sheet class="pa-3" rounded>
          <span class="text-h4">Hourly Stats</span>
          <div class="hourly-chart-container">
        <hourly-stats-chart />
          </div>
        </v-sheet>
      </v-col>
      <v-col cols="12" md="6" xl="4">
        <v-sheet class="pa-3" rounded>
          <span class="text-h4">Top Artists</span>
          <div class="chart-container">
            <artists-radar-chart />
          </div>
          <stats-table :stats="statsStore.artists" name="Artist" />
        </v-sheet>
      </v-col>
      <v-col cols="12" md="6" xl="4">
        <v-sheet class="pa-3" rounded>
          <span class="text-h4">Top Tracks</span>
          <div class="chart-container">
            <tracks-radar-chart />
          </div>
          <stats-table :stats="statsStore.tracks" name="Track" />
        </v-sheet>
      </v-col>
      <v-col cols="12" md="6" offset-md="3" xl="4" offset-xl="0">
        <v-sheet class="pa-3" rounded>
          <span class="text-h4">Top Albums</span>
          <div class="chart-container">
            <albums-radar-chart />
          </div>
          <stats-table :stats="statsStore.albums" name="Album" />
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import {onMounted} from "vue";
import {useStatsStore} from "@/store/stats";
import ArtistsRadarChart from "@/components/ArtistsRadarChart.vue";
import TracksRadarChart from "@/components/TracksRadarChart.vue";
import PageTitleCol from "@/components/PageTitleCol.vue";
import AlbumsRadarChart from "@/components/AlbumsRadarChart.vue";
import StatsTable from "@/components/StatsTable.vue";
import HourlyStatsChart from "@/components/HourlyStatsChart.vue";

const btnGroupModel = defineModel('btnGroupModel', { type: Number });
const statsStore = useStatsStore();

function createAfterDate(value: string) {
  const date = new Date();

  switch (value) {
    case '1h':
      date.setHours(date.getHours() - 1);
      break;
    case '24h':
      date.setDate(date.getDate() - 1);
      break;
    case '7d':
      date.setDate(date.getDate() - 7);
      break;
    case '30d':
      date.setDate(date.getDate() - 30);
      break;
    case '90d':
      date.setDate(date.getDate() - 90);
      break;
    case '365d':
      date.setDate(date.getDate() - 365);
      break;
    case 'all':
      date.setTime(0);
      break;
  }

  return Math.trunc(date.getTime() / 1000);
}

async function onTimeRangeChange() {
  statsStore.selectedTimeRange = btnGroupModel.value!;
  await statsStore.fetchStats(createAfterDate(statsStore.timeRange[statsStore.selectedTimeRange]));
}

onMounted(() => {
  btnGroupModel.value = statsStore.selectedTimeRange;
  onTimeRangeChange();
  statsStore.fetchHourlyStats();
});
</script>

<style scoped>
.hourly-chart-container {
  height: 300px;
  width: 100%;
}

.chart-container {
  height: 400px;
  width: 100%;
}
</style>
