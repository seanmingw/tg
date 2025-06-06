<template>
  <div v-if="formVisibleProps">
    <el-form-item
      :label="$t('table.trojanGoSni').toString()"
      prop="trojanGoSni"
    >
      <el-input v-model="nodeProps.trojanGoSni" clearable />
    </el-form-item>
    <el-form-item
      :label="$t('table.trojanGoMuxEnable').toString()"
      prop="trojanGoMuxEnable"
    >
      <el-switch
        v-model="nodeProps.trojanGoMuxEnable"
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
      :label="$t('table.trojanGoWebsocketEnable').toString()"
      prop="trojanGoWebsocketEnable"
    >
      <el-switch
        v-model="nodeProps.trojanGoWebsocketEnable"
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
      :label="$t('table.trojanGoWebsocketPath').toString()"
      prop="trojanGoWebsocketPath"
      v-if="isTrojanGoEnableWebsocket(nodeProps)"
    >
      <el-input v-model="nodeProps.trojanGoWebsocketPath" clearable />
    </el-form-item>
    <el-form-item
      :label="$t('table.trojanGoWebsocketHost').toString()"
      prop="trojanGoWebsocketHost"
      v-if="isTrojanGoEnableWebsocket(nodeProps)"
    >
      <el-input v-model="nodeProps.trojanGoWebsocketHost" clearable />
    </el-form-item>
    <el-form-item
      :label="$t('table.trojanGoSsEnable').toString()"
      prop="trojanGoSsEnable"
      v-if="isTrojanGoEnableWebsocket(nodeProps)"
    >
      <el-switch
        v-model="nodeProps.trojanGoSsEnable"
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
      :label="$t('table.trojanGoSsMethod').toString()"
      prop="trojanGoSsMethod"
      v-if="
        isTrojanGoEnableWebsocket(nodeProps) && isTrojanGoEnableSs(nodeProps)
      "
    >
      <el-select
        v-model="nodeProps.trojanGoSsMethod"
        :placeholder="$t('table.trojanGoSsMethod').toString()"
        controls-position="right"
      >
        <el-option
          :label="item"
          :value="item"
          :key="item"
          v-for="item in trojanGoSsMethods"
        ></el-option>
      </el-select>
    </el-form-item>
    <el-form-item
      :label="$t('table.trojanGoSsPassword').toString()"
      prop="trojanGoSsPassword"
      v-if="
        isTrojanGoEnableWebsocket(nodeProps) && isTrojanGoEnableSs(nodeProps)
      "
    >
      <el-input v-model="nodeProps.trojanGoSsPassword" clearable />
    </el-form-item>
  </div>
</template>

<script>
import { isTrojanGoEnableSs, isTrojanGoEnableWebsocket } from '@/utils/node'

export default {
  name: 'TrojanGoForm',
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
      trojanGoSsMethods: [
        'AES-128-GCM',
        'AES-256-GCM',
        'CHACHA20-IETF-POLY1305'
      ]
    }
  },
  methods: {
    isTrojanGoEnableWebsocket,
    isTrojanGoEnableSs
  }
}
</script>

<style scoped></style>
