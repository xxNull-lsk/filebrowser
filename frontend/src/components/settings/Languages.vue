<template>
  <el-select v-model="curr_locale" placeholder="请选择">
    <el-option
      v-for="(language, value) in locales"
      :key="value"
      :label="$t('languages.' + language)"
      :value="value">
    </el-option>
  </el-select>
</template>

<script>

export default {
  name: 'languages',
  props: [ 'locale' ],
  data() {
    let dataObj = {
      locales: {
        ar: 'ar',
        de: 'de',
        en: 'en',
        es: 'es',
        fr: 'fr',
        is: 'is',
        it: 'it',
        ja: 'ja',
        ko: 'ko',
        'nl-be': 'nlBE',
        pl: 'pl',
        'pt-br': 'ptBR',
        pt: 'pt',
        ro: 'ro',
        ru: 'ru',
        'sv-se': 'svSE',
        'zh-cn': 'zhCN',
        'zh-tw': 'zhTW'
      },
      curr_locale: this.locale
    };

    Object.defineProperty(dataObj, "locales", { configurable: false, writable: false });

    return dataObj;
  },
  watch: {
    locale(newVal) {
      this.curr_locale = newVal;
    },
    curr_locale(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.$emit('update:locale', this.curr_locale)
      }
    }
  }
}
</script>
