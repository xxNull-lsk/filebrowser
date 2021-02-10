<template>
  <header v-if="!isEditor && !isPreview">
    <div>
      <button @click="openSidebar" :aria-label="$t('buttons.toggleSidebar')" :title="$t('buttons.toggleSidebar')" class="action">
        <i class="material-icons">menu</i>
      </button>
      <img :src="logoURL" :alt="name">
      <span class="logo-name">{{ name }}</span>
      <search v-if="isLogged"></search>
    </div>
    <div>
      <template v-if="isLogged || isSharing">
        <button v-show="!isSharing" @click="openSearch" :aria-label="$t('buttons.search')" :title="$t('buttons.search')" class="search-button action">
          <i class="material-icons">search</i>
        </button>

        <button @click="openMore" id="more" :aria-label="$t('buttons.more')" :title="$t('buttons.more')" class="action">
          <i class="material-icons">more_vert</i>
        </button>

        <!-- Menu that shows on listing AND mobile when there are files selected -->
        <div id="file-selection" v-if="isMobile && isListing && !isSharing">
          <span v-if="selectedCount > 0">{{ selectedCount }} selected</span>
          <share-button v-show="showShareButton"></share-button>
          <rename-button v-show="showRenameButton"></rename-button>
          <copy-button v-show="showCopyButton"></copy-button>
          <move-button v-show="showMoveButton"></move-button>
          <delete-button v-show="showDeleteButton"></delete-button>
          <restore-button v-show="isTrash"></restore-button>
        </div>

        <!-- This buttons are shown on a dropdown on mobile phones -->
        <div>
        <div id="dropdown" :class="{ active: showMore }">
          <div v-if="!isListing || !isMobile" class="header-actions">
            <share-button v-show="showShareButton"></share-button>
            <rename-button v-show="showRenameButton"></rename-button>
            <copy-button v-show="showCopyButton"></copy-button>
            <move-button v-show="showMoveButton"></move-button>
            <delete-button v-show="showDeleteButton"></delete-button>
            <restore-button v-show="isTrash"></restore-button>
          </div>

          <new-file-button v-show="showNewFileButton"/>
          <new-dir-button v-show="showNewDirButton"/>
          <favorite-button v-show="showFavoriteButton"/>
          <shell-button v-if="isExecEnabled && !isSharing && !isTrash && user.perm.execute" />
          <switch-button v-show="isListing"></switch-button>
          <download-button v-show="showDownloadButton"></download-button>
          <upload-button v-show="showUpload"></upload-button>
          <info-button v-show="isFiles"></info-button>

          <button v-show="isListing || (isSharing && req.isDir)" @click="toggleMultipleSelection" :aria-label="$t('buttons.selectMultiple')" :title="$t('buttons.selectMultiple')" class="action" >
            <i class="material-icons">check_circle</i>
            <span>{{ $t('buttons.select') }}</span>
          </button>
        </div>
        </div>

      </template>

      <div v-show="showOverlay" @click="resetPrompts" class="overlay"></div>
      <div class="user-info" @mouseleave="mouseleave">
        <p class="menu-title" @click="selectorToggle">
            <i class="material-icons">person</i>
            <span class="user-name">{{ this.$store.state.user.userName }} </span>
            <i class="material-icons">arrow_drop_down</i>
        </p>
        <div class="menu-items" v-show="showMenu">

          <router-link class="action" to="/settings" :aria-label="$t('sidebar.settings')" :title="$t('sidebar.settings')">
            <i class="material-icons">settings_applications</i>
            <span>{{ $t('sidebar.settings') }}</span>
          </router-link>

          <button v-if="authMethod == 'json'" @click="logout" class="action" id="logout" :aria-label="$t('sidebar.logout')" :title="$t('sidebar.logout')">
            <i class="material-icons">exit_to_app</i>
            <span>{{ $t('sidebar.logout') }}</span>
          </button>
        </div>
      </div>
      
    </div>
  </header>
</template>

