<template>
  <div>
    <v-row no-gutters class="peerDialog" v-if="msg.isPeer == true">
      <v-col cols="8">{{msg.text}}</v-col>
      <v-col cols="4"></v-col>
    </v-row>
    <v-row no-gutters class="dialog" v-else>
      <v-col cols="4"></v-col>
      <v-col cols="8">{{msg.text}}</v-col>
    </v-row>
  </div>
</template>

<style scoped>
.dialog {
  margin: 20px 8px;
  text-align: right;
}

.peerDialog {
  margin: 20px 8px;
}
</style>

<script>
export default {
  name: "DialogDetail",
  props: ["msg"],
  data: () => ({})
};
</script>
