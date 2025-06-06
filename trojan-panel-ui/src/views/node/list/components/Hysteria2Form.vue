<template>
  <div v-if="formVisibleProps">
    <el-form-item
      :label="$t('table.hysteria2ObfsPassword').toString()"
      prop="hysteria2ObfsPassword"
    >
      <el-input v-model="nodeProps.hysteria2ObfsPassword" />
    </el-form-item>
    <el-form-item
      :label="$t('table.hysteria2UpMbps').toString()"
      prop="hysteria2UpMbps"
    >
      <el-input-number
        v-model.number="nodeProps.hysteria2UpMbps"
        controls-position="right"
        type="number"
      />
    </el-form-item>
    <el-form-item
      :label="$t('table.hysteria2DownMbps').toString()"
      prop="hysteria2DownMbps"
    >
      <el-input-number
        v-model.number="nodeProps.hysteria2DownMbps"
        controls-position="right"
        type="number"
      />
    </el-form-item>
    <el-form-item
      :label="$t('table.hysteria2ServerName').toString()"
      prop="hysteria2ServerName"
    >
      <el-input v-model="nodeProps.hysteria2ServerName" />
    </el-form-item>
    <el-form-item
      :label="$t('table.hysteria2Insecure').toString()"
      prop="hysteria2Insecure"
    >
      <el-switch
        v-model="nodeProps.hysteria2Insecure"
        active-color="#13ce66"
        inactive-color="#ff4949"
        :active-text="$t('table.enable').toString()"
        :inactive-text="$t('table.disable').toString()"
        :active-value="1"
        :inactive-value="0"
      >
      </el-switch>
    </el-form-item>
  </div>
</template>

<script>
export default {
  name: 'Hysteria2Form',
  props: {
    nodeProps: {
      type: Object,
      require: true
    },
    formVisibleProps: {
      type: Boolean,
      require: true
    }
  },
  data() {
    return {}
  }
}
</script>

<style scoped></style>