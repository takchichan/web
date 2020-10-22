import Vue from "vue";
import VueRouter from "vue-router";

import Home from "./components/Home.vue";
import Dialog from "./components/Dialog.vue";

Vue.use(VueRouter);

const routes = [
    {
        path: "/",
        component: Home
    },
    {
        path: "/dialog",
        component: Dialog
    }
]

var router = new VueRouter({
    routes
})
export default router;