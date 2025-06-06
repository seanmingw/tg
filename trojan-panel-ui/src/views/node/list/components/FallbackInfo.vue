<template>
  <el-dialog
    :title="$t('table.detail').toString()"
    :visible="dialogVisibleProps"
    @close="$emit('update:dialogVisibleProps', false)"
    width="30%"
  >
    <el-form :model="fallbackPrpops" label-position="left">
      <el-form-item label="name">
        <el-tag>{{ fallbackPrpops.name }}</el-tag>
      </el-form-item>
      <el-form-item label="alpn">
        <el-tag>{{ fallbackPrpops.alpn }}</el-tag>
      </el-form-item>
      <el-form-item label="path">
        <el-tag>{{ fallbackPrpops.path }}</el-tag>
      </el-form-item>
      <el-form-item label="dest">
        <el-tag>{{ fallbackPrpops.dest }}</el-tag>
      </el-form-item>
      <el-form-item label="xver">
        <el-tag>{{ fallbackPrpops.xver }}</el-tag>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button
        type="primary"
        @click="$emit('update:dialogVisibleProps', false)"
        >{{ $t('table.confirm') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script>
export default {
  name: 'FallbackInfo',
  props: {
    dialogVisibleProps: {
      type: Boolean,
      required: true
    },
    fallbackPrpops: {
      type: Object,
      required: true
    }
  }
}
</script>

<style scoped></style>
