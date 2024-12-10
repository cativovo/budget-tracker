<script setup lang="ts">
import { onMounted, ref } from 'vue'

const list = ref<string[]>([])

onMounted(async () => {
  const res = await fetch('/api/entries')
  const j: string[] = await res.json()
  list.value = j
})
</script>

<template>
  <div class="about">
    <h1>This is an about page</h1>
    <ul>
      <li v-for="v in list" :key="v">
        {{ v }}
      </li>
    </ul>
  </div>
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