<script>
import Search from './Search'
import InfoButton from './buttons/Info'
import DeleteButton from './buttons/Delete'
import RenameButton from './buttons/Rename'
import UploadButton from './buttons/Upload'
import DownloadButton from './buttons/Download'
import SwitchButton from './buttons/SwitchView'
import MoveButton from './buttons/Move'
import CopyButton from './buttons/Copy'
import ShareButton from './buttons/Share'
import ShellButton from './buttons/Shell'
import NewFileButton from './buttons/NewFile'
import NewDirButton from './buttons/NewDir'
import FavoriteButton from './buttons/Favorite'
import RestoreButton from './buttons/Restore'
import {mapGetters, mapState} from 'vuex'
import { name, logoURL, enableExec } from '@/utils/constants'
import * as api from '@/api'
import buttons from '@/utils/buttons'
import * as auth from '@/utils/auth'
import { authMethod } from '@/utils/constants'

export default {
  name: 'header-layout',
  components: {
    Search,
    InfoButton,
    DeleteButton,
    ShareButton,
    RenameButton,
    DownloadButton,
    CopyButton,
    UploadButton,
    SwitchButton,
    MoveButton,
    ShellButton,
    NewFileButton,
    NewDirButton,
    FavoriteButton,
    RestoreButton
  },
  data: function () {
    return {
      width: window.innerWidth,
      showMenu: false,
      pluginData: {
        api,
        buttons,
        'store': this.$store,
        'router': this.$router
      }
    }
  },
  created () {
    window.addEventListener('resize', () => {
      this.width = window.innerWidth
    })
  },
  computed: {
    ...mapGetters([
      'selectedCount',
      'isFiles',
      'isTrash',
      'isEditor',
      'isPreview',
      'isListing',
      'isLogged',
      'isSharing'
    ]),
    ...mapState([
      'req',
      'user',
      'loading',
      'reload',
      'multiple'
    ]),
    logoURL: () => logoURL,
    name: () => name,
    authMethod: () => authMethod,
    isExecEnabled: () => enableExec,
    isMobile () {
      return this.width <= 736
    },
    showUpload () {
      return !this.isTrash && this.isListing && this.user.perm.create
    },
    showDownloadButton () {
      return (!this.isTrash && this.isFiles && this.user.perm.download) || (this.isSharing && this.selectedCount > 0)
    },
    showDeleteButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount !== 0 && this.user.perm.delete)
        : this.user.perm.delete)
    },
    showRenameButton () {
      return !this.isTrash && this.isFiles && (this.isListing
        ? (this.selectedCount === 1 && this.user.perm.rename)
        : this.user.perm.rename)
    },
    showShareButton () {
      return !this.isTrash && this.isFiles && (this.isListing
        ? (this.selectedCount === 1 && this.user.perm.share)
        : this.user.perm.share)
    },
    showMoveButton () {
      return !this.isTrash && this.isFiles && (this.isListing
        ? (this.selectedCount > 0 && this.user.perm.rename)
        : this.user.perm.rename)
    },
    showCopyButton () {
      return !this.isTrash && this.isFiles && (this.isListing
        ? (this.selectedCount > 0 && this.user.perm.create)
        : this.user.perm.create)
    },
    showNewFileButton () {
      return !this.isTrash && this.isFiles && (this.isListing
        ? (this.selectedCount === 0 && this.user.perm.create)
        : this.user.perm.create)
    },
    showNewDirButton () {
      return !this.isTrash && this.isFiles && (this.isListing
        ? (this.selectedCount === 0 && this.user.perm.create)
        : this.user.perm.create)
    },
    showFavoriteButton () {
      return !this.isTrash && this.isFiles
    },
    showMore () {
      return (!this.isTrash && this.isFiles || this.isSharing) && this.$store.state.show === 'more'
    },
    showOverlay () {
      return this.showMore
    }
  },
  methods: {
    openSidebar () {
      this.$store.commit('showHover', 'sidebar')
    },
    openMore () {
      this.$store.commit('showHover', 'more')
    },
    openSearch () {
      this.$store.commit('showHover', 'search')
    },
    toggleMultipleSelection () {
      this.$store.commit('multiple', !this.multiple)
      this.resetPrompts()
    },
    resetPrompts () {
      this.$store.commit('closeHovers')
    },
    logout: auth.logout,
    selectorToggle(){
        this.showMenu = !this.showMenu
    },
    mouseleave(){
        this.showMenu = false
    }
  }
}
</script>
