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
import Dialog from "./components/Dialog.vue";

export default {
  name: "App",
  components: {
    Dialog: Dialog
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
