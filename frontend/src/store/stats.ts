import {defineStore} from "pinia";
import {Stat} from "@/types/Stat";
import httpClient from "@/store/httpClient";
import {colors, rgbToRgba} from "@/utils/colors";
import {useMessageStore} from "@/store/message";

export const useStatsStore = defineStore('stats', {
  state: () => ({
    artists: [] as Stat[],
    tracks: [] as Stat[],
    albums: [] as Stat[],
    timeRange: ['1h', '24h', '7d', '30d', '90d', '365d', 'all'],
    selectedTimeRange: 1,
  }),
  getters: {

    artistChartData(): any {
      return {
        labels: this.artists.map((stat: Stat) => stat.name),
        datasets: [
          {
            label: 'Top Artists',
            data: this.artists.map((stat: Stat) => stat.datapointCount * 10 / 60),
            backgroundColor: rgbToRgba(colors.primary, 0.2),
            borderColor: colors.primary,
            borderWidth: 1
          }
        ]
      }
    },
    tracksChartData(): any {
      return {
        labels: this.tracks.map((stat: Stat) => stat.name.replace('(', '-(').split('-').map((s: string) => s.trim())),
        datasets: [
          {
            label: 'Top Tracks',
            data: this.tracks.map((stat: Stat) => stat.datapointCount * 10 / 60),
            backgroundColor: rgbToRgba(colors.secondary, 0.2),
            borderColor: colors.secondary,
            borderWidth: 1
          }
        ]
      }
    },
    albumsChartData(): any {
      return {
        labels: this.albums.map((stat: Stat) => stat.name.replace('(', '-(').split('-').map((s: string) => s.trim())),
        datasets: [
          {
            label: 'Top Albums',
            data: this.albums.map((stat: Stat) => stat.datapointCount * 10 / 60),
            backgroundColor: rgbToRgba(colors.tertiary, 0.2),
            borderColor: colors.tertiary,
            borderWidth: 1
          }
        ]
      }
    }
  },
  actions: {
    async fetchStats(after: number) {
      const messageStore = useMessageStore();
      try {
        interface StatsResponse {
          artists: Stat[];
          tracks: Stat[];
          albums: Stat[];
        }
        const response = await httpClient.get<StatsResponse>(`/auth/stats?after=${after}`);

        if(response.status === 200) {
          this.artists = response.data.artists.slice(0,10); //.filter(stat => stat.datapointCount >= 6);
          this.tracks = response.data.tracks.slice(0,10); //.filter(stat => stat.datapointCount >= 6);
          this.albums = response.data.albums.slice(0, 10);
        }
      } catch (e) {
        messageStore.showMessage("Failed to fetch stats", "error")
      }
    },
  },
})
