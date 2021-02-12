<template>
  <select v-on:change="change" :value="icon_theme">
    <option value="material-icon">{{ $t('settings.themes.materialIcon') }}</option>
    <option value="font-icon">{{ $t('settings.themes.fontIcon') }}</option>
    <option value="color-icon">{{ $t('settings.themes.colorIcon') }}</option>
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