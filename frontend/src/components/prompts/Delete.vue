<template>
  <div class="card floating">
    <div class="card-content">
      <p v-if="req.kind !== 'listing'">{{ $t(this.isTrash?'prompts.deleteMessageSingle':'prompts.deleteToTrashMessageSingle') }}</p>
      <p v-else>{{ $t(this.isTrash?'prompts.deleteMessageMultiple':'prompts.deleteToTrashMessageMultiple', { count: selectedCount}) }}</p>
      <p v-if="!isTrash">
        <el-checkbox v-model="force_delete">
        <label class="checkbox-title">
          {{  $t('buttons.forceDelete')  }}
        </label>
        </el-checkbox>
      </p>
    </div>
    <div class="card-action">
      <button @click="$store.commit('closeHovers')"
        class="button button--flat button--grey"
        :aria-label="$t('buttons.cancel')"
        :title="$t('buttons.cancel')">{{ $t('buttons.cancel') }}</button>
      <button @click="submit"
        class="button button--flat button--red"
        :aria-label="$t('buttons.delete')"
        :title="$t('buttons.delete')">{{ $t('buttons.delete') }}</button>
    </div>
  </div>
</template>

<script>
import {mapGetters, mapMutations, mapState} from 'vuex'
import { files as api } from '@/api'
import buttons from '@/utils/buttons'

export default {
  name: 'delete',
  data: function () {
    return {
      force_delete: false
    }
  },
  computed: {
    ...mapGetters(['isListing', 'isTrash', 'selectedCount']),
    ...mapState(['req', 'selected'])
  },
  methods: {
    ...mapMutations(['closeHovers']),
    submit: async function () {
      buttons.loading('delete')

      try {
        if (!this.isListing) {
          if (this.isTrash || this.force_delete) {
            await api.remove(this.$route.path)
            buttons.success('delete')
            this.$root.$emit('preview-deleted')
          } else {
            await api.trash(this.$route.path)
            this.$root.$emit('preview-trash')
          }

          this.closeHovers()
          return
        }

        this.closeHovers()

        if (this.selectedCount === 0) {
          return
        }

        let promises = []
        for (let index of this.selected) {
          if (this.isTrash || this.force_delete){
            promises.push(api.remove(this.req.items[index].url))
          }
          else {
            promises.push(api.trash(this.req.items[index].url))
          }
        }

        await Promise.all(promises)
        buttons.success('delete')
        this.$store.commit('setReload', true)
      } catch (e) {
        buttons.done('delete')
        this.$showError(e)
        if (this.isListing) this.$store.commit('setReload', true)
      }
    }
  }
}
</script>
