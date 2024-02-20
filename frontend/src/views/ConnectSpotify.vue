<script setup lang="ts">
import {computed, onMounted} from "vue";
import {useUserStore} from "@/store/user";

const userStore = useUserStore();
onMounted(() => {
  userStore.getSpotifyAuthUrl();
});

const loading = computed(() => userStore.connectSpotifyUrl === '');
</script>

<template>
  <v-card max-width="750px" class="mx-auto">
    <v-card-title>Connect Spotify</v-card-title>
    <v-card-text>
      <v-row>
        <v-col cols="12">
          <v-btn block color="primary" :href="userStore.connectSpotifyUrl" :loading="loading">Connect to Spotify</v-btn>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<style scoped>

</style>
