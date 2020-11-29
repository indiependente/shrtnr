import Vue from 'vue'
import App from './App.vue'
import VueParticles from 'vue-particles'
import VueClipboard from 'vue-clipboard2'
import VueDarkMode from '@vue-a11y/dark-mode'
import { BootstrapVue, BootstrapVueIcons } from 'bootstrap-vue'
Vue.config.productionTip = false;
VueClipboard.config.autoSetContainer = true;

Vue.use(VueDarkMode);
Vue.use(VueParticles);
Vue.use(VueClipboard);
Vue.use(BootstrapVue)
Vue.use(BootstrapVueIcons)

new Vue({
  render: h => h(App)
}).$mount('#app');
