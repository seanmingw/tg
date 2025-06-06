<template>
  <el-dialog
    :title="$t('table.nodeQRCode').toString()"
    :visible="dialogVisibleProps"
    @close="$emit('update:dialogVisibleProps', false)"
  >
    <el-image
      style="width: 256px; height: 256px"
      :src="qrCodeSrcProps"
    ></el-image>
    <div slot="footer" class="dialog-footer">
      <el-button
        type="primary"
        @click="$emit('update:dialogVisibleProps', false)"
      >
        {{ $t('table.confirm') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script>
export default {
  name: 'NodeQrcode',
  props: {
    dialogVisibleProps: {
      type: Boolean,
      required: true
    },
    qrCodeSrcProps: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped></style>
