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
      <v-col cols="12" md="6" xl="4">
        <v-sheet class="pa-3" rounded>
          <span class="text-h4">Top Artists</span>
          <div class="chart-container">
            <artists-radar-chart />
          </div>
          <v-table>
            <thead>
              <tr>
                <th></th>
                <th>Artist</th>
                <th>Time</th>
                <th class="text-center">
                  <v-icon icon="mdi-spotify" color="green" />
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(artist, index) in statsStore.artists" :key="artist.name">
                <td>{{index + 1}}</td>
                <td>{{artist.name}}</td>
                <td>{{ secondsToString(artist.datapointCount * 10)}}</td>
                <td class="text-center"><a :href="artist.spotifyUrl" target="_blank">
                  open
                </a></td>
              </tr>
            </tbody>
          </v-table>
        </v-sheet>
      </v-col>
      <v-col cols="12" md="6" xl="4">
        <v-sheet class="pa-3" rounded>
          <span class="text-h4">Top Tracks</span>
          <div class="chart-container">
            <tracks-radar-chart />
          </div>
          <v-table>
            <thead>
              <tr>
                <th></th>
                <th>Track</th>
                <th>Time</th>
                <th class="text-center">
                  <v-icon color="green" icon="mdi-spotify" />
                </th>
              </tr>
            </thead>
            <tbody>
            <tr v-for="(track, index) in statsStore.tracks" :key="track.name">
              <td>{{index + 1}}</td>
              <td>{{track.name}}</td>
              <td>{{ secondsToString(track.datapointCount * 10)}}</td>
              <td class="text-center"><a :href="track.spotifyUrl" target="_blank">open</a></td>
            </tr>
            </tbody>
          </v-table>
        </v-sheet>
      </v-col>
      <v-col cols="12" md="6" offset-md="3" xl="4" offset-xl="0">
        <v-sheet class="pa-3" rounded>
          <span class="text-h4">Top Albums</span>
          <div class="chart-container">
            <albums-radar-chart />
          </div>
          <v-table>
            <thead>
            <tr>
              <th></th>
              <th>Album</th>
              <th>Time</th>
              <th class="text-center">
                <v-icon icon="mdi-spotify" color="green" />
              </th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="(album, index) in statsStore.albums" :key="album.name">
              <td>{{index + 1}}</td>
              <td>{{album.name}}</td>
              <td>{{ secondsToString(album.datapointCount * 10)}}</td>
              <td class="text-center"><a :href="album.spotifyUrl" target="_blank">
                open
              </a></td>
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
import ArtistsRadarChart from "@/components/ArtistsRadarChart.vue";
import TracksRadarChart from "@/components/TracksRadarChart.vue";
import PageTitleCol from "@/components/PageTitleCol.vue";
import AlbumsRadarChart from "@/components/AlbumsRadarChart.vue";

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
});
</script>

<style scoped>
.chart-container {
  height: 400px;
  width: 100%;
}
</style>
