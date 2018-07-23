import Vue from "vue";
import ElementUI from "element-ui";
import VueResource from "vue-resource";
import "element-ui/lib/theme-chalk/index.css";
import App from "./App.vue";

Vue.use(ElementUI);
Vue.use(VueResource);

new Vue({
  el: "#app",
  render: h => h(App)
});
