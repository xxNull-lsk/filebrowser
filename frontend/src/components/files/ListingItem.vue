<template>
  <div class="item"
  role="button"
  tabindex="0"
  :draggable="isDraggable"
  @dragstart="dragStart"
  @dragover="dragOver"
  @drop="drop"
  @click="itemClick"
  @dblclick="dblclick"
  @touchstart="touchstart"
  :data-dir="isDir"
  :aria-label="name"
  :aria-selected="isSelected">
    <div>
      <lazy-component v-if="type==='image' && isThumbsEnabled && !isSharing" :timeout="1000">
        <img :src="thumbnailUrl">
      </lazy-component>
      <i v-else-if="isFontIcon()" :class="icon"></i>
      <svg v-else-if="isColorIcon()" class="color-icon" aria-hidden="true">
        <use :xlink:href="icon"></use>
      </svg>
      <i v-else class="material-icons">{{ icon }}</i>
    </div>

    <div>
      <p class="name">{{ name }}</p>

      <p v-if="isDir" class="size" data-order="-1"></p>
      <p v-else class="size" :data-order="humanSize()">{{ humanSize() }}</p>

      <p class="modified">
        <time :datetime="modified" :title="localTime()">{{ humanTime() }}</time>
      </p>
    </div>
  </div>
</template>

<script>
import { baseURL, enableThumbs, iconTheme } from '@/utils/constants'
import { mapMutations, mapGetters, mapState } from 'vuex'
import filesize from 'filesize'
import moment from 'moment'
import { files as api } from '@/api'
import * as upload  from '@/utils/upload'

