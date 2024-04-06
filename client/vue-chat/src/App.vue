<template>
  <div class="v-container">
    <router-view></router-view>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  methods: {},
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

.v-container {
  width: 100%;
  height: 100%;
  min-width: 550px;
  max-width: 880px;

  align-self: center;

  display: flex;
  justify-content: center;
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
