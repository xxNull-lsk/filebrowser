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

        <router-link class="action" to="/files/.trash" :aria-label="$t('sidebar.trash')" :title="$t('sidebar.trash')">
          <i class="material-icons">delete</i>
          <span>{{ $t('sidebar.trash') }}</span>
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
        <span class="favorite_title">{{$t('sidebar.favorite')}}</span>
        <div v-for="favorite in favorites" class="favorite" v-bind:key="favorite">
          <router-link class="favorite_left" :key="favorite.hash" :to="'/files' + favorite.path" :aria-label="favorite.name" :title="favorite.name">
            <div class="favorite_content">
              <i class="material-icons">{{ favorite.type }}</i>
              <span>{{ favorite.name }}</span>
            </div>
          </router-link>
          
          <div class="favorite_right">
            <button class="delete" @click="deleteFavorite(favorite.hash)" id="deleteFavorite">
              <i class="material-icons">remove_circle_outline</i>
            </button>
          </div>
        </div>
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

    <div class="credits">
      <span>
        <span v-if="disableExternal">File Browser</span>
        <a v-else rel="noopener noreferrer" target="_blank" href="https://github.com/xxnull-lsk/filebrowser">File Browser</a>
        <span> {{ version }}</span>
      </span>
      <span><a @click="help">{{ $t('sidebar.help') }}</a></span>
    </div>
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
    this.getFavorites()

    this.$root.$on("favorite-created", this.getFavorites);
  },
  beforeDestroy() {
    this.$root.$off("favorite-created", this.deleted);
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
    deleteFavorite(hash){
      favorite_api.remove(hash).then(this.getFavorites)
    },
    async getFavorites() {
      try {
        const favorites = await favorite_api.list();
        this.favorites = favorites;
      } catch (e) {
        console.error(e)
      }
    },
    help () {
      this.$store.commit('showHover', 'help')
    },
    logout: auth.logout
  }
}
</script>
