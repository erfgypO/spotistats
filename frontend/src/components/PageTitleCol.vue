<script setup lang="ts">
const props = defineProps({
  title: {
    required: true,
    type: String
  },
  mdWidth: {
    default: 6,
    type: Number
  }
})
</script>

<template>
  <v-col cols="12" :md="mdWidth" class="d-flex align-center" v-if="$vuetify.display.mdAndUp">
    <div class="text-h3">{{ props.title }}</div>
  </v-col>
</template>
