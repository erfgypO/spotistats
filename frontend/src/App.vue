<template>
  <v-app theme="dark">
    <v-app-bar color="primary">
      <v-app-bar-title>
        <router-link to="/" style="text-decoration: none; color: inherit">
          Soundstats
        </router-link>
        </v-app-bar-title>
      <v-spacer />
      <v-btn icon="mdi-account" v-if="appStore.isAuthenticated" to="/account" />
    </v-app-bar>
    <v-main>
      <router-view />
      <v-snackbar v-model="messageStore.messageBarVisible"
                  :color="messageBarColor"
                  timeout="2500">
        {{ messageStore.messageBarText }}
      </v-snackbar>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import {onMounted, watch} from "vue";
import {useUserStore} from "@/store/user";
import {useAppStore} from "@/store/app";
import {useMessageStore} from "@/store/message";
import {storeToRefs} from "pinia";

const userStore = useUserStore();
const appStore = useAppStore();
const messageStore = useMessageStore();
onMounted(async () => {
  if (appStore.expiresAt * 1000 > Date.now()) {
    await userStore.fetchUser();
  }
});

const { messageBarColor, messageBarText } = storeToRefs(messageStore);

watch(messageBarText, () => {
  messageStore.messageBarVisible = true;
});
</script>

<style>

</style>
