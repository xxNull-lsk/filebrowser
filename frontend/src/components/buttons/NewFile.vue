<template>
  <button @click="$store.commit('showHover', 'newFile')" class="action" :aria-label="$t('sidebar.newFile')"
    :title="$t('sidebar.newFile')" id="new-file-button"
  >
    <i class="material-icons">note_add</i>
    <span>{{ $t("sidebar.newFile") }}</span>
  </button>
</template>

<script>
export default {
  name: "new-file-button",
};
</script>
