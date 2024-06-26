<template>
  <div class="v-chat-wrapper">
    <div class="v-chat-container">
      <div class="v-chat-container-header">
        <div class="v-user-description">
          <vIcon
            width="36px"
            height="36px"
            font_size="24px"
            :title="user.name"
            :color="user.color"
          />
          <span :style="{ color: '#e9edef', padding: '6px 0 0 10px' }">{{
            user.name
          }}</span>
          <button class="v-btn-log-out" @click="log_out">Log out</button>
        </div>
      </div>
      <div class="v-chat-container-messages" @scroll="onScroll">
        <VueEternalLoading :load="load" class="v-loader" />
        <ul class="messages">
          <li
            class="message"
            v-for="(m, idx) in getMessages"
            :key="'m-' + idx"
            style="clear: both"
          >
            <div
              :class="{
                'msg-from-me': m.sender_id === user.id,
                'msg-from-other': m.sender_id !== user.id,
              }"
            >
              <div
                :style="{
                  color: members?.get(m.sender_id)?.color,
                  margin: '0 0 5px 0',
                }"
              >
                {{ members?.get(m.sender_id)?.name }}
              </div>
              <div>
                {{ m.message }}
              </div>
              <div :style="{ color: '#bbbbbb', margin: '5px 0 0 0' }">
                {{ new Date(m.created_at) }}
              </div>
            </div>
          </li>
        </ul>
        <div ref="bottomMessage"></div>
      </div>
      <div class="v-chat-container-input">
        <textarea
          type="text"
          wrap="soft"
          class="send-input"
          placeholder="Type a message"
          v-model="messageToSend"
          maxlength="256"
          @keyup.enter="send_message"
        />
      </div>
    </div>
    <vFooter
      :host="store.state.host"
      :port="`${store.state.port}`"
      :connected="connection_ready"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { useStore } from "vuex";

import vFooter from "../components/layouts/v-footer.vue";
import vIcon from "../components/v-icon.vue";
import { VueEternalLoading, LoadAction } from "@ts-pro/vue-eternal-loading";

import WSsocket from "../lib/socket/wssocket";
import { TSMap } from "typescript-map";
import { Connect } from "../lib/reqresp-conn/retry_conn";
import { ResponseToast, NotifySystem } from "../lib/toasts/notifications";
import { successResponse, unauthResponse } from "../lib/reqresp-conn/reqresp";
import { IMessage } from "../store/models/message";
import { IPublicUser } from "../store/models/user";
import { Request } from "../lib/reqresp-conn/conn";
import {
  IMessagesRequest,
  IMessagesResponse,
  messagesEndpoint,
  memberEndpoint,
} from "../store/endpoints/endpoints";
import { IEvent } from "../store/models/event";

const MESSAGES_CNT = 25;

