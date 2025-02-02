import { sync } from 'vuex-router-sync'
import store from '@/store'
import router from '@/router'
import i18n from '@/i18n'
import Vue from 'vue'
import { recaptcha, loginPage } from '@/utils/constants'
import { login, validateLogin } from '@/utils/auth'
import App from '@/App'
import { Select, Option, Checkbox } from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

import '@/utils/vue'

sync(store, router)

async function start () {
  if (loginPage) {
    await validateLogin()
  } else {
    await login('', '', '')
  }

  if (recaptcha) {
    await new Promise (resolve => {
      const check = () => {
        if (typeof window.grecaptcha === 'undefined') {
          setTimeout(check, 100)
        } else {
          resolve()
        }
      }

      check()
    })
  }

  Vue.use(Select);
  Vue.use(Option);
  Vue.use(Checkbox);

  new Vue({
    el: '#app',
    store,
    router,
    i18n,
    template: '<App/>',
    components: { App }
  })
}

start()
