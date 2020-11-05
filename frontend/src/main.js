import Vue from 'vue'
import App from './App.vue'
import VueParticles from 'vue-particles'
import VueClipboard from 'vue-clipboard2'
import VueDarkMode from '@vue-a11y/dark-mode'

Vue.config.productionTip = false;
VueClipboard.config.autoSetContainer = true; // add this line

Vue.use(VueDarkMode);
Vue.use(VueParticles);
Vue.use(VueClipboard);

new Vue({
  render: h => h(App)
}).$mount('#app');
