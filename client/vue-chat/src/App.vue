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
  width: 100%;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans",
    Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";

  display: flex;
  justify-content: center;
}

* {
  margin: 0;
  padding: 0;
}

textarea {
  resize: none;
}

body {
  width: 100%;
  height: 100%;
}

html {
  width: 100%;
  height: 100%;
  background-color: #0f1828;
}
</style>
