<template>
  <button
    @click="showFavorite"
    :aria-label="$t('buttons.favorite')"
    :title="$t('buttons.favorite')"
    class="action"
    id="favorite-button"
  >
    <i class="material-icons">{{favoriteInfo == undefined?"favorite-border":"favorite"}}</i>
    <span>{{ $t("buttons.favorite") }}</span>
  </button>
</template>

<script>
import { favorites as favorite_api } from '@/api'

export default {
  name: "favorite-button",
  data: function () {
    return {
      favoriteInfo: undefined,
    };
  },
  async beforeMount() {
    try {
      this.favoriteInfo = await favorite_api.get(this.$route.path);
    } catch (e) {
      this.$showError(e)
    }
  },
  methods: {
    showFavorite: function () {
      this.$store.commit("showHover", "favorite");
    },
  },
};
</script>
