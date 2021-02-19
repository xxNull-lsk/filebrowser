<template>
<el-select v-model="curr_icon_theme" placeholder="请选择">
    <el-option
      v-for="item in options"
      :key="item.value"
      :label="item.label"
      :value="item.value">
    </el-option>
  </el-select>
</template>

<script>
export default {
  name: 'icon-themes',
  props: [ 'icon_theme' ],
  data() {
    return {
      options: [{
        value: 'material-icon',
        label: this.$t('settings.iconThemes.material'),
      },{/*
        value: 'font-icon',
        label: this.$t('settings.iconThemes.font'),
      },{*/
        value: 'color-icon',
        label: this.$t('settings.iconThemes.color'),
      },
      ],
      curr_icon_theme: this.icon_theme
    }
  },
  watch: {
    icon_theme(newVal) {
      this.curr_icon_theme = newVal;
    },
    curr_icon_theme(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.$emit('update:icon_theme', this.curr_icon_theme)
      }
    }
  }
}
</script>