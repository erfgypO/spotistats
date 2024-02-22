// Utilities
import { defineStore } from 'pinia'
import httpClient from "@/store/httpClient";
import {ErrorResponse, LoginResponse} from "@/types/response";
import {useUserStore} from "@/store/user";
import {AxiosError} from "axios";
import {useMessageStore} from "@/store/message";

export const useAppStore = defineStore('app', {
  state: () => ({
    accessToken: "",
    expiresAt: 0,
  }),
  getters: {
    isAuthenticated(): boolean {
      return this.accessToken.length > 0;
    }
  },
  actions: {
    async login(username: string, password: string) {
      try {
        const response = await httpClient.post<LoginResponse>('/sign-in', {
          username,
          password
        });

        if(response.status === 200) {
          this.accessToken = response.data.accessToken;
          this.expiresAt = response.data.expiresAt;
          localStorage.setItem('token', this.accessToken);

          const userStore = useUserStore();
          await userStore.fetchUser();
        }
      } catch (e) {
        const response = (e as AxiosError)?.response?.data as ErrorResponse;

        if(response) {
          useMessageStore().showMessage(response.error, "error");
        }
      }
    },
    async signUp(username: string, password: string) {
      try {
        const response = await httpClient.post<LoginResponse>('/sign-up', {
          username,
          password
        });

        if(response.status === 200) {
          this.accessToken = response.data.accessToken;
          this.expiresAt = response.data.expiresAt;
          localStorage.setItem('token', this.accessToken);

          const userStore = useUserStore();
          await userStore.fetchUser();
        }
      } catch (e) {
        const response = (e as AxiosError)?.response?.data as ErrorResponse;

        if(response) {
          useMessageStore().showMessage(response.error, "error");
        }
      }
    }
  },
  persist: true,
})
