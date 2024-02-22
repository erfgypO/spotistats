<script setup lang="ts">
import PageTitleCol from "@/components/PageTitleCol.vue";
import {useUserStore} from "@/store/user";
import {computed, onMounted, ref} from "vue";

const userStore = useUserStore();

onMounted(() => {
  userStore.fetchUser();
});

const newPassword = defineModel('newPassword', { default: '' });
const newPasswordRepeat = defineModel('newPasswordRepeat', { default: '' });

const disableUpdatePasswordBtn = computed(() => {
  return newPassword.value.length < 8 || newPasswordRepeat.value !== newPassword.value;
});

const loading = ref(false);
async function updatePassword() {
  try {
    loading.value = true;
    await userStore.updatePassword(newPassword.value);
  } finally {
    loading.value = false;
    newPassword.value = '';
    newPasswordRepeat.value = '';
  }
}
</script>

<template>
  <v-container fluid>
    <v-row>
      <page-title-col title="Account" :md-width="12" />
      <v-col cols="12" md="6">
        <v-sheet class="pa-3" rounded>
          <v-row>
            <v-col cols="12" class="pb-0">
              <span class="text-h4">User</span>
            </v-col>
            <v-col cols="12" lg="6">
              <v-text-field readonly v-model="userStore.username" label="Soundstats username" hide-details />
            </v-col>
            <v-col cols="12" lg="6">
              <v-text-field readonly v-model="userStore.displayName" label="Spotify username" hide-details />
            </v-col>
            <v-col cols="12" lg="6" style="height: 60px">
              <v-checkbox readonly v-model="userStore.connectedToSpotify" label="Connected to spotify"
                          hide-details color="green" />
            </v-col>
          </v-row>
        </v-sheet>
      </v-col>
      <v-col cols="12" md="6">
        <v-sheet class="pa-3" rounded>
          <v-row>
            <v-col cols="12" class="pb-0">
              <span class="text-h4">Update password</span>
            </v-col>
            <v-col cols="12" lg="6">
              <v-text-field label="New password" hide-details v-model="newPassword" type="password" />
            </v-col>
            <v-col cols="12" lg="6">
              <v-text-field label="Repeat new password" hide-details v-model="newPasswordRepeat" type="password" />
            </v-col>
            <v-col cols="12" lg="6">
              <v-btn color="primary" @click="updatePassword" block text="Update password" :disabled="disableUpdatePasswordBtn" />
            </v-col>
          </v-row>
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>

</style>
