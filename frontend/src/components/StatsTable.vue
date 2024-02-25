<script setup lang="ts">
import {secondsToString} from "@/utils/secondsToString";
import {Stat} from "@/types/Stat";
import {PropType} from "vue";

const props = defineProps({
  name: {
    required: true,
    type: String
  },
  stats: {
    required: true,
    type: Array as PropType<Array<Stat>>
  }
})
</script>

<template>
  <v-table>
    <thead>
    <tr>
      <th></th>
      <th>{{ props.name }}</th>
      <th>Time</th>
      <th class="text-center">
        <v-icon icon="mdi-spotify" color="green" />
      </th>
    </tr>
    </thead>
    <tbody>
    <tr v-for="(stat, index) in props.stats" :key="props.name + (stat.name !== '' ? stat.name : index.toString())">
      <td><span v-if="stat.name !== ''">{{index + 1}}</span></td>
      <td>{{stat.name}}</td>
      <td>{{ secondsToString(stat.datapointCount * 10)}}</td>
      <td class="text-center">
        <a :href="stat.spotifyUrl" target="_blank" v-if="stat.spotifyUrl !== ''">
        open
      </a>
      </td>
    </tr>
    </tbody>
  </v-table>
</template>

<style scoped>

</style>
