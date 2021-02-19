<template>
  <el-select v-model="curr_theme" placeholder="请选择">
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
  name: 'themes',
  props: [ 'theme' ],
  data() {
    return {
      options: [{
        value: '',
        label: this.$t('settings.themes.light'),
      },{
        value: 'dark',
        label: this.$t('settings.themes.dark'),
      },{
        value: 'color',
        label: this.$t('settings.themes.color'),
      },
      ],
      curr_theme: this.theme
    }
  },
  watch: {
    theme(newVal) {
      this.curr_theme = newVal;
    },
    curr_theme(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.$emit('update:theme', this.curr_theme)
      }
    }
  }
}
</script>