<template>
  <footer class="v-footer">
    <div
      v-if="connected"
      style="background-color: #50ff50"
      class="v-footer-indicator"
    ></div>
    <div
      v-else
      style="background-color: #ff5050"
      class="v-footer-indicator"
    ></div>
    <span>{{ host + " : " + port }}</span>
  </footer>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    host: String,
    port: String,
    connected: Boolean,
  },
});
</script>

<style scoped>
.v-footer {
  margin: 0 5px 10px 0;
  width: 100%;

  display: flex;

  background-color: #0d2e69;

  color: #ffffff;
  font-size: 14px;
}

.v-footer-indicator {
  margin: 0 5px 0px 5px;
  width: 8px;
  height: 8px;

  align-self: center;
}
</style>
