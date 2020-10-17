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
      <div id="content" style="min-height: 100%; max-height: 100%;">
        <v-row no-gutters class="dialog">
          <v-col
            cols="8"
          >.col-12 .col-sm-6 .col-md-8.col-12 .col-sm-6 .col-md-8.col-12 .col-sm-6 .col-md-8.col-12 .col-sm-6 .col-md-8.col-12 .col-sm-6 .col-md-8</v-col>
          <v-col cols="4"></v-col>
        </v-row>
        <v-row no-gutters class="dialog">
          <v-col cols="4"></v-col>
          <v-col
            cols="8"
          >.col-sm-6 .col-md-8.col-12 .col-sm-6 .col-md-8.col-12 .col-sm-6 .col-md-8.col-12 .col-sm-6 .col-md-8.col-12 .col-sm-6 .col-md-8</v-col>
        </v-row>
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
          <v-icon size="24px">{{ icon }}</v-icon>
        </template>
      </v-text-field>
    </v-footer>
  </v-app>
</template>

<style scoped>
.dialog {
  margin: 20px 8px;
}
</style>

<script>
export default {
  name: "App",
  data: () => ({
    message: "Hey!",
    icon: "mdi-facebook"
  })
};
</script>
