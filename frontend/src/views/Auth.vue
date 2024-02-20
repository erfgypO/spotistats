<script setup lang="ts">
import SignInDialog from "@/components/SignInDialog.vue";
import SignUpDialog from "@/components/SignUpDialog.vue";
const tab = defineModel<number>("tab");
</script>

<template>
  <v-card max-width="750px" class="mx-auto">
    <!--<v-card-title>Login</v-card-title>-->
    <v-card-text>
      <v-tabs align-tabs="center" color="primary" v-model="tab">
        <v-tab :value="1">Sign in</v-tab>
        <v-tab :value="2">Sign up</v-tab>
      </v-tabs>
      <v-window v-model="tab">
        <v-window-item :value="1">
          <sign-in-dialog />
        </v-window-item>
        <v-window-item :value="2">
          <sign-up-dialog />
        </v-window-item>
      </v-window>
    </v-card-text>
  </v-card>
</template>

<style scoped>

</style>
