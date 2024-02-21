import {defineStore} from "pinia";
import {Stat} from "@/types/Stat";
import httpClient from "@/store/httpClient";
import {ChartData} from "@/types/ChartData";
import colors from "@/utils/colors";

export const useStatsStore = defineStore('stats', {
  state: () => ({
    artists: [] as Stat[],
    tracks: [] as Stat[],
  }),
  getters: {
    artistChartData(): ChartData {
      return {
        series: [
          {
            name: 'Top Artists',
            data: this.artists.map((stat: Stat) => stat.percentage * 100),
          }
        ],
        chartOptions: {
          labels: this.artists.map((stat: Stat) => stat.name),
          theme: {
            mode: 'dark'
          },
          stroke: {
            show: true,
            width: 2,
            colors: [colors.primary],
            dashArray: 0
          },
          fill: {
            opacity: 0.5,
            colors: [colors.primary]
          }
        }
      }
    },
    tracksChartData(): ChartData {
      return {
        series: [
          {
            name: 'Top Tracks',
            data: this.tracks.map((stat: Stat) => stat.percentage * 100),
          }
        ],
        chartOptions: {
          labels: this.tracks.map((stat: Stat) => stat.name),
          theme: {
            mode: 'dark'
          },
          stroke: {
            show: true,
            width: 2,
            colors: [colors.secondary],
            dashArray: 0
          },
          fill: {
            opacity: 0.5,
            colors: [colors.secondary]
          }
        }
      }
    }
  },
  actions: {
    async fetchStats(after: number) {
      try {
        interface StatsResponse {
          artists: Stat[];
          tracks: Stat[];
        }
        const response = await httpClient.get<StatsResponse>(`/auth/stats?after=${after}`);

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
