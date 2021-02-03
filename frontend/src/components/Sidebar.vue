<template>
  <nav :class="{active}">
    <template v-if="isLogged">

      <div>
        <router-link class="action" to="/files/" :aria-label="$t('sidebar.myFiles')" :title="$t('sidebar.myFiles')">
          <i class="material-icons">folder</i>
          <span>{{ $t('sidebar.myFiles') }}</span>
        </router-link>

        <router-link class="action" to="/settings/shares" :aria-label="$t('settings.shareManagement')" :title="$t('settings.shareManagement')">
          <i class="material-icons">share</i>
          <span>{{ $t('settings.shareManagement') }}</span>
        </router-link>
      </div>

      <div>
        <router-link class="action" to="/settings" :aria-label="$t('sidebar.settings')" :title="$t('sidebar.settings')">
          <i class="material-icons">settings_applications</i>
          <span>{{ $t('sidebar.settings') }}</span>
        </router-link>

        <button v-if="authMethod == 'json'" @click="logout" class="action" id="logout" :aria-label="$t('sidebar.logout')" :title="$t('sidebar.logout')">
          <i class="material-icons">exit_to_app</i>
          <span>{{ $t('sidebar.logout') }}</span>
        </button>
      </div>

      <div>
          <router-link class="favorite" v-for="favorite in favorites" :key="favorite.hash" to="favorite.path" :aria-label="favorite.name" :title="favorite.name">
            <i class="material-icons">{{ favorite.type }}</i>
            <span>{{ favorite.name }}</span>
          </router-link>
      </div>
    </template>
    <template v-else>
      <router-link class="action" to="/login" :aria-label="$t('sidebar.login')" :title="$t('sidebar.login')">
        <i class="material-icons">exit_to_app</i>
        <span>{{ $t('sidebar.login') }}</span>
      </router-link>

      <router-link v-if="signup" class="action" to="/login" :aria-label="$t('sidebar.signup')" :title="$t('sidebar.signup')">
        <i class="material-icons">person_add</i>
        <span>{{ $t('sidebar.signup') }}</span>
      </router-link>
    </template>

    <p class="credits">
      <span>
        <span v-if="disableExternal">File Browser</span>
        <a v-else rel="noopener noreferrer" target="_blank" href="https://github.com/filebrowser/filebrowser">File Browser</a>
        <span> {{ version }}</span>
      </span>
      <span><a @click="help">{{ $t('sidebar.help') }}</a></span>
    </p>
  </nav>
</template>

<script>
import { mapState, mapGetters } from 'vuex'
import { favorites as favorite_api } from "@/api";
import * as auth from '@/utils/auth'
import { version, signup, disableExternal, noAuth, authMethod } from '@/utils/constants'

export default {
  name: 'sidebar',
  data: function () {
    return {
      favorites: [],
    };
  },
  async beforeMount() {
    try {
      const favorites = await favorite_api.list();
      this.favorites = favorites;
    } catch (e) {
      console.error(e)
    }
  },
  computed: {
    ...mapState([ 'user' ]),
    ...mapGetters([ 'isLogged' ]),
    active () {
      return this.$store.state.show === 'sidebar'
    },
    signup: () => signup,
    version: () => version,
    disableExternal: () => disableExternal,
    noAuth: () => noAuth,
    authMethod: () => authMethod
  },
  methods: {
    help () {
      this.$store.commit('showHover', 'help')
    },
    logout: auth.logout
  }
}
</script>
