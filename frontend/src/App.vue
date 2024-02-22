<template>
  <v-app theme="dark">
    <v-app-bar color="primary">
      <v-app-bar-title>Soundstats</v-app-bar-title>
    </v-app-bar>
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import {onMounted} from "vue";
import {useUserStore} from "@/store/user";
import {useAppStore} from "@/store/app";

const userStore = useUserStore();
const appStore = useAppStore();

onMounted(async () => {
  if (appStore.expiresAt * 1000 > Date.now()) {
    await userStore.fetchUser();
  }
});
</script>

<style>

</style>