export default {
  name: 'item',
  data: function () {
    return {
      touches: 0,
      icon: ''
    }
  },
  props: ['name', 'isDir', 'url', 'type', 'size', 'modified', 'index'],
  mounted () {
    if (this.isFontIcon()){
      this.fontIcons()
    } else if (this.isColorIcon()) {
      this.colorIcons()
    } else {
      this.imaterialIcons()
    }
  },
  computed: {
    ...mapState(['user', 'selected', 'req', 'jwt']),
    ...mapGetters(['selectedCount', 'isSharing', 'isTrash']),
    singleClick () {
      if (this.isSharing || this.isTrash) return false
      return this.user.singleClick
    },
    isSelected () {
      return (this.selected.indexOf(this.index) !== -1)
    },
    isDraggable () {
      return !this.isTrash && !this.isSharing && this.user.perm.rename
    },
    canDrop () {
      if (!this.isDir || this.isSharing || this.isTrash) return false

      for (let i of this.selected) {
        if (this.req.items[i].url === this.url) {
          return false
        }
      }

      return true
    },
    thumbnailUrl () {
      const path = this.url.replace(/^\/files\//, '')
      return `${baseURL}/api/preview/thumb/${path}?auth=${this.jwt}&inline=true`
    },
    isThumbsEnabled () {
      return enableThumbs
    },
  },
  methods: {
    ...mapMutations(['addSelected', 'removeSelected', 'resetSelected']),
    humanSize: function () {
      return filesize(this.size)
    },
    localTime: function () {
      var d = new Date(this.modified)
      return d.toLocaleString()
    },
    humanTime: function () {
      return moment(this.modified).fromNow()
    },
    dragStart: function () {
      if (this.selectedCount === 0) {
        this.addSelected(this.index)
        return
      }

      if (!this.isSelected) {
        this.resetSelected()
        this.addSelected(this.index)
      }
    },
    dragOver: function (event) {
      if (!this.canDrop) return

      event.preventDefault()
      let el = event.target

      for (let i = 0; i < 5; i++) {
        if (!el.classList.contains('item')) {
          el = el.parentElement
        }
      }

      el.style.opacity = 1
    },
    drop: async function (event) {
      if (!this.canDrop) return
      event.preventDefault()

      if (this.selectedCount === 0) return

      let el = event.target
      for (let i = 0; i < 5; i++) {
        if (el !== null && !el.classList.contains('item')) {
          el = el.parentElement
        }
      }

      let items = []

      for (let i of this.selected) {
        items.push({
          from: this.req.items[i].url,
          to: this.url + this.req.items[i].name,
          name: this.req.items[i].name
        })
      }      

      let base = el.querySelector('.name').innerHTML + '/'
      let path = this.$route.path + base
      let baseItems = (await api.fetch(path)).items

      let action = (overwrite, rename) => {
        api.move(items, overwrite, rename).then(() => {
          this.$store.commit('setReload', true)
        }).catch(this.$showError)
      }

      let conflict = upload.checkConflict(items, baseItems)

      let overwrite = false
      let rename = false

      if (conflict) {
        this.$store.commit('showHover', {
          prompt: 'replace-rename',
          confirm: (event, option) => {
            overwrite = option == 'overwrite'
            rename = option == 'rename'

            event.preventDefault()
            this.$store.commit('closeHovers')
            action(overwrite, rename)
          }
        })

        return
      }

      action(overwrite, rename)
    },
    itemClick: function(event) {
      if (this.singleClick && !this.$store.state.multiple) this.open()
      else this.click(event)
    },
    click: function (event) {
      if (!this.singleClick && this.selectedCount !== 0) event.preventDefault()
      if (this.$store.state.selected.indexOf(this.index) !== -1) {
        this.removeSelected(this.index)
        return
      }

      if (event.shiftKey && this.selected.length > 0) {
        let fi = 0
        let la = 0

        if (this.index > this.selected[0]) {
          fi = this.selected[0] + 1
          la = this.index
        } else {
          fi = this.index
          la = this.selected[0] - 1
        }

        for (; fi <= la; fi++) {
          if (this.$store.state.selected.indexOf(fi) == -1) {
            this.addSelected(fi)
          }
        }

        return
      }

      if (!this.singleClick && !event.ctrlKey && !event.metaKey && !this.$store.state.multiple) this.resetSelected()
      this.addSelected(this.index)
    },
    dblclick: function () {
      if (!this.singleClick) this.open()
    },
    touchstart () {
      setTimeout(() => {
        this.touches = 0
      }, 300)

      this.touches++
      if (this.touches > 1) {
        this.open()
      }
    },
    open: function () {
      if (this.isTrash)return
      this.$router.push({path: this.url})
    },
    isFontIcon () {
      return iconTheme === "font-icon"
    },
    isColorIcon () {
      return iconTheme === "color-icon"
    },
    isMaterial () {
      return iconTheme === "material-icon"
    },
    imaterialIcons() {
      if (this.isDir) this.icon = 'folder'
      else if (this.type === 'image') this.icon = 'insert_photo'
      else if (this.type === 'audio') this.icon = 'volume_up'
      else if (this.type === 'video') this.icon = 'movie'
      else this.icon = 'insert_drive_file'
    },
    colorIcons () {
        if (this.isDir){
          this.icon = '#icon-folder'
          return
        }

        var items = this.name.split('.')
        var ext_name = items[items.length - 1].toLowerCase()
        var icons = {
          '7z': 'icon-7z',
          'avi': 'icon-avi',
          'apk': 'icon-apk',
          'bat': 'icon-bat',
          'c': 'icon-c',
          'cad': 'icon-cad',
          'cc': 'icon-cpp',
          'chm': 'icon-chm',
          'cpp': 'icon-cpp',
          'cxx': 'icon-cpp',
          'css': 'icon-css',
          'csv': 'icon-csv',
          'doc': 'icon-doc',
          'docx': 'icon-docx',
          'dll': 'icon-dll',
          'exe': 'icon-exe',
          'gif': 'icon-gif',
          'gz': 'icon-gz',
          'gzip': 'icon-gzip',
          'go': 'icon-go',
          'h': 'icon-h',
          'htm': 'icon-html',
          'html': 'icon-html',
          'hpp': 'icon-hpp',
          'ini': 'icon-ini',
          'iso': 'icon-iso',
          'jar': 'icon-jar',
          'jpg': 'icon-jpg',
          'jpeg': 'icon-jpg',
          'json': 'icon-json',
          'js': 'icon-js',
          'java': 'icon-java',
          'lua': 'icon-lua',
          'log': 'icon-log',
          'mp3': 'icon-mp3',
          'mp4': 'icon-mp4',
          'mkv': 'icon-mkv',
          'msi': 'icon-msi',
          'png': 'icon-png',
          'psd': 'icon-psd',
          'ppt': 'icon-ppt',
          'pptx': 'icon-pptx',
          'py': 'icon-py',
          'pyc': 'icon-pyc',
          'php': 'icon-php',
          'pdf': 'icon-pdf',
          'rar': 'icon-rar',
          'sh': 'icon-sh',
          'sql': 'icon-sql',
          'txt': 'icon-txt',
          'tar': 'icon-tar',
          'vue': 'icon-vue',
          'xls': 'icon-xls',
          'xlsx': 'icon-xlsx',
          'xz': 'icon-xz',
          'xml': 'icon-xml',
          'yml': 'icon-yml',
          'zip': 'icon-zip',
        }
        if (icons[ext_name] != undefined) {
          this.icon = icons[ext_name]
        }
        else if (this.type === 'image'){
          this.icon = 'icon-image'
        }
        else if (this.type === 'audio'){
          this.icon = 'icon-audio'
        }
        else if (this.type === 'video'){
          this.icon = 'icon-video'
        } else {
          this.icon = 'icon-file'
        }
        this.icon = '#' + this.icon
    },
    fontIcons () {
        if (this.isDir){
          this.icon = 'iconfont icon-folder'
          return
        }

        var items = this.name.split('.')
        var ext_name = items[items.length - 1].toLowerCase()
        var icons = {
          '7z': 'icon-7z',
          'avi': 'icon-avi',
          'apk': 'icon-apk',
          'bat': 'icon-bat',
          'c': 'icon-c',
          'cad': 'icon-cad',
          'cc': 'icon-cpp',
          'chm': 'icon-chm',
          'cpp': 'icon-cpp',
          'cxx': 'icon-cpp',
          'css': 'icon-css',
          'csv': 'icon-csv',
          'doc': 'icon-doc',
          'docx': 'icon-docx',
          'dll': 'icon-dll',
          'exe': 'icon-exe',
          'gif': 'icon-gif',
          'gz': 'icon-gz',
          'gzip': 'icon-gzip',
          'go': 'icon-go',
          'h': 'icon-h',
          'htm': 'icon-html',
          'html': 'icon-html',
          'hpp': 'icon-hpp',
          'ini': 'icon-ini',
          'iso': 'icon-iso',
          'jar': 'icon-jar',
          'jpg': 'icon-jpg',
          'jpeg': 'icon-jpg',
          'json': 'icon-json',
          'js': 'icon-js',
          'java': 'icon-java',
          'lua': 'icon-lua',
          'log': 'icon-log',
          'mp3': 'icon-mp3',
          'mp4': 'icon-mp4',
          'mkv': 'icon-mkv',
          'msi': 'icon-msi',
          'png': 'icon-png',
          'psd': 'icon-psd',
          'ppt': 'icon-ppt',
          'pptx': 'icon-pptx',
          'py': 'icon-py',
          'pyc': 'icon-pyc',
          'php': 'icon-php',
          'pdf': 'icon-pdf',
          'rar': 'icon-rar',
          'sh': 'icon-sh',
          'sql': 'icon-sql',
          'txt': 'icon-txt',
          'tar': 'icon-tar',
          'vue': 'icon-vue',
          'xls': 'icon-xls',
          'xlsx': 'icon-xlsx',
          'xz': 'icon-xz',
          'xml': 'icon-xml',
          'yml': 'icon-yml',
          'zip': 'icon-zip',
        }
        if (icons[ext_name] != undefined) {
          this.icon = icons[ext_name]
        }
        else if (this.type === 'image'){
          this.icon = 'icon-image'
        }
        else if (this.type === 'audio'){
          this.icon = 'icon-audio'
        }
        else if (this.type === 'video'){
          this.icon = 'icon-video'
        } else {
          this.icon = 'icon-file'
        }
        this.icon = 'iconfont ' + this.icon
    }
  },
}
</script>