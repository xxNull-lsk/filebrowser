<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ $t('prompts.favorite') }}</h2>
    </div>

    <div class="card-content">
      <p>{{ $t('prompts.favoriteMessage') }}</p>
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
      name: ''
    };
  },
  computed: {
    ...mapState(['req', 'selected', 'selectedCount']),
    ...mapGetters(['isListing'])
  },
  methods: {
    submit: async function() {
      let path = ''

      if (this.selectedCount != undefined && this.selectedCount !== 0) {
        path = this.req.items[this.selected[0]].url
      } else {
        path = this.req.url
      }
      "".s
      console.error(path)

      try {
        await favorite_api.create(path, this.name)
        this.$router.push({ path: this.$route.path })
      } catch (e) {
        console.error(e)
        //this.$showError(e)
      }

      this.$store.commit('closeHovers')
    }
  }
};
</script>

