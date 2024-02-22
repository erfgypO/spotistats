<script setup lang="ts">
import {useRouter} from "vue-router";
import {computed, ref} from "vue";
import {useAppStore} from "@/store/app";
import {useUserStore} from "@/store/user";

const username = defineModel<string>("username");
const password = defineModel<string>("password");

const router = useRouter();

const disableButtons = computed(() => {
  return (username.value?.length ?? 0) < 1 || (password.value?.length ?? 0) < 8;
});

const appStore = useAppStore();

const loadingBtn = ref('');
const userStore = useUserStore();
async function login() {
  loadingBtn.value = 'login';
  await appStore.login(username.value!, password.value!);
  loadingBtn.value = '';

  if(appStore.expiresAt * 1000 > Date.now()) {
    await router.push(userStore.connectedToSpotify ? '/' : '/connect' );
  } else {
    console.log('Login failed');
  }
}

async function onEnter() {
  if (!disableButtons.value) {
    await login();
  }
}
</script>

<template>
  <v-form>
    <v-row>
      <v-col cols="12">
        <v-text-field label="Username" v-model="username" hide-details @keydown.enter="onEnter" />
      </v-col>
      <v-col cols="12">
        <v-text-field label="Password" v-model="password" type="password" hide-details @keydown.enter="onEnter" />
      </v-col>
      <v-col cols="12" md="12">
        <v-btn block
               color="primary"
               :disabled="disableButtons"
               @click="login" :loading="loadingBtn === 'login'"
               text="Sign in"
        />
      </v-col>
    </v-row>
  </v-form>
</template>

<style scoped>

</style>
