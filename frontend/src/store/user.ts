import {defineStore} from "pinia";
import httpClient from "@/store/httpClient";
import {UserResponse} from "@/types/response";

export const useUserStore = defineStore('user', {
  state: () => ({
    username: "",
    displayName: "",
    connectedToSpotify: false,
    datapointCount: 0,
    connectSpotifyUrl: ""
  }),
  actions: {
    async fetchUser() {
      try {
        const response = await httpClient.get<UserResponse>('/user/me');
        if(response.status === 200) {
          this.username = response.data.username;
          this.displayName = response.data.displayName;
          this.connectedToSpotify = response.data.connectedToSpotify;
          this.datapointCount = response.data.datapointCount;
        }
      } catch (e) {
        console.error(e);
      }
    },
    async getSpotifyAuthUrl() {
      try {
        interface ConnectSpotifyResponse {
          url: string;
        }
        const response = await httpClient.get<ConnectSpotifyResponse>('/auth-url');
        if(response.status === 200) {
          this.connectSpotifyUrl = response.data.url;
        }
      } catch (e) {
        console.error(e);
      }
    }
  },
  persist: true,
});
