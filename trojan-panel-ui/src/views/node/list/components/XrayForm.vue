<template>
  <div v-if="formVisibleProps">
    <el-form-item
        :label="$t('table.xrayProtocol').toString()"
        prop="xrayProtocol"
    >
      <el-select
          v-model="nodeProps.xrayProtocol"
          controls-position="right"
          @change="xrayProtocolChange"
      >
        <el-option
            :label="item"
            :value="item"
            :key="index"
            v-for="(item, index) in xrayProtocols"
        ></el-option>
      </el-select>
    </el-form-item>
    <el-form-item
        :label="$t('table.xrayStreamSettingsNetwork').toString()"
        prop="xrayStreamSettingsEntity.network"
        v-if="
        !(isXrayShadowsocksAEAD(nodeProps) || isXrayShadowsocks2022(nodeProps))
      "
    >
      <el-select
          v-model="nodeProps.xrayStreamSettingsEntity.network"
          controls-position="right"
          @change="xrayStreamSettingsNetworkChange"
      >
        <el-option
            :label="item"
            :value="item"
            :key="item"
            v-for="item in xrayStreamSettingsNetworks"
        ></el-option>
      </el-select>
    </el-form-item>

    <XrayFormWebSocket
        :form-visible-props="isXrayWs(nodeProps) && !isXrayShadowsocks(nodeProps)"
        :node-props="nodeProps"
    />

    <el-form-item
        :label="$t('table.xrayStreamSettingsSecurity').toString()"
        prop="xrayStreamSettingsEntity.security"
        v-if="!isXrayShadowsocks(nodeProps)"
    >
      <el-select
          v-model="nodeProps.xrayStreamSettingsEntity.security"
          controls-position="right"
          @change="xrayStreamSettingsSecurityChange"
      >
        <el-option
            :label="item"
            :value="item"
            :key="item"
            v-for="item in xrayStreamSettingsSecuritys"
        ></el-option>
      </el-select>
    </el-form-item>

    <XrayFormTls
        :form-visible-props="isXrayStreamSettingsSecurityTls(nodeProps)"
        :node-props="nodeProps"
    />

    <XrayFormReality
        :form-visible-props="isXrayStreamSettingsSecurityReality(nodeProps)"
        :node-props="nodeProps"
    />

    <el-form-item
        :label="$t('table.xrayFlow').toString()"
        prop="xrayFlow"
        v-if="showXrayFlow(nodeProps)"
    >
      <el-select v-model="nodeProps.xrayFlow" controls-position="right">
        <el-option
            :label="item"
            :value="item"
            :key="index"
            v-for="(item, index) in xrayFlows"
        ></el-option>
      </el-select>
    </el-form-item>
    <el-form-item
        :label="$t('table.xraySSMethod').toString()"
        prop="xraySSMethod"
        v-if="
        isXrayShadowsocks(nodeProps) ||
        isXrayShadowsocksAEAD(nodeProps) ||
        isXrayShadowsocks2022(nodeProps)
      "
    >
      <el-select v-model="nodeProps.xraySSMethod" controls-position="right">
        <el-option
            :label="item"
            :value="item"
            :key="index"
            v-for="(item, index) in xraySSMethods"
        ></el-option>
      </el-select>
    </el-form-item>
    <el-form-item
        :label="$t('table.xraySSNetwork').toString()"
        prop="xraySettingsEntity.network"
        v-if="
        isXrayShadowsocksAEAD(nodeProps) ||
        isXrayShadowsocks2022(nodeProps)
      "
    >
      <el-select
          v-model="nodeProps.xraySettingsEntity.network"
          controls-position="right"
      >
        <el-option
            :label="item"
            :value="item"
            :key="index"
            v-for="(item, index) in xraySettingsNetworks"
        ></el-option>
      </el-select>
    </el-form-item>
    <el-form-item
        :label="$t('table.xrayFallbacks').toString()"
        prop="xraySettingsEntity.fallbacks"
        v-if="showFallback(nodeProps)"
    >
      <el-tag
          v-for="(item, index) in nodeProps.xraySettingsEntity.fallbacks"
          :key="index"
          :disable-transitions="true"
          type="success"
          @close="deleteFallbackProps(item)"
          effect="dark"
          size="medium"
          closable
          @click="handleFallbackDetailProps(item)"
      >
        {{ item.dest }}
      </el-tag>
      <el-button
          type="primary"
          size="mini"
          icon="el-icon-plus"
          @click="handleCreateFallbackProps"
      ></el-button>
    </el-form-item>
    <el-form-item :label="$t('table.xraySocksUser').toString()" prop="xraySettingsEntity.accounts[0].user" v-if="isXraySocks(nodeProps)">
      <el-input v-model="nodeProps.xraySettingsEntity.accounts[0].user" />
    </el-form-item>
    <el-form-item :label="$t('table.xraySocksPass').toString()" prop="xraySettingsEntity.accounts[0].pass" v-if="isXraySocks(nodeProps)">
      <el-input v-model="nodeProps.xraySettingsEntity.accounts[0].pass" />
    </el-form-item>
    <el-form-item
        :label="$t('table.xraySocksUdp').toString()"
        prop="xraySettingsEntity.udp"
        v-if="isXraySocks(nodeProps)"
    >
      <el-switch
          v-model="nodeProps.xraySettingsEntity.udp"
          active-color="#13ce66"
          inactive-color="#ff4949"
          :active-text="$t('table.enable').toString()"
          :inactive-text="$t('table.disable').toString()"
          :active-value="true"
          :inactive-value="false"
      >
      </el-switch>
    </el-form-item>
  </div>
