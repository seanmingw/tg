<template>
  <div>
    <el-dialog
      :title="textMap[dialogStatusProps]"
      :visible="dialogFormVisibleProps"
      @close="$emit('update:dialogFormVisibleProps', false)"
    >
      <el-form
        ref="dataForm"
        :rules="dialogStatusProps === 'create' ? createRules : updateRules"
        :model="nodeProps"
        label-position="left"
      >
        <el-form-item :label="$t('table.nodeName').toString()" prop="name">
          <el-input v-model="nodeProps.name" clearable />
        </el-form-item>
        <el-form-item
          :label="$t('table.nodeServer').toString()"
          prop="nodeServerId"
        >
          <el-select v-model="nodeProps.nodeServerId" controls-position="right">
            <el-option
              :label="item.name"
              :value="item.id"
              :key="item.id"
              v-for="item in nodeServersProps"
            ></el-option>
          </el-select>
          <el-button
            type="primary"
            size="mini"
            icon="el-icon-plus"
            @click="toAddNodeServer"
          ></el-button>
        </el-form-item>
        <el-form-item :label="$t('table.nodeDomain').toString()" prop="domain">
          <el-input v-model="nodeProps.domain" clearable />
        </el-form-item>
        <el-form-item :label="$t('table.nodePort').toString()" prop="port">
          <el-input-number
            v-model.number="nodeProps.port"
            controls-position="right"
            type="number"
          />
        </el-form-item>

        <el-form-item
          :label="$t('table.nodePriority').toString()"
          prop="priority"
        >
          <el-input-number
            v-model.number="nodeProps.priority"
            controls-position="right"
            type="number"
          />
        </el-form-item>

        <el-form-item
          :label="$t('table.nodeType').toString()"
          prop="nodeTypeId"
        >
          <el-select v-model="nodeProps.nodeTypeId" controls-position="right">
            <el-option
              :label="item.name"
              :value="item.id"
              :key="item.id"
              v-for="item in nodeTypesProps"
            ></el-option>
          </el-select>
        </el-form-item>

        <XrayForm
          :node-props="nodeProps"
          :form-visible-props="isXray(nodeProps)"
          :handle-create-fallback-props="handleCreateFallback"
          :handle-fallback-detail-props="handleFallbackDetail"
          :delete-fallback-props="deleteFallback"
        />

        <TrojanGoForm
          :node-props="nodeProps"
          :form-visible-props="isTrojanGo(nodeProps)"
        />

        <HysteriaForm
          :node-props="nodeProps"
          :form-visible-props="isHysteria(nodeProps)"
        />

        <NaiveProxyForm :form-visible-props="isNaiveProxy(nodeProps)" />

        <Hysteria2Form
          :node-props="nodeProps"
          :form-visible-props="isHysteria2(nodeProps)"
        />
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="$emit('update:dialogFormVisibleProps', false)"
          >{{ $t('table.cancel') }}
        </el-button>
        <el-button
          type="primary"
          @click="dialogStatusProps === 'create' ? createData() : updateData()"
        >
          {{ $t('table.confirm') }}
        </el-button>
      </div>
    </el-dialog>

    <FallbackInfo
      :dialog-visible-props.sync="dialogFallbackDetailVisible"
      :fallback-prpops="fallback"
    />

    <FallbackForm
      ref="fallbackForm"
      :create-fallback-props="createFallback"
      :dialog-visible-props.sync="dialogFallbackFormVisible"
    />
  </div>
</template>

<script>
import FallbackInfo from '@/views/node/list/components/FallbackInfo'
import FallbackForm from '@/views/node/list/components/FallbackForm'
import { createNode, updateNodeById } from '@/api/node'
import XrayForm from '@/views/node/list/components/XrayForm'
import TrojanGoForm from '@/views/node/list/components/TrojanGoForm'
import {
  handleXraySettings,
  handleXrayStreamSettings,
  isHysteria,
  isHysteria2,
  isNaiveProxy,
  isTrojanGo,
  isXray,
  isXrayShadowsocks2022,
  isXrayShadowsocksAEAD,
  isXrayStreamSettingsSecurityReality
} from '@/utils/node'
import HysteriaForm from '@/views/node/list/components/HysteriaForm'
import NaiveProxyForm from '@/views/node/list/components/NaiveProxyForm'
import Hysteria2Form from '@/views/node/list/components/Hysteria2Form'

