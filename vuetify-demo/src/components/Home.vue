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
  data: () => ({
    usrs: []
  }),
  methods: {
  }
};
</script>
