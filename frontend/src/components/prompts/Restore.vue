<template>
  <div class="card floating">
    <div class="card-content">
      <p v-if="req.kind !== 'listing'">{{ $t('prompts.restoreMessageSingle') }}</p>
      <p v-else>{{ $t('prompts.restoreMessageMultiple', { count: selectedCount}) }}</p>
    </div>
    <div class="card-action">
      <button @click="$store.commit('closeHovers')"
        class="button button--flat button--grey"
        :aria-label="$t('buttons.cancel')"
        :title="$t('buttons.cancel')">{{ $t('buttons.cancel') }}</button>
      <button @click="submit"
        class="button button--flat button--red"
        :aria-label="$t('buttons.restore')"
        :title="$t('buttons.restore')">{{ $t('buttons.restore') }}</button>
    </div>
  </div>
</template>

<script>
import {mapGetters, mapMutations, mapState} from 'vuex'
import { files as api } from '@/api'
import buttons from '@/utils/buttons'

export default {
  name: 'delete',
  data: () => ({
    error: null
  }),
  computed: {
    ...mapGetters(['isListing', 'isTrash', 'selectedCount']),
    ...mapState(['req', 'selected'])
  },
  methods: {
    ...mapMutations(['closeHovers']),
    submit: async function () {
      buttons.loading('restore')

      try {
        this.closeHovers()

        if (this.selectedCount === 0) {
          return
        }

        let promises = []
        for (let index of this.selected) {
            promises.push(api.untrash(this.req.items[index].path))
        }

        await Promise.all(promises)
        buttons.success('restore')
        this.$store.commit('setReload', true)
      } catch (e) {
        buttons.done('restore')
        this.error = e
        this.$store.commit('setReload', true)
      }
    }
  }
}
</script>