export default defineComponent({
  setup() {
    const store = useStore();
    const user = store.state.user;
    const connection_ready = ref(false);
    const messageToSend = ref("");

    const wssock = new WSsocket(store.state.host, store.state.port);
    const members = ref(new Map<number, IPublicUser>([]));
    const messages = ref(store.state.messages);

    return {
      store: store,
      user: user,
      wssock: wssock,
      messageToSend: messageToSend,
      connection_ready: connection_ready,
      members: members,
      messages: messages,
      detach_scroll: true,
    };
  },
  mounted() {
    this.try_connect();
    const retryConnection = () => {
      this.try_connect();
      if (this.wssock !== undefined) {
        setTimeout(() => {
          retryConnection();
        }, this.store.state.retryTimeout);
      }
    };
    setTimeout(() => {
      retryConnection();
    }, this.store.state.retryTimeout);

    Request(
      this.store.state.host,
      this.store.state.port,
      this.store.state.requestTimeout,
      {
        method: "GET",
        url: memberEndpoint,
        proto: "http",

        header: new TSMap([["Content-Type", "application/json"]]),
        body: "",
      }
    )
      .then((value) => {
        ResponseToast.notify(value.status_code, value.status);

        if (successResponse(value)) {
          try {
            const members: IPublicUser[] = JSON.parse(value.body);
            this.members = new Map<number, IPublicUser>(
              members.map((member: IPublicUser) => [member.id, member])
            );
          } catch (e) {
            console.log(e);
          }
        }
      })
      .catch((e) => {
        console.log(e);
      });
    this.scroll_to_end();
  },
  components: {
    vIcon,
    vFooter,
    VueEternalLoading,
  },
  computed: {
    getMessages(): IMessage[] {
      return this.messages;
    },
  },
  methods: {
    load: async function ({ loaded }: LoadAction) {
      console.log("...loading");

      const request: IMessagesRequest = {
        auth_token: this.store.state.token,
        timestamp:
          this.messages.length > 0
            ? this.messages[0].created_at
            : new Date(Date.now()),
        limit: MESSAGES_CNT,
      };
      const response = await Request(
        this.store.state.host,
        this.store.state.port,
        this.store.state.requestTimeout,
        {
          method: "GET",
          url: messagesEndpoint,
          proto: "http",

          header: new TSMap([["Content-Type", "application/json"]]),
          body: JSON.stringify(request),
        }
      );

      if (unauthResponse(response)) {
        this.log_out();
        return;
      }

      if (!successResponse(response)) {
        return;
      }

      let receivedMessages: IMessagesResponse;
      try {
        receivedMessages = JSON.parse(response.body);
        NotifySystem.notify("info", response.body);
        this.store.commit("unshiftMessages", receivedMessages.messages);
      } catch (e) {
        console.log(e);
        // loaded(MESSAGES_CNT, MESSAGES_CNT);
        return;
      }
      if (this.detach_scroll) {
        setTimeout(() => {
          this.scroll_to_end();
        }, 50);
      }

      loaded(
        receivedMessages.messages?.length === undefined
          ? 0
          : receivedMessages.messages.length,
        MESSAGES_CNT
      );
    },
    send_message() {
      if (this.messageToSend.trim() === "") {
        return;
      }

      if (!this.connection_ready || !this.wssock) {
        NotifySystem.notify(
          "warning",
          "cannot send message: disconnected from the server"
        );
        return;
      }

      const msg: IEvent = {
        type: "NewMessageEvent",
        payload: {
          id: "",
          sender_id: this.user.id,
          message_kind: 2,
          message: this.messageToSend.trim(),
          created_at: new Date(Date.now()),
          updated_at: new Date(Date.now()),
        },
      };

      try {
        this.wssock.socketSend(JSON.stringify(msg), 0x1, true);
      } catch (e) {
        console.log(e);
      }

      this.messageToSend = "";
    },
    onScroll({ target: { scrollTop, clientHeight, scrollHeight } }) {
      if (scrollTop + clientHeight >= scrollHeight - 30) {
        this.detach_scroll = true;
      } else {
        this.detach_scroll = false;
      }
    },
    try_connect() {
      if (this.wssock?.try_connecting || this.wssock?.connected) {
        return;
      }
      this.connection_ready = false;

      Connect(
        this.wssock,
        this.store.state.token,
        this.store.state.retryTimeout,

        () => {
          this.connection_ready = true;
        },
        (data: Buffer) => {
          try {
            const msg: IEvent = JSON.parse(data.toString());
            this.store.commit("appendMessage", msg.payload);
            if (this.detach_scroll) {
              setTimeout(() => {
                this.scroll_to_end();
              }, 50);
            }

            if (!this.members?.get(msg.payload.sender_id)) {
              Request(
                this.store.state.host,
                this.store.state.port,
                this.store.state.requestTimeout,
                {
                  method: "GET",
                  url: memberEndpoint + "/" + msg.payload.sender_id,
                  proto: "http",

                  header: new TSMap([["Content-Type", "application/json"]]),
                  body: "",
                }
              )
                .then((value) => {
                  ResponseToast.notify(value.status_code, value.status);

                  if (successResponse(value)) {
                    try {
                      const member: IPublicUser = JSON.parse(value.body);
                      this.members.set(member.id, member);
                    } catch (e) {
                      console.log(e);
                    }
                  }
                })
                .catch((e) => {
                  console.log(e);
                });
            }
          } catch (e) {
            console.log(e);
          }
        },
        () => {
          this.connection_ready = false;
          NotifySystem.notify("warning", "disconnected");
        },
        (e: Error) => {
          NotifySystem.notify("error", e.toString());
        }
      )
        .then((value) => {
          ResponseToast.notify(value.status_code, value.status);
          if (successResponse(value)) {
            NotifySystem.notify("success", "connected to server");
            this.connection_ready = true;
          } else if (unauthResponse(value)) {
            this.log_out();
          } else {
            NotifySystem.notify("warning", "cannot connect to server");
          }
        })
        .catch((e) => {
          console.log(e);
        });
    },
    scroll_to_end() {
      this.$refs.bottomMessage.scrollIntoView({ behaivor: "smooth" });
    },
    log_out() {
      this.store.commit("logOut");

      try {
        this.wssock.closeSocket();
      } catch (e) {
        console.log(e);
      }

      this.connection_ready = false;
      this.wssock = undefined;
      this.$router.push("/");
    },
  },
});
</script>

<style>
.v-loader {
  text-align: center;
  color: #e9edef;
}

.v-btn-log-out {
  width: 75px;
  margin: 0 0 0 10px;

  background: #046a62;

  color: #e9edef;
  border: none;
  border-radius: 15px;
}

.v-btn-log-out:hover {
  background: #152033;
}

.v-chat-wrapper {
  width: 100%;
  height: 100%;

  display: flex;
  flex-direction: column;
  justify-content: center;
}

.v-chat-container {
  width: 100%;
  height: 98%;

  background: #152033;
  border-radius: 20px;
}

.v-chat-container-header {
  background: #202c33;

  width: 100%;
  height: 8%;

  border-radius: 15px;
}

.v-user-description {
  padding: 20px;

  display: flex;
  flex-wrap: wrap;
}

.v-chat-container-messages {
  padding: 10px 0 10px 0;

  width: 100%;
  height: 80%;

  position: sticky;
  overflow-y: scroll;
  overflow-x: hidden;
  background-size: cover;
}

.msg-from-me {
  border-radius: 7.5px;
  max-width: 40%;
  font-size: 16px;
  line-height: 19px;
  color: #e9edef;
  background: #046a62;
  padding: 5px;
  margin: 20px 20px 10px 0px;
  float: right;

  word-wrap: break-word;
}

.msg-from-other {
  padding: 15px;

  border-radius: 7.5px;
  max-width: 40%;
  font-size: 16px;
  line-height: 19px;
  color: #e9edef;
  background: #202c33;
  padding: 5px;
  margin: 20px 0px 10px 20px;
  float: left;

  word-wrap: break-word;
}

.v-chat-container-input {
  margin: 10px 0 0 0;
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

/* width */
::-webkit-scrollbar {
  width: 5px;
}

/* Track */
::-webkit-scrollbar-track {
  background: #202c33;
}

/* Handle */
::-webkit-scrollbar-thumb {
  background: #2a3942;
}

/* Handle on hover */
::-webkit-scrollbar-thumb:hover {
  background: #555;
}

.v-footer {
  width: 100%;
}
</style>
