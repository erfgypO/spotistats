// Utilities
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    accessToken: "",
    expiresAt: 0,
  }),
})
