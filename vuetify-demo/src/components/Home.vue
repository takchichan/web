<template>
  <v-app id="inspire">
    <v-card class="mx-auto" max-width="500">
      <v-toolbar color="pink" dark>
        <v-app-bar-nav-icon></v-app-bar-nav-icon>

        <v-toolbar-title>Inbox</v-toolbar-title>

        <v-spacer></v-spacer>

        <v-btn icon>
          <v-icon>mdi-magnify</v-icon>
        </v-btn>

        <v-btn icon>
          <v-icon>mdi-checkbox-marked-circle</v-icon>
        </v-btn>
      </v-toolbar>

      <v-list two-line>
        <v-list-item-group active-class="pink--text" multiple>
          <template v-for="(item, index) in usrs">
            <v-list-item :key="item.nickName" two-line @click="clickUsr(item)">
              <template>
                <v-list-item-content>
                  <v-list-item-title v-text="item.nickName"></v-list-item-title>
                  <v-list-item-subtitle v-text="item.location"></v-list-item-subtitle>
                </v-list-item-content>
              </template>
            </v-list-item>

            <v-divider v-if="index < usrs.length - 1" :key="index"></v-divider>
          </template>
        </v-list-item-group>
      </v-list>
    </v-card>
  </v-app>
</template>

<script>
export default {
  name: "Home",
  created() {
    this.timerId = setInterval(this.onGetChatStatus, 10000);
    this.onGetUsrList();
  },
  destroyed() {
    clearInterval(this.timerId);
  },
  data: () => ({
    timerId: -1,
    usrs: []
  }),
  methods: {
    onGetChatStatus() {
      let params = {
        nickName: "唐师傅",
        location: "混沌星球",
      };

      this.$httpUtil.get(
        this.$httpUtil.uri.chat,
        params,
        this.onGetChatStatusSucc,
        this.onGetChatStatusFail
      );
    },
    onGetChatStatusSucc(data, params) {
      // 如果是被邀请，则将offer填入signal函数，产生answer
      console.log(data, params);

      // ChatStateWaitForSendingAnswer
      if(data.Initiator == false && data.State == 4) {
        this.$router.push({path: ('/dialog?peer=' + data.Peer + "&initiative=false&signal=" + data.Offer)});
      }
    },
    onGetChatStatusFail(error, params) {
      console.log(error, params);
    },
    onGetUsrList() {
      let params = {};

      this.$httpUtil.get(
        this.$httpUtil.uri.usr,
        params,
        this.onGetUsrListSucc,
        this.onGetUsrListFail
      );
    },
    onGetUsrListSucc(data, params) {
      console.log(data, params);
      this.usrs = data;
    },
    onGetUsrListFail(error, params) {
      console.log(error, params);
    },
    clickUsr(usr) {
      console.log(usr);
      this.$router.push({path: ('/dialog?peer=' + usr.nickName + "&initiative=true")});
    }
  }
};
</script>
