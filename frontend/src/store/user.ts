import {defineStore} from "pinia";
import httpClient from "@/store/httpClient";
import {ErrorResponse, UserResponse} from "@/types/response";
import {AxiosError} from "axios";
import {useMessageStore} from "@/store/message";

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
        const response = await httpClient.get<UserResponse>('/auth/user/me');
        if(response.status === 200) {
          this.username = response.data.username;
          this.displayName = response.data.displayName;
          this.connectedToSpotify = response.data.connectedToSpotify;
          this.datapointCount = response.data.datapointCount;
        }

      } catch (e) {
        const response = (e as AxiosError)?.response?.data as ErrorResponse;

        if(response) {
          useMessageStore().showMessage(response.error, "error");
        }
      }
    },
    async updatePassword(password: string) {
      try {
        const response = await httpClient.put('/auth/user/update-password', {
          password
        });

        if(response.status === 200) {
          useMessageStore().showMessage("Password updated", "success");
        } else {
          useMessageStore().showMessage("Failed to update password", "error");
        }
      } catch (e) {
        const response = (e as AxiosError)?.response?.data as ErrorResponse;

        if(response) {
          useMessageStore().showMessage(response.error, "error");
        }
      }
    },
    async fetchSpotifyAuthUrl() {
      try {
        interface ConnectSpotifyResponse {
          url: string;
        }
        const response = await httpClient.get<ConnectSpotifyResponse>('/auth/url');
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