</template>

<script>
import FallbackInfo from '@/views/node/list/components/FallbackInfo'
import FallbackForm from '@/views/node/list/components/FallbackForm'

import {
  isXrayShadowsocks,
  isXrayShadowsocks2022,
  isXrayShadowsocksAEAD,
  isXraySocks,
  isXrayStreamSettingsSecurityReality,
  isXrayStreamSettingsSecurityTls, isXrayTrojan,
  isXrayVless,
  isXrayWs,
  showFallback,
  showXrayFlow
} from '@/utils/node.js'
import XrayFormReality from '@/views/node/list/components/XrayFormReality'
import XrayFormTls from '@/views/node/list/components/XrayFormTls'
import XrayFormWebSocket from '@/views/node/list/components/XrayFormWebSocket'

export default {
  name: 'XrayForm',
  components: {
    XrayFormWebSocket,
    XrayFormTls,
    XrayFormReality,
    FallbackInfo,
    FallbackForm
  },
  props: {
    nodeProps: {
      type: Object,
      require: true
    },
    formVisibleProps: {
      type: Boolean,
      require: true
    },
    handleCreateFallbackProps: {
      type: Function,
      required: true
    },
    handleFallbackDetailProps: {
      type: Function,
      required: true
    },
    deleteFallbackProps: {
      type: Function,
      required: true
    }
  },
  computed: {
    xrayStreamSettingsSecuritys() {
      let securitys = ['none', 'tls']
      if (isXrayVless(this.nodeProps)) {
        securitys.push('reality')
      }
      if (isXrayTrojan(this.nodeProps)) {
        for (let i = 0; i < securitys.length; i++) {
          if (securitys[i] === 'none') {
            securitys.splice(i, 1)
          }
        }
      }
      return securitys
    },
    xrayFlows() {
      return ['none', 'xtls-rprx-vision']
    },
    xrayStreamSettingsNetworks() {
      let xrayStreamSettingsNetworks = [
        'tcp',
        // 'kcp',
        'ws'
        // 'http',
        // 'domainsocket',
        // 'quic',
        // 'grpc'
      ]
      if (this.isXrayShadowsocks(this.nodeProps) ||
          this.isXraySocks(this.nodeProps)) {
        this.nodeProps.xrayStreamSettingsEntity.network = 'tcp'
        return ['tcp']
      } else {
        return xrayStreamSettingsNetworks
      }
    }
  },
  data() {
    return {
      xraySettingsNetworks: ['tcp', 'udp', 'tcp,udp'],
      xraySSMethods: [
        'none',
        'aes-128-gcm',
        'aes-256-gcm',
        'chacha20-ietf-poly1305',
        'xchacha20-ietf-poly1305'
      ],
      xrayProtocols: [
        // 'dokodemo-door',
        // 'http',
        'vless',
        'vmess',
        'trojan',
        'shadowsocks',
        'socks'
      ]
    }
  },
  methods: {
    isXraySocks,
    isXrayShadowsocks,
    isXrayShadowsocksAEAD,
    isXrayShadowsocks2022,
    showFallback,
    showXrayFlow,
    isXrayWs,
    isXrayStreamSettingsSecurityTls,
    isXrayStreamSettingsSecurityReality,
    xrayProtocolChange() {
      if (isXrayVless(this.nodeProps)) {
        this.nodeProps.xrayStreamSettingsEntity.security = 'reality'
      } else if (isXrayShadowsocks(this.nodeProps) || isXraySocks(this.nodeProps)) {
        this.nodeProps.xrayStreamSettingsEntity.security = 'none'
      } else {
        this.nodeProps.xrayStreamSettingsEntity.security = 'tls'
      }
    },
    xrayStreamSettingsNetworkChange() {
      if (this.nodeProps.xrayStreamSettingsEntity.network === 'ws') {
        this.nodeProps.xrayStreamSettingsEntity.security = 'tls'
      }
    },
    xrayStreamSettingsSecurityChange() {
      if (
          isXrayVless(this.nodeProps) &&
          (this.nodeProps.xrayStreamSettingsEntity.security === 'tls' ||
              this.nodeProps.xrayStreamSettingsEntity.security === 'reality')
      ) {
        this.nodeProps.xrayFlow = 'xtls-rprx-vision'
      } else {
        this.nodeProps.xrayFlow = 'none'
      }
    }
  }
}
</script>

<style scoped>
.el-tag + .el-tag {
  margin-left: 10px;
}

.el-button {
  margin-left: 10px;
}
</style>
