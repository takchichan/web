import Vue from 'vue';
import App from './App.vue';
import vuetify from './plugins/vuetify';

import router from "./router.js";
import httpUtil from './httpUtil.js';

Vue.config.productionTip = false;

Vue.prototype.$httpUtil = httpUtil;

new Vue({
  vuetify,
  router,
  render: h => h(App)
}).$mount('#app');
