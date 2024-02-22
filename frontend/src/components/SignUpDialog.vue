<script setup lang="ts">
import {computed, ref} from "vue";
import {useAppStore} from "@/store/app";
import {useRouter} from "vue-router";

const appStore = useAppStore();

const username = defineModel<string>("username");
const password = defineModel<string>("password");
const passwordConfirm = defineModel<string>("passwordConfirm");

const disabled = computed(() => {
  return !username.value || !password.value || !passwordConfirm.value || password.value !== passwordConfirm.value;
});

const loadingBtn = ref(false);

const router = useRouter();

async function signUp() {
  loadingBtn.value = true;
  await appStore.signUp(username.value!, password.value!);
  loadingBtn.value = false;

  if(appStore.expiresAt * 1000 > Date.now()) {
    await router.push('/');
  } else {
    console.log('Sign up failed');
  }
}

async function onEnter() {
  if (!disabled.value) {
    await signUp();
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
    <v-col cols="12">
      <v-text-field label="Confirm Password" v-model="passwordConfirm" type="password" hide-details @keydown.enter="onEnter" />
    </v-col>
    <v-col cols="12" md="12">
      <v-btn block
             color="primary"
             :disabled="disabled"
             text="Sign up"
             :loading="loadingBtn"
             @click="signUp"
      />
    </v-col>
  </v-row>
</v-form>
</template>

<style scoped>

</style>
