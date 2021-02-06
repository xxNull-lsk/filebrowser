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
      <i v-else-if="use_material_icons" class="material-icons">{{ icon }}</i>
      <i v-else :class="icon"></i>
    </div>

    <div>
      <p class="name">{{ name }}</p>

      <p v-if="isDir" class="size" data-order="-1"></p>
      <p v-else class="size" :data-order="humanSize()">{{ humanSize() }}</p>

      <p class="modified">
        <time :datetime="modified">{{ humanTime() }}</time>
      </p>
    </div>
  </div>
</template>

<script>
import { baseURL, enableThumbs } from '@/utils/constants'
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
      use_material_icons: true,
      icon: ''
    }
  },
  props: ['name', 'isDir', 'url', 'type', 'size', 'modified', 'index'],
  mounted () {
      if (this.isDir){
        this.use_material_icons = false
        this.icon = 'fa fa-folder-o'
        return
      }

      var items = this.name.split('.')
      var ext_name = items[items.length - 1].toLowerCase()
      if (ext_name === 'apk'){
        this.use_material_icons = false
        this.icon = 'fa fa-android'
      } 
      else if (ext_name === 'txt' || ext_name === 'log'
              || ext_name === "json" || ext_name === "yml"){
        this.use_material_icons = false
        this.icon = 'fa fa-file-text-o'
      } 
      else if (ext_name === 'doc' || ext_name === 'docx'){
        this.use_material_icons = false
        this.icon = 'fa fa-file-word-o'
      } 
      else if (ext_name === 'xls' || ext_name === 'xlsx'){
        this.use_material_icons = false
        this.icon = 'fa fa-file-excel-o'
      } 
      else if (ext_name === 'ppt' || ext_name === 'pptx'){
        this.use_material_icons = false
        this.icon = 'fa fa-file-powerpoint-o'
      } 
      else if (ext_name === 'pdf'){
        this.use_material_icons = false
        this.icon = 'fa fa-file-pdf-o'
      } 
      else if (ext_name === 'zip' || ext_name == "rar" || ext_name == "7z" || ext_name == "gz" || ext_name == "tar"){
        this.use_material_icons = false
        this.icon = 'fa fa-file-zip-o'
      }
      else if (ext_name === 'c' || ext_name == "cpp"
              || ext_name == "h" || ext_name == "hpp"
              || ext_name == "py" || ext_name == "php"
              || ext_name == "js" || ext_name == "java"
              || ext_name == "html" || ext_name == "xml"
              || ext_name == "sh" || ext_name == "css"
              || ext_name == "vue" || ext_name == "md"){
        this.use_material_icons = false
        this.icon = 'fa fa-file-code-o'
      }
      else if (this.type === 'image'){
        this.use_material_icons = false
        this.icon = 'fa fa-file-photo-o'
      }
      else if (this.type === 'audio'){
        this.use_material_icons = false
        this.icon = 'fa fa-file-sound-o'
      }
      else if (this.type === 'video'){
        this.use_material_icons = false
        this.icon = 'fa fa-file-video-o'
      } else {
      this.use_material_icons = false
      this.icon = 'fa fa-file-o'
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
    }
  },
  methods: {
    ...mapMutations(['addSelected', 'removeSelected', 'resetSelected']),
    humanSize: function () {
      return filesize(this.size)
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
    }
  }
}
</script>