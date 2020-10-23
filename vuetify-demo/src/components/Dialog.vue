<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <v-btn href="https://github.com/vuetifyjs/vuetify/releases/latest" target="_blank" text>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
      <v-spacer></v-spacer>
      <span class="mr-2">Latest Release</span>
    </v-app-bar>
    <v-main>
      <div id="content" class="scroll-y" style="min-height: 100%; max-height: 100%;">
        <DialogDetail v-for="(value, index) in messages" :key="index" :msg="value"></DialogDetail>
      </div>
    </v-main>
    <v-footer dark color="primary" app>
      <v-text-field
        v-model="message"
        outlined
        clearable
        label="Message"
        type="text"
        hide-details="true"
      >
        <template v-slot:append>
          <v-icon size="24px" @click="sendMsg">{{ icon }}</v-icon>
        </template>
      </v-text-field>
    </v-footer>
  </v-app>
</template>

<script>
import DialogDetail from "./DialogDetail.vue";
import Peer from "simple-peer";

export default {
  name: "Dialog",
  components: {
    DialogDetail: DialogDetail
  },
  destroyed() {
    clearInterval(this.timerId);
  },
  created() {
    if (
      typeof this.$route.query.peer != "undefined" &&
      typeof this.$route.query.initiative != "undefined"
    ) {
      this.peer = this.$route.query.peer;
      this.initiator = this.$route.query.initiative == "true";
      let signal = this.$route.query.signal;

      console.log("initiator: ", this.initiator);

      // offer
      this.channel = new Peer({
        initiator: this.initiator,
        trickle: false,
      });

      this.channel._debug = console.log

      // exchange signal
      this.channel.on("signal", data => {
        // when host has signaling data, give it to peer somehow
        console.log("signal: ", data);
        // offer answer and ice candidate

        if (
          data.type != "undefined" &&
          (data.type == "offer" || data.type == "answer")
        ) {
          let act = 1;

          if (data.type == "answer") {
            act = 2;
          }

          this.onSendSignal(act, this.initiator, this.peer, data);
          this.timerId = setInterval(this.onGetChatStatus, 5000);
        }
      });

      this.channel.on("connect", () => {
        // wait for 'connect' event before using the data channel
        console.log("connect to peer successfully");
        // backend mark as connected
        this.onSendSignal(5, this.initiator, this.peer, "");
        this.channel.send('whatever' + Math.random());
      });

      this.channel.on('data', data => {
        console.log('data: ' + data);
        let msg = {
          isPeer: true,
          text: data
        };

        this.messages.push(msg);
      })
      // generate answer signal by offer signal
      if (this.initiator == false) {
        console.log("generate answer signal", JSON.parse(signal));
        this.channel.signal(JSON.parse(signal));
      }

      this.channel.on('error', (err) => {
        console.log("channel error: ", err);
      })
    }
  },
  data: () => ({
    timerId: -1,
    peer: null,
    channel: null,
    initiator: false,
    message: "Hey!",
    icon: "mdi-facebook",
    messages: [
      {
        isPeer: false,
        text: "hey! how are you?"
      },
      {
        isPeer: true,
        text: "fine, and you?"
      },
      {
        isPeer: false,
        text: "666"
      }
    ]
  }),
  methods: {
    sendMsg() {
      let msg = {
        isPeer: false,
        text: this.message
      };

      this.messages.push(msg);
      this.channel.send(this.message);
    },
    onSendSignal(act, initiator, peer, signal) {
      let params = {
        action: act,
        peer: peer,
        initiator: initiator,
        signal: JSON.stringify(signal)
      };

      this.$httpUtil.post(
        this.$httpUtil.uri.chat,
        params,
        this.onSendSignalSucc,
        this.onSendSignalFail
      );
    },
    onSendSignalSucc(data, params) {
      console.log(data, params);
    },
    onSendSignalFail(error, params) {
      console.log(error, params);
    },
    onGetChatStatus() {
      let params = {};

      this.$httpUtil.get(
        this.$httpUtil.uri.chat,
        params,
        this.onGetChatStatusSucc,
        this.onGetChatStatusFail
      );
    },
    onGetChatStatusSucc(data, params) {
      console.log("onGetChatStatusSucc", data, params);

      // ChatStateReceivedAnswer
      if (data.Initiator == true && data.State == 3) {
        // using answer signal
        console.log("using answer signal:", data.Answer);
        this.channel.signal(JSON.parse(data.Answer));
      }

      // ChatStateConnected
      if (data.State == 6) {
        clearInterval(this.timerId);
      }
    },
    onGetChatStatusFail(error, params) {
      console.log(error, params);
    }
  }
};
</script>
