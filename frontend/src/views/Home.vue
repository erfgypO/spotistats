<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" sm="8" lg="10">
        <span class="text-h3">Stats</span>
      </v-col>
      <v-col cols="12" sm="4" lg="2">
        <v-select variant="outlined"  v-model="selectedTimeRange" :items="timeRange" label="Time range" item-title="text" @update:modelValue="onTimeRangeChange" />
      </v-col>
      <v-col cols="12" md="6">
        <v-sheet class="pa-3" rounded>
          <span class="text-h4">Top Artists</span>
          <v-table>
            <thead>
              <tr>
                <th></th>
                <th>Artist</th>
                <th>Time</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(artist, index) in statsStore.artists" :key="artist.name">
                <td>{{index + 1}}</td>
                <td>{{artist.name}}</td>
                <td>{{ secondsToString(artist.datapointCount * 10)}}</td>
              </tr>
            </tbody>
          </v-table>
        </v-sheet>
      </v-col>
      <v-col cols="12" md="6">
        <v-sheet class="pa-3" rounded>
          <span class="text-h4">Top Tracks</span>
          <v-table>
            <thead>
              <tr>
                <th></th>
                <th>Track</th>
                <th>Time</th>
              </tr>
            </thead>
            <tbody>
            <tr v-for="(track, index) in statsStore.tracks" :key="track.name">
              <td>{{index + 1}}</td>
              <td>{{track.name}}</td>
              <td>{{ secondsToString(track.datapointCount * 10)}}</td>
            </tr>
            </tbody>
          </v-table>
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import {onMounted} from "vue";
import {useStatsStore} from "@/store/stats";
import { secondsToString } from "@/utils/secondsToString";

const timeRange = [
  {text: 'Last 1 hour', value: '1h'},
  {text: 'Last 24 hours', value: '24h'},
  {text: 'Last 7 days', value: '7d'},
  {text: 'Last 30 days', value: '30d'},
  {text: 'Last 90 days', value: '90d'},
  {text: 'Last 365 days', value: '365d'},
  {text: 'All time', value: 'all'}
];

const selectedTimeRange = defineModel('selectedTimeRange', { default: '1h'});
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
  await statsStore.fetchStats(createAfterDate(selectedTimeRange.value));
}

onMounted(() => {
  statsStore.fetchStats(createAfterDate(selectedTimeRange.value));
});
</script>
