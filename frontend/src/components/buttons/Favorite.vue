<template>
  <button
    @click="showFavorite"
    :aria-label="$t('buttons.favorite')"
    :title="$t('buttons.favorite')"
    class="action"
    id="favorite-button"
  >
    <i class="material-icons">favorite</i>
    <span>{{ $t("buttons.favorite") }}</span>
  </button>
</template>

<script>
export default {
  name: "favorite-button",
  methods: {
    showFavorite: function () {
      this.$store.commit("showHover", "favorite");
    },
  },
};
</script>
