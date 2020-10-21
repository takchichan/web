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
        <Dialog v-for="(value, index) in messages" :key="index" :msg="value"></Dialog>
        <v-layout v-scroll:#content="onScroll" style="height: 1000px;"></v-layout>
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
import Dialog from "./Dialog.vue";
import Peer from "simple-peer";

export default {
  name: "App",
  components: {
    Dialog: Dialog
  },
  created() {
    // offer
    var peer1 = new Peer({ initiator: true });
    // anwser
    var peer2 = new Peer();

    // exchange signal

    peer1.on("signal", data => {
      // when peer1 has signaling data, give it to peer2 somehow
      console.log("peer1", data);
      peer2.signal(data);
    });

    peer2.on("signal", data => {
      // when peer2 has signaling data, give it to peer1 somehow
      console.log("peer2", data);
      peer1.signal(data);
    });

    peer1.on("connect", () => {
      // wait for 'connect' event before using the data channel
      peer1.send("hey peer2, how is it going?");
    });

    peer2.on("data", data => {
      // got a data channel message
      console.log("got a message from peer1: " + data);
    });
  },
  data: () => ({
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
    },
    onScroll(e) {
      console.log(e);
      this.offsetTop = e.target.scrollTop;
    }
  }
};
</script>
