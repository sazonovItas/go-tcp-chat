<template>
  <v-app />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import vApp from "./components/v-app.vue";

export default defineComponent({
  methods: {},
  components: {
    vApp,
  },
});
</script>

<style>
#app {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans",
    Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
}

html {
  background-color: #0f1828;
}
</style>
