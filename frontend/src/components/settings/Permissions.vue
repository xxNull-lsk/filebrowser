<template>
  <div>
    <h3>{{ $t('settings.permissions') }}</h3>
    <p class="small">{{ $t('settings.permissionsHelp') }}</p>

    <p><el-checkbox v-model="admin"> {{ $t('settings.administrator') }}</p>

    <p><el-checkbox :disabled="admin" v-model="perm.create"> {{ $t('settings.perm.create') }}</el-checkbox></p>
    <p><el-checkbox :disabled="admin" v-model="perm.delete"> {{ $t('settings.perm.delete') }}</el-checkbox></p>
    <p><el-checkbox :disabled="admin" v-model="perm.download"> {{ $t('settings.perm.download') }}</el-checkbox></p>
    <p><el-checkbox :disabled="admin" v-model="perm.modify"> {{ $t('settings.perm.modify') }}</el-checkbox></p>
    <p v-if="isExecEnabled"><el-checkbox :disabled="admin" v-model="perm.execute"> {{ $t('settings.perm.execute') }}</el-checkbox></p>
    <p><el-checkbox :disabled="admin" v-model="perm.rename"> {{ $t('settings.perm.rename') }}</el-checkbox></p>
    <p><el-checkbox :disabled="admin" v-model="perm.share"> {{ $t('settings.perm.share') }}</el-checkbox></p>
  </div>
</template>

<script>
import { enableExec } from '@/utils/constants'
export default {
  name: 'permissions',
  props: ['perm'],
  computed: {
    admin: {
      get () {
        return this.perm.admin
      },
      set (value) {
        if (value) {
          for (const key in this.perm) {
            this.perm[key] = true
          }
        }

        this.perm.admin = value
      }
    },
    isExecEnabled: () => enableExec
  }
}
</script>
