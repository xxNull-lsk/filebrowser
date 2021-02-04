<template>
  <button @click="$store.commit('showHover', 'newDir')" class="action" :aria-label="$t('sidebar.newFolder')"
    :title="$t('sidebar.newFolder')" id="new-dir-button"
  >
    <i class="material-icons">create_new_folder</i>
    <span>{{ $t("sidebar.newFolder") }}</span>
  </button>
</template>

<script>
export default {
  name: "new-dir-button",
};
</script>
