<template>
  <button @click="show" :aria-label="$t('buttons.restore')" :title="$t('buttons.restore')" class="action" id="restore-button">
    <i class="material-icons">restore</i>
    <span>{{ $t('buttons.restore') }}</span>
  </button>
</template>

<script>
export default {
  name: 'restore-button',
  methods: {
    show: function () {
      this.$store.commit('showHover', 'restore')
    }
  }
}
</script>
