<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ $t('prompts.favorite') }}</h2>
    </div>

    <div class="card-content">
      <p>{{ $t('prompts.favoriteMessage', [this.path]) }}</p>
      <input class="input input--block" v-focus type="text" @keyup.enter="submit" v-model.trim="name">
    </div>

    <div class="card-action">
      <button
        class="button button--flat button--grey"
        @click="$store.commit('closeHovers')"
        :aria-label="$t('buttons.cancel')"
        :title="$t('buttons.cancel')"
      >{{ $t('buttons.cancel') }}</button>
      <button
        class="button button--flat"
        @click="submit"
        :aria-label="$t('buttons.create')"
        :title="$t('buttons.create')"
      >{{ $t('buttons.create') }}</button>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex'
import { favorites as favorite_api } from '@/api'

export default {
  name: 'favorite',
  data: function() {
    return {
      name: '',
      path: ''
    };
  },
  computed: {
    ...mapState(['req', 'selected', 'selectedCount']),
    ...mapGetters(['isListing'])
  },
  async beforeMount() {
      let path = ''

      if (this.selected != undefined && this.selected.length > 0) {
        path = this.req.items[this.selected[0]].url
      } else {
        path = this.req.url
      }
      path = path.replace(/\/$/g, '')
      this.path = path
      var items = path.split('/')
      if (items.length > 0) {
        this.name = items[items.length - 1]
      }
  },
  methods: {
    submit: async function() {
      let path = ''

      if (this.selected != undefined && this.selected.length > 0) {
        path = this.req.items[this.selected[0]].url
      } else {
        path = this.req.url
      }

      try {
        path = path.replace(/^.files/g, '')
        await favorite_api.create(path, this.name)
      } catch (e) {
        console.error(e)
        //this.$showError(e)
      }

      this.$store.commit('closeHovers')
    }
  }
};
</script>

