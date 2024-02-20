<script setup lang="ts">
import {computed, ref} from "vue";
import {useAppStore} from "@/store/app";
import {useRouter} from "vue-router";

const username = defineModel<string>("username");
const password = defineModel<string>("password");

const router = useRouter();

const disableButtons = computed(() => {
  return (username.value?.length ?? 0) < 1 || (password.value?.length ?? 0) < 8;
});

const appStore = useAppStore();

const loadingBtn = ref('');

async function login() {
  loadingBtn.value = 'login';
  await appStore.login(username.value!, password.value!);
  loadingBtn.value = '';

  if(appStore.expiresAt * 1000 > Date.now()) {
    await router.push('/');
  } else {
    console.log('Login failed');
  }
}

</script>

<template>
  <v-card max-width="750px" class="mx-auto">
    <v-card-title>Login</v-card-title>
    <v-card-text>
      <v-form>
        <v-row>
          <v-col cols="12">
            <v-text-field label="Username" v-model="username" hide-details />
          </v-col>
          <v-col cols="12">
            <v-text-field label="Password" v-model="password" type="password" hide-details />
          </v-col>
          <v-col cols="12" md="6">
            <v-btn block color="primary" :disabled="disableButtons" @click="login" :loading="loadingBtn === 'login'">Login</v-btn>
          </v-col>
          <v-col cols="12" md="6">
            <v-btn block color="primary" variant="outlined" :disabled="disableButtons">Create Account</v-btn>
          </v-col>
        </v-row>
      </v-form>
    </v-card-text>
  </v-card>
</template>

<style scoped>

</style>
