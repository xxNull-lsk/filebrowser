<template>
  <select v-on:change="change" :value="icon_theme">
    <option value="material-icon">{{ $t('settings.iconThemes.material') }}</option>
    <option value="font-icon">{{ $t('settings.iconThemes.font') }}</option>
    <option value="color-icon">{{ $t('settings.iconThemes.color') }}</option>
  </select>
</template>

<script>
export default {
  name: 'icon-themes',
  props: [ 'icon_theme' ],
  methods: {
    change (event) {
      this.$emit('update:icon_theme', event.target.value)
    }
  }
}
</script>