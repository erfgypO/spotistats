import {defineStore} from "pinia";
import {HourlyStat, Stat} from "@/types/Stat";
import httpClient from "@/store/httpClient";
import {colors, rgbToRgba} from "@/utils/colors";
import {useMessageStore} from "@/store/message";

export const useStatsStore = defineStore('stats', {
  state: () => ({
    artists: [] as Stat[],
    tracks: [] as Stat[],
    albums: [] as Stat[],
    hourlyStats: [] as HourlyStat[],
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
    },
    hourlyChartData(): any {
      const labels = [];
      for(let i = 0; i < 24; i++) {
        labels.push(i);
      }

      const tracks =[...new Set(this.hourlyStats.map((stat: HourlyStat) => stat.songName))];
      const datasets = [];

      for(const track of tracks) {
        const data = []//this.hourlyStats.filter((stat: HourlyStat) => stat.songName === track).map((stat: HourlyStat) => stat.seconds);

        for(let i = 0; i < 24; i++) {
          const stat = this.hourlyStats.find((stat: HourlyStat) => stat.songName === track && stat.hour === i);
          data.push(stat ? stat.seconds / 60 : 0);
        }

        datasets.push({
          label: track,
          data,
          backgroundColor: rgbToRgba(colors.quaternary, 0.1),
          borderColor: colors.quaternary,
          borderWidth: 2
        });
      }

      return {
        labels,
        datasets
      }
    },
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
          this.artists = response.data.artists.slice(0,10);

          if (this.artists.length < 10) {
            for(let i = this.artists.length; i < 10; i++) {
              this.artists.push({name: "", percentage: 0, datapointCount: 0, spotifyUrl: ""});
            }
          }
          this.tracks = response.data.tracks.slice(0,10);
          if (this.tracks.length < 10) {
            for(let i = this.tracks.length; i < 10; i++) {
              this.tracks.push({name: "", percentage: 0, datapointCount: 0, spotifyUrl: ""});
            }
          }
          this.albums = response.data.albums.slice(0, 10);
          if (this.albums.length < 10) {
            for(let i = this.albums.length; i < 10; i++) {
              this.albums.push({name: "", percentage: 0, datapointCount: 0, spotifyUrl: ""});
            }
          }
        }
      } catch (e) {
        messageStore.showMessage("Failed to fetch stats", "error")
      }
    },
    async fetchHourlyStats() {
      const messageStore = useMessageStore();
      try {

        const date = new Date();
        date.setHours(0, 0, 0, 0);

        const after = Math.floor(date.getTime() / 1000);
        const offset = date.getTimezoneOffset() / 60 * -1;
        const response = await httpClient.get<HourlyStat[]>(`/auth/stats/hourly?after=${after}&tzOffset=${offset}`);

        if(response.status === 200) {
          this.hourlyStats = response.data.map((stat: HourlyStat) => {
            if(stat.hour === 24) {
              stat.hour = 0;
            }

            return stat;
          }).toSorted((a: HourlyStat, b: HourlyStat) => a.hour - b.hour);
        }
      } catch (e) {
        messageStore.showMessage("Failed to fetch hourly stats", "error")
      }
    }
  },
})
