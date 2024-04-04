<template>
  <div class="v-chat-wrapper">
    <div class="v-chat-container">
      <div class="v-chat-container-header">
        <button @click="appendMessage">add message</button>
      </div>
      <div class="v-chat-container-messages">
        <div v-for="(m, idx) in messages" :key="'m-' + idx" style="clear: both">
          <div :class="{
          'msg-from-me': m.sender_id == user.id,
          'msg-from-other': m.from == 'other',
        }">
            {{ m.message }}
          </div>
        </div>
      </div>
      <div class="v-chat-container-input">
        <textarea type="text" wrap="soft" class="send-input" placeholder="Type a message" v-model="messageToSend"
          @keyup.enter="send_message" />
      </div>
    </div>
    <vFooter :host="wssock?.host" :port="wssock?.port" :connected="connection_ready" />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { useStore } from "vuex";

import vFooter from "../components/layouts/v-footer.vue";

import WSsocket from "../lib/socket/wssocket";

export default defineComponent({
  setup() {
    const store = useStore();
    const user = store.state.user;
    const connection_ready = ref(false);
    const wssock = new WSsocket(
      store.state.host,
      store.state.port,
      () => {
        connection_ready.value = true;
      },
      undefined,
      () => {
        connection_ready.value = false;
      }
    );
    wssock.connectSocket();
    wssock.setTimeout(store.state.retryTimeout, () => {
      wssock.connectSocket();
    });
    const messageToSend = ref("");

    return {
      store: store,
      user: user,
      wssock: wssock,
      messageToSend: messageToSend,
      connection_ready: connection_ready,
      messages: store.state.messages,
    };
  },
  components: {
    vFooter,
  },
  watch: {},
  computed: {
    connectedSock(): boolean {
      return this.wssock?.connected;
    },
  },
  methods: {
    send_message() {
      this.store.commit("appendMessage", {
        guid: "145324140slfhalfhj",
        sender_id: 1,
        convesation_id: 1,
        message_kind: 1,
        message: this.messageToSend,
        created_at: Date.now(),
        updated_at: Date.now(),
      });
      this.messageToSend = "";
    },
  },
});
</script>

<style>
.v-chat-wrapper {
  width: 100%;
  height: 100%;

  display: flex;
  flex-direction: column;
  justify-content: center;
}

.v-chat-container {
  width: 100%;
  height: 100%;

  background-image: url("../assets/background.webp");
  background: #152033;
  border-radius: 20px;
}

.v-chat-container-header {
  background: #202c33;

  width: 100%;
  height: 12%;

  border-radius: 15px;
}

.v-chat-container-messages {
  width: 100%;
  height: 80%;

  overflow-y: scroll;
  background-size: cover;
}

.msg-from-me {
  border-radius: 7.5px;
  max-width: 65%;
  font-size: 16px;
  line-height: 19px;
  color: #e9edef;
  background: #046a62;
  padding: 5px;
  margin: 20px 20px 5px 0px;

  float: right;
}

.msg-from-other {
  border-radius: 7.5px;
  max-width: 65%;
  font-size: 16px;
  line-height: 19px;
  color: #e9edef;
  background: fade(#202c33, 90%);
  padding: 5px;
  margin: 20px 0px 5px 20px;
  float: left;
}

.v-chat-container-input {
  width: 100%;
  height: 8%;
  resize: none;

  background: #202c33;
  border-radius: 15px;
  display: flex;
}

.send-input {
  padding: 9px 12px 11px;
  margin: 5px 10px;
  border: 1px solid #2a3942;
  background: #2a3942;
  border-radius: 8px;
  font-size: 18px;
  flex-grow: 1;

  color: white;
}

.v-footer {
  width: 100%;
}
</style>