export default {
  name: 'NodeForm',
  components: {
    Hysteria2Form,
    NaiveProxyForm,
    HysteriaForm,
    TrojanGoForm,
    XrayForm,
    FallbackInfo,
    FallbackForm
  },
  props: {
    dialogFormVisibleProps: {
      type: Boolean,
      required: true
    },
    nodeProps: {
      type: Object,
      required: true
    },
    dialogStatusProps: {
      type: String,
      required: true
    },
    nodeServersProps: {
      type: Array,
      required: true
    },
    nodeTypesProps: {
      type: Array,
      required: true
    },
    getListProps: {
      type: Function,
      required: true
    }
  },
  data() {
    const validateXraySettingsEntityNetwork = (rule, value, callback) => {
      if (
        (isXrayShadowsocksAEAD(this.nodeProps) ||
          isXrayShadowsocks2022(this.nodeProps)) &&
        this.nodeProps.xraySettingsEntity.network.trim().length === 0
      ) {
        callback(new Error(this.$t('valid.xraySSNetwork').toString()))
      } else {
        callback()
      }
    }
    const validateXrayStreamSettingsEntityNetwork = (rule, value, callback) => {
      if (
        !isXrayShadowsocksAEAD(this.nodeProps) &&
        !isXrayShadowsocks2022(this.nodeProps) &&
        this.nodeProps.xrayStreamSettingsEntity.network.trim().length === 0
      ) {
        callback(new Error(this.$t('valid.xrayNetwork').toString()))
      } else {
        callback()
      }
    }
    const validateXrayStreamSettingsEntityRealitySettingsDest = (
      rule,
      value,
      callback
    ) => {
      if (
        isXrayStreamSettingsSecurityReality(this.nodeProps) &&
        !this.nodeProps.xrayStreamSettingsEntity.realitySettings.dest
      ) {
        callback(
          new Error(
            this.$t(
              'valid.xrayStreamSettingsEntityRealitySettingsDest'
            ).toString()
          )
        )
      } else {
        callback()
      }
    }
    const validateXrayStreamSettingsEntityRealitySettingsXver = (
      rule,
      value,
      callback
    ) => {
      console.log(this.nodeProps)
      if (
        isXrayStreamSettingsSecurityReality(this.nodeProps) &&
        this.nodeProps.xrayStreamSettingsEntity.realitySettings.xver !== 0 &&
        this.nodeProps.xrayStreamSettingsEntity.realitySettings.xver !== 1 &&
        this.nodeProps.xrayStreamSettingsEntity.realitySettings.xver !== 2
      ) {
        callback(
          new Error(
            this.$t(
              'valid.xrayStreamSettingsEntityRealitySettingsXver'
            ).toString()
          )
        )
      } else {
        callback()
      }
    }
    const validateXrayStreamSettingsEntityRealitySettingsServerNames = (
      rule,
      value,
      callback
    ) => {
      if (
        isXrayStreamSettingsSecurityReality(this.nodeProps) &&
        (!this.nodeProps.xrayStreamSettingsEntity.realitySettings.serverNames ||
          this.nodeProps.xrayStreamSettingsEntity.realitySettings.serverNames
            .length === 0)
      ) {
        callback(
          new Error(
            this.$t(
              'valid.xrayStreamSettingsEntityRealitySettingsServerNames'
            ).toString()
          )
        )
      } else {
        callback()
      }
    }
    const validateXrayStreamSettingsEntityRealitySettingsPrivateKey = (
      rule,
      value,
      callback
    ) => {
      if (
        isXrayStreamSettingsSecurityReality(this.nodeProps) &&
        !this.nodeProps.xrayStreamSettingsEntity.realitySettings.privateKey
      ) {
        callback(
          new Error(
            this.$t(
              'valid.xrayStreamSettingsEntityRealitySettingsPrivateKey'
            ).toString()
          )
        )
      } else {
        callback()
      }
    }
    const validateXrayStreamSettingsEntityRealitySettingsShortIds = (
      rule,
      value,
      callback
    ) => {
      if (
        isXrayStreamSettingsSecurityReality(this.nodeProps) &&
        (!this.nodeProps.xrayStreamSettingsEntity.realitySettings.shortIds ||
          this.nodeProps.xrayStreamSettingsEntity.realitySettings.shortIds
            .length === 0)
      ) {
        callback(
          new Error(
            this.$t(
              'valid.xrayStreamSettingsEntityRealitySettingsShortIds'
            ).toString()
          )
        )
      } else {
        callback()
      }
    }
    return {
      fallback: {
        name: '',
        alpn: '',
        path: undefined,
        dest: '80',
        xver: 0
      },
      dialogFallbackFormVisible: false,
      dialogFallbackDetailVisible: false,
      textMap: {
        update: this.$t('table.edit'),
        create: this.$t('table.add')
      },
      createRules: {
        name: [
          {
            required: true,
            message: this.$t('valid.nodeName'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 20,
            message: this.$t('valid.nodeNameRange'),
            trigger: ['change', 'blur']
          }
        ],
        nodeServerId: [
          {
            required: true,
            message: this.$t('valid.nodeServerId'),
            trigger: ['change', 'blur']
          }
        ],
        domain: [
          {
            required: true,
            message: this.$t('valid.nodeDomain'),
            trigger: ['change', 'blur']
          },
          {
            min: 4,
            max: 64,
            message: this.$t('valid.nodeDomainRange'),
            trigger: ['change', 'blur']
          }
        ],
        port: [
          {
            required: true,
            message: this.$t('valid.nodePort'),
            trigger: ['change', 'blur']
          },
          {
            pattern:
              /^([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/,
            message: this.$t('valid.nodePortRange'),
            trigger: ['change', 'blur']
          }
        ],
        priority: [
          {
            required: true,
            message: this.$t('valid.nodePriority'),
            trigger: ['change', 'blur']
          },
          {
            pattern: /^-?\d+$/,
            message: this.$t('valid.nodePriorityRange'),
            trigger: ['change', 'blur']
          }
        ],
        nodeTypeId: [
          {
            required: true,
            message: this.$t('valid.nodeType'),
            trigger: ['change', 'blur']
          }
        ],
        xrayProtocol: [
          {
            required: true,
            message: this.$t('valid.xrayProtocol'),
            trigger: ['change', 'blur']
          }
        ],
        xraySSMethod: [
          {
            required: true,
            message: this.$t('valid.xraySSMethod'),
            trigger: ['change', 'blur']
          }
        ],
        realityPbk: [
          {
            required: true,
            message: this.$t('valid.realityPbk'),
            trigger: ['change', 'blur']
          }
        ],
        'xraySettingsEntity.network': [
          {
            validator: validateXraySettingsEntityNetwork,
            trigger: ['change', 'blur']
          }
        ],
        'xraySettingsEntity.udp': [
          {
            required: true,
            message: this.$t('valid.xraySocksUdp'),
            trigger: ['change', 'blur']
          }
        ],
        'xraySettingsEntity.accounts[0].user': [
          {
            required: true,
            message: this.$t('valid.xraySocksUser'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 64,
            message: this.$t('valid.xraySocksUserRange'),
            trigger: ['change', 'blur']
          }
        ],
        'xraySettingsEntity.accounts[0].pass': [
          {
            required: true,
            message: this.$t('valid.xraySocksPass'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 64,
            message: this.$t('valid.xraySocksPassRange'),
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.network': [
          {
            validator: validateXrayStreamSettingsEntityNetwork,
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.security': [
          {
            required: true,
            message: this.$t('valid.xraySecurity'),
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.wsSettings.path': [
          {
            min: 2,
            max: 64,
            message: this.$t('valid.xrayWsRange'),
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.wsSettings.headers.Host': [
          {
            min: 2,
            max: 64,
            message: this.$t('valid.xrayWsRange'),
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.realitySettings.dest': [
          {
            required: true,
            message: this.$t(
              'valid.xrayStreamSettingsEntityRealitySettingsDest'
            ),
            trigger: ['change', 'blur']
          },
          {
            validator: validateXrayStreamSettingsEntityRealitySettingsDest,
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.realitySettings.xver': [
          {
            validator: validateXrayStreamSettingsEntityRealitySettingsXver,
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.realitySettings.serverNames': [
          {
            required: true,
            message: this.$t(
              'valid.xrayStreamSettingsEntityRealitySettingsServerNames'
            ),
            trigger: ['change', 'blur']
          },
          {
            validator:
              validateXrayStreamSettingsEntityRealitySettingsServerNames,
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.realitySettings.privateKey': [
          {
            required: true,
            message: this.$t(
              'valid.xrayStreamSettingsEntityRealitySettingsPrivateKey'
            ),
            trigger: ['change', 'blur']
          },
          {
            validator:
              validateXrayStreamSettingsEntityRealitySettingsPrivateKey,
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.realitySettings.shortIds': [
          {
            required: true,
            message: this.$t(
              'valid.xrayStreamSettingsEntityRealitySettingsShortIds'
            ),
            trigger: ['change', 'blur']
          },
          {
            validator: validateXrayStreamSettingsEntityRealitySettingsShortIds,
            trigger: ['change', 'blur']
          }
        ],
        trojanGoSni: [
          {
            min: 4,
            max: 64,
            message: this.$t('valid.trojanGoSni'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoMuxEnable: [
          {
            required: true,
            message: this.$t('valid.trojanGoMux'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoWebsocketEnable: [
          {
            required: true,
            message: this.$t('valid.trojanGoWs'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoWebsocketPath: [
          {
            min: 2,
            max: 64,
            message: this.$t('valid.trojanGoWsRange'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoWebsocketHost: [
          {
            min: 2,
            max: 64,
            message: this.$t('valid.trojanGoWsHostRange'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoSsEnable: [
          {
            required: true,
            message: this.$t('valid.trojanGoSs'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoSsMethod: [
          {
            required: true,
            message: this.$t('valid.trojanGoSsMethod'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoSsPassword: [
          {
            min: 2,
            max: 32,
            message: this.$t('valid.trojanGoSsPasswordRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaProtocol: [
          {
            required: true,
            message: this.$t('valid.hysteriaProtocol'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 16,
            message: this.$t('valid.hysteriaProtocolRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaObfs: [
          {
            min: 0,
            max: 64,
            message: this.$t('valid.hysteriaObfsRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaUpMbps: [
          {
            required: true,
            message: this.$t('valid.hysteriaUpMbps'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 1,
            max: 9999999999,
            message: this.$t('valid.hysteriaUpMbpsRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaDownMbps: [
          {
            required: true,
            message: this.$t('valid.hysteriaDownMbps'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 1,
            max: 9999999999,
            message: this.$t('valid.hysteriaDownMbpsRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaServerName: [
          {
            min: 0,
            max: 64,
            message: this.$t('valid.hysteriaServerNameRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaInsecure: [
          {
            required: true,
            message: this.$t('valid.hysteriaInsecure'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaFastOpen: [
          {
            required: true,
            message: this.$t('valid.hysteriaFastOpen'),
            trigger: ['change', 'blur']
          }
        ],
        hysteria2ObfsPassword: [
          {
            min: 4,
            max: 64,
            message: this.$t('valid.hysteria2ObfsPasswordRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteria2UpMbps: [
          {
            required: true,
            message: this.$t('valid.hysteria2UpMbps'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 1,
            max: 9999999999,
            message: this.$t('valid.hysteria2UpMbpsRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteria2DownMbps: [
          {
            required: true,
            message: this.$t('valid.hysteria2DownMbps'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 1,
            max: 9999999999,
            message: this.$t('valid.hysteria2DownMbpsRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteria2ServerName: [
          {
            min: 0,
            max: 64,
            message: this.$t('valid.hysteria2ServerNameRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteria2Insecure: [
          {
            required: true,
            message: this.$t('valid.hysteria2Insecure'),
            trigger: ['change', 'blur']
          }
        ]
      },
      updateRules: {
        name: [
          {
            required: true,
            message: this.$t('valid.nodeName'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 20,
            message: this.$t('valid.nodeNameRange'),
            trigger: ['change', 'blur']
          }
        ],
        nodeServerId: [
          {
            required: true,
            message: this.$t('valid.nodeServerId'),
            trigger: ['change', 'blur']
          }
        ],
        domain: [
          {
            required: true,
            message: this.$t('valid.nodeDomain'),
            trigger: ['change', 'blur']
          },
          {
            min: 4,
            max: 64,
            message: this.$t('valid.nodeDomainRange'),
            trigger: ['change', 'blur']
          }
        ],
        port: [
          {
            required: true,
            message: this.$t('valid.nodePort'),
            trigger: ['change', 'blur']
          },
          {
            pattern:
              /^([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/,
            message: this.$t('valid.nodePortRange'),
            trigger: ['change', 'blur']
          }
        ],
        priority: [
          {
            required: true,
            message: this.$t('valid.nodePriority'),
            trigger: ['change', 'blur']
          },
          {
            pattern: /^-?\d+$/,
            message: this.$t('valid.nodePriorityRange'),
            trigger: ['change', 'blur']
          }
        ],
        nodeTypeId: [
          {
            required: true,
            message: this.$t('valid.nodeType'),
            trigger: ['change', 'blur']
          }
        ],
        xrayProtocol: [
          {
            required: true,
            message: this.$t('valid.xrayProtocol'),
            trigger: ['change', 'blur']
          }
        ],
        xraySSMethod: [
          {
            required: true,
            message: this.$t('valid.xraySSMethod'),
            trigger: ['change', 'blur']
          }
        ],
        realityPbk: [
          {
            required: true,
            message: this.$t('valid.realityPbk'),
            trigger: ['change', 'blur']
          }
        ],
        'xraySettingsEntity.network': [
          {
            required: true,
            message: this.$t('valid.xraySSNetwork'),
            trigger: ['change', 'blur']
          }
        ],
        'xraySettingsEntity.udp': [
          {
            required: true,
            message: this.$t('valid.xraySocksUdp'),
            trigger: ['change', 'blur']
          }
        ],
        'xraySettingsEntity.accounts[0].user': [
          {
            required: true,
            message: this.$t('valid.xraySocksUser'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 64,
            message: this.$t('valid.xraySocksUserRange'),
            trigger: ['change', 'blur']
          }
        ],
        'xraySettingsEntity.accounts[0].pass': [
          {
            required: true,
            message: this.$t('valid.xraySocksPass'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 64,
            message: this.$t('valid.xraySocksPassRange'),
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.network': [
          {
            required: true,
            message: this.$t('valid.xrayNetwork'),
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.security': [
          {
            required: true,
            message: this.$t('valid.xraySecurity'),
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.wsSettings.path': [
          {
            min: 2,
            max: 64,
            message: this.$t('valid.xrayWsRange'),
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.wsSettings.headers.Host': [
          {
            min: 2,
            max: 64,
            message: this.$t('valid.xrayWsRange'),
            trigger: ['change', 'blur']
          }
        ],
        'xrayStreamSettingsEntity.realitySettings.dest': {
          validator: validateXrayStreamSettingsEntityRealitySettingsDest,
          trigger: ['change', 'blur']
        },
        'xrayStreamSettingsEntity.realitySettings.serverNames': {
          validator: validateXrayStreamSettingsEntityRealitySettingsServerNames,
          trigger: ['change', 'blur']
        },
        'xrayStreamSettingsEntity.realitySettings.privateKey': {
          validator: validateXrayStreamSettingsEntityRealitySettingsPrivateKey,
          trigger: ['change', 'blur']
        },
        'xrayStreamSettingsEntity.realitySettings.shortIds': {
          validator: validateXrayStreamSettingsEntityRealitySettingsShortIds,
          trigger: ['change', 'blur']
        },
        trojanGoSni: [
          {
            min: 4,
            max: 64,
            message: this.$t('valid.trojanGoSni'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoMuxEnable: [
          {
            required: true,
            message: this.$t('valid.trojanGoMux'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoWebsocketEnable: [
          {
            required: true,
            message: this.$t('valid.trojanGoWs'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoWebsocketPath: [
          {
            min: 2,
            max: 64,
            message: this.$t('valid.trojanGoWsRange'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoWebsocketHost: [
          {
            min: 2,
            max: 64,
            message: this.$t('valid.trojanGoWsHostRange'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoSsEnable: [
          {
            required: true,
            message: this.$t('valid.trojanGoSs'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoSsMethod: [
          {
            required: true,
            message: this.$t('valid.trojanGoSsMethod'),
            trigger: ['change', 'blur']
          }
        ],
        trojanGoSsPassword: [
          {
            min: 2,
            max: 32,
            message: this.$t('valid.trojanGoSsPasswordRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaProtocol: [
          {
            required: true,
            message: this.$t('valid.hysteriaProtocol'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 16,
            message: this.$t('valid.hysteriaProtocolRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaObfs: [
          {
            min: 0,
            max: 64,
            message: this.$t('valid.hysteriaObfsRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaUpMbps: [
          {
            required: true,
            message: this.$t('valid.hysteriaUpMbps'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 1,
            max: 9999999999,
            message: this.$t('valid.hysteriaUpMbpsRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaDownMbps: [
          {
            required: true,
            message: this.$t('valid.hysteriaDownMbps'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 1,
            max: 9999999999,
            message: this.$t('valid.hysteriaDownMbpsRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaServerName: [
          {
            min: 0,
            max: 64,
            message: this.$t('valid.hysteriaServerNameRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaInsecure: [
          {
            required: true,
            message: this.$t('valid.hysteriaInsecure'),
            trigger: ['change', 'blur']
          }
        ],
        hysteriaFastOpen: [
          {
            required: true,
            message: this.$t('valid.hysteriaFastOpen'),
            trigger: ['change', 'blur']
          }
        ],
        hysteria2ObfsPassword: [
          {
            min: 4,
            max: 64,
            message: this.$t('valid.hysteria2ObfsPasswordRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteria2UpMbps: [
          {
            required: true,
            message: this.$t('valid.hysteria2UpMbps'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 1,
            max: 9999999999,
            message: this.$t('valid.hysteria2UpMbpsRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteria2DownMbps: [
          {
            required: true,
            message: this.$t('valid.hysteria2DownMbps'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 1,
            max: 9999999999,
            message: this.$t('valid.hysteria2DownMbpsRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteria2ServerName: [
          {
            min: 0,
            max: 64,
            message: this.$t('valid.hysteria2ServerNameRange'),
            trigger: ['change', 'blur']
          }
        ],
        hysteria2Insecure: [
          {
            required: true,
            message: this.$t('valid.hysteria2Insecure'),
            trigger: ['change', 'blur']
          }
        ]
      }
    }
  },
  methods: {
    isXray,
    isTrojanGo,
    isHysteria,
    isNaiveProxy,
    isHysteria2,
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.nodeProps.xraySettings = handleXraySettings(this.nodeProps)
          this.nodeProps.xrayStreamSettings = handleXrayStreamSettings(
            this.nodeProps
          )
          createNode(this.nodeProps).then(() => {
            this.getListProps()
            this.$emit('update:dialogFormVisibleProps', false)
            this.$notify({
              title: 'Success',
              message: this.$t('confirm.createSuccess').toString(),
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.nodeProps.xraySettings = handleXraySettings(this.nodeProps)
          this.nodeProps.xrayStreamSettings = handleXrayStreamSettings(
            this.nodeProps
          )
          const nodePropsData = Object.assign({}, this.nodeProps)
          updateNodeById(nodePropsData).then(() => {
            this.getListProps()
            this.$emit('update:dialogFormVisibleProps', false)
            this.$notify({
              title: 'Success',
              message: this.$t('confirm.modifySuccess').toString(),
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleCreateFallback() {
      this.$refs.fallbackForm.resetTemp()
      this.dialogFallbackFormVisible = true
    },
    handleFallbackDetail(fallback) {
      this.dialogFallbackDetailVisible = true
      this.fallback = fallback
    },
    deleteFallback(fallback) {
      for (
        let i = 0;
        i < this.nodeProps.xraySettingsEntity.fallbacks.length;
        i++
      ) {
        let nodeProps = this.nodeProps.xraySettingsEntity.fallbacks[i]
        if (
          fallback.alpn === nodeProps.alpn &&
          fallback.name === nodeProps.name &&
          fallback.path === nodeProps.path &&
          fallback.dest === nodeProps.dest &&
          fallback.xver === nodeProps.xver
        ) {
          this.nodeProps.xraySettingsEntity.fallbacks.splice(i, 1)
        }
      }
    },
    createFallback(fallback) {
      this.nodeProps.xraySettingsEntity.fallbacks.push(fallback)
    },
    toAddNodeServer() {
      this.$router.push({ path: '/server-manage/server-list' })
    }
  }
}
</script>

<style scoped>
.el-button {
  margin-left: 10px;
}
</style>
