<template>
  <div v-if="formVisibleProps">
    <el-form-item>
      <aside>
        {{ $t('table.naiveproxyTip') }}
      </aside>
    </el-form-item>
  </div>
</template>

<script>
export default {
  name: 'NaiveProxyForm',
  props: {
    formVisibleProps: {
      type: Boolean,
      require: true
    }
  }
}
</script>

<style scoped></style>
