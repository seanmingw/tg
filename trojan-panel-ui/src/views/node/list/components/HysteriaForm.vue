<template>
  <div v-if="formVisibleProps">
    <el-form-item
      :label="$t('table.hysteriaProtocol').toString()"
      prop="hysteriaProtocol"
    >
      <el-select
        v-model="nodeProps.hysteriaProtocol"
        :placeholder="$t('table.hysteriaProtocol').toString()"
        controls-position="right"
      >
        <el-option
          :label="item"
          :value="item"
          :key="item"
          v-for="item in hysteriaProtocols"
        ></el-option>
      </el-select>
    </el-form-item>
    <el-form-item
      :label="$t('table.hysteriaObfs').toString()"
      prop="hysteriaObfs"
    >
      <el-input v-model="nodeProps.hysteriaObfs" />
    </el-form-item>
    <el-form-item
      :label="$t('table.hysteriaUpMbps').toString()"
      prop="hysteriaUpMbps"
    >
      <el-input-number
        v-model.number="nodeProps.hysteriaUpMbps"
        controls-position="right"
        type="number"
      />
    </el-form-item>
    <el-form-item
      :label="$t('table.hysteriaDownMbps').toString()"
      prop="hysteriaDownMbps"
    >
      <el-input-number
        v-model.number="nodeProps.hysteriaDownMbps"
        controls-position="right"
        type="number"
      />
    </el-form-item>
    <el-form-item
      :label="$t('table.hysteriaServerName').toString()"
      prop="hysteriaServerName"
    >
      <el-input v-model="nodeProps.hysteriaServerName" />
    </el-form-item>
    <el-form-item
      :label="$t('table.hysteriaInsecure').toString()"
      prop="hysteriaInsecure"
    >
      <el-switch
        v-model="nodeProps.hysteriaInsecure"
        active-color="#13ce66"
        inactive-color="#ff4949"
        :active-text="$t('table.enable').toString()"
        :inactive-text="$t('table.disable').toString()"
        :active-value="1"
        :inactive-value="0"
      >
      </el-switch>
    </el-form-item>
    <el-form-item
      :label="$t('table.hysteriaFastOpen').toString()"
      prop="hysteriaFastOpen"
    >
      <el-switch
        v-model="nodeProps.hysteriaFastOpen"
        active-color="#13ce66"
        inactive-color="#ff4949"
        :active-text="$t('table.enable').toString()"
        :inactive-text="$t('table.disable').toString()"
        :active-value="1"
        :inactive-value="0"
      >
      </el-switch>
    </el-form-item>
    <el-form-item>
      <aside>
        {{ $t('table.hysteriaTip') }}
      </aside>
    </el-form-item>
  </div>
</template>

<script>
export default {
  name: 'HysteriaForm',
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
    return {
      hysteriaProtocols: ['udp', 'wechat-video', 'faketcp']
    }
  }
}
</script>

<style scoped></style>
