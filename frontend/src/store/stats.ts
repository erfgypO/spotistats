import {defineStore} from "pinia";
import {Stat} from "@/types/Stat";
import httpClient from "@/store/httpClient";

export const useStatsStore = defineStore('stats', {
  state: () => ({
    artists: [] as Stat[],
    tracks: [] as Stat[],
  }),
  actions: {
    async fetchStats(after: number) {
      try {
        interface StatsResponse {
          artists: Stat[];
          tracks: Stat[];
        }
        const response = await httpClient.get<StatsResponse>(`/stats?after=${after}`);

        if(response.status === 200) {
          this.artists = response.data.artists.slice(0,10); //.filter(stat => stat.datapointCount >= 6);
          this.tracks = response.data.tracks.slice(0,10); //.filter(stat => stat.datapointCount >= 6);
        }
      } catch (e) {
        console.error(e);
      }
    },
  },
})
