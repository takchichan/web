import Vue from "vue";
import VueRouter from "vue-router";

import DialogPage from "./components/DialogPage.vue";

Vue.use(VueRouter);

const routes = [
    {
        path: "/dialog",
        component: DialogPage
    }
]

var router = new VueRouter({
    routes
})
export default router;