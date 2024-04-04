<template>
  <div class="v-container-icon" :style="style">
    {{ title.toUpperCase()[0] }}
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    width: String,
    height: String,
    font_size: String,

    title: String,
    color: String,
  },
  computed: {
    style(): string {
      return (
        "background-color: " +
        this.color +
        "; width: " +
        this.width +
        "; height: " +
        this.height +
        "; font-size: " +
        this.font_size +
        ";"
      );
    },
  },
});
</script>

<style scoped>
.v-container-icon {
  display: flex;
  flex-direction: column;
  justify-content: center;

  text-align: center;

  color: #ffffff;
  border-radius: 50%;
  font-family: 600;
}
</style>
