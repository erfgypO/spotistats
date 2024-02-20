// Utilities
import { defineStore } from 'pinia'
import httpClient from "@/store/httpClient";
import {LoginResponse} from "@/types/response";
import {useUserStore} from "@/store/user";

export const useAppStore = defineStore('app', {
  state: () => ({
    accessToken: "",
    expiresAt: 0,
  }),
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
        console.error(e);
      }
    }
  },
  persist: true,
})
