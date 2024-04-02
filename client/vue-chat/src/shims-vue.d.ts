/* eslint-disable */

declare module "*.vue" {
  import type { DefineComponent, ComponentCustomProperties } from "vue";
  import { Store } from "vuex";
  const component: DefineComponent<{}, {}, any>;
  export default component;

  interface ComponentCustomProperties {
    $store: Store<any>;
  }
}
