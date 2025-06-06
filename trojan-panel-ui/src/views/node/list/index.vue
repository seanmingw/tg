<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.name"
        :placeholder="$t('table.nodeName')"
        style="width: 200px"
        class="filter-item"
        clearable
        @keyup.enter.native="handleFilter"
        @clear="handleFilter"
      />
      <el-select
        v-model="listQuery.nodeServerId"
        :placeholder="$t('table.nodeServerName').toString()"
        class="filter-item"
        @change="handleFilter"
        @clear="handleFilter"
      >
        <el-option
          :label="item.name"
          :value="item.id"
          :key="item.id"
          v-for="item in nodeServers"
        ></el-option>
      </el-select>
      <el-button
        class="filter-item"
        type="primary"
        icon="el-icon-search"
        @click="handleFilter"
      >
        {{ $t('table.search') }}
      </el-button>
      <el-button
        class="filter-item"
        type="primary"
        icon="el-icon-edit"
        @click="handleCreate"
        v-if="checkPermission(['sysadmin', 'admin'])"
      >
        {{ $t('table.add') }}
      </el-button>
      <el-button
        class="filter-item"
        type="success"
        icon="el-icon-link"
        @click="handleClashSubscribe"
      >
        {{ $t('table.clashSubscribe') }}
      </el-button>
    </div>
    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%"
    >
      <el-table-column
        :label="$t('table.id').toString()"
        sortable="custom"
        align="center"
        width="80"
        type="index"
      />
      <el-table-column
        :label="$t('table.nodeName').toString()"
        width="120"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.nodeServerName').toString()"
        width="150"
        align="center"
        v-if="checkPermission(['sysadmin', 'admin'])"
      >
        <template slot-scope="{ row }">
          <span>{{ nodeServerFind(nodeServers, row.nodeServerId) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.nodeDomain').toString()"
        width="150"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ row.domain }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.nodePort').toString()"
        width="80"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ row.port }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.nodePriority').toString()"
        width="80"
        align="center"
        v-if="checkPermission(['sysadmin', 'admin'])"
      >
        <template slot-scope="{ row }">
          <span>{{ row.priority }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.nodeType').toString()"
        width="120"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ nodeTypeFind(nodeTypes, row.nodeTypeId) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.nodeStatus').toString()"
        width="100"
        align="center"
        v-if="checkPermission(['sysadmin', 'admin'])"
      >
        <template slot-scope="{ row }">
          <el-tag :type="row.status | statusTypeFilter">
            <span>{{ statusComputed(row.status) }}</span>
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.createTime').toString()"
        width="150"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ timeStampToDate(row.createTime, false) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.actions').toString()"
        align="center"
        class-name="small-padding fixed-width"
      >
        <template slot-scope="{ row, $index }">
          <el-button
            type="primary"
            size="mini"
            @click="handleUpdate(row)"
            v-if="checkPermission(['sysadmin', 'admin'])"
          >
            {{ $t('table.edit') }}
          </el-button>
          <el-button type="primary" size="mini" @click="handleQRCode(row)">
            {{ $t('table.nodeQRCode') }}
          </el-button>
          <el-button type="success" size="mini" @click="handleCopyURL(row)">
            {{ $t('table.nodeURL') }}
          </el-button>
          <el-button type="success" size="mini" @click="handleDetail(row)">
            {{ $t('table.detail') }}
          </el-button>
          <el-button
            size="mini"
            type="danger"
            @click="handleDelete(row, $index)"
            v-if="checkPermission(['sysadmin', 'admin'])"
          >
            {{ $t('table.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <Pagination
      v-if="total > 0"
      :total="total"
      :page.sync="listQuery.pageNum"
      :limit.sync="listQuery.pageSize"
      @pagination="getList"
    />

    <NodeForm
      ref="nodeForm"
      :dialog-form-visible-props.sync="dialogFormVisible"
      :node-props="temp"
      :dialog-status-props="dialogStatus"
      :node-servers-props="nodeServers"
      :node-types-props="nodeTypes"
      :get-list-props="getList"
    />

    <NodeDetail
      :dialog-visible-props.sync="dialogInfoVisible"
      :node-info-props="nodeDetail"
      :node-servers-props="nodeServers"
      :node-types-props="nodeTypes"
    />

    <NodeQrcode
      :dialog-visible-props.sync="dialogQRCodeVisible"
      :qr-code-src-props="qrCodeSrc"
    />
  </div>
</template>

<script>
import Pagination from '@/components/Pagination'
import NodeDetail from '@/views/node/list/components/NodeDetail'
import NodeQrcode from '@/views/node/list/components/NodeQrcode'
import { Message, MessageBox } from 'element-ui'
import copy from 'copy-to-clipboard'
import {
  deleteNodeById,
  nodeDefault,
  nodeQRCode,
  nodeURL,
  selectNodeById,
  selectNodeInfo,
  selectNodePage
} from '@/api/node'
import { selectNodeTypeList } from '@/api/node-type'
import {
  handleNodeDetail,
  handleNodeUpdate,
  nodeServerFind,
  nodeTypeFind
} from '@/utils/node'
import checkPermission from '@/utils/permission'
import { timeStampToDate } from '@/utils'
import { clashSubscribe } from '@/api/account'
import { selectNodeServerList } from '@/api/node-server'
import FallbackForm from '@/views/node/list/components/FallbackForm'
import FallbackInfo from '@/views/node/list/components/FallbackInfo'
import NodeForm from '@/views/node/list/components/NodeForm'

export default {
  name: 'index',
  components: {
    NodeForm,
    FallbackInfo,
    FallbackForm,
    NodeQrcode,
    NodeDetail,
    Pagination
  },
  filters: {
    statusTypeFilter(status) {
      return status > 0 ? 'success' : 'danger'
    }
  },
  computed: {
    statusComputed() {
      return function (status) {
        return status === 1
          ? this.$t('table.nodeStatusSuccess')
          : this.$t('table.nodeStatusError')
      }
    }
  },
  data() {
    return {
      tableKey: 0,
      listLoading: true,
      list: null,
      total: 0,
      listQuery: {
        pageNum: 1,
        pageSize: 20,
        name: undefined,
        nodeServerId: undefined
      },
      temp: {
        id: undefined,
        nodeServerId: undefined,
        nodeSubId: undefined,
        nodeTypeId: 1,
        name: '',
        domain: '',
        port: 443,
        priority: 100,

        xrayProtocol: 'vless',
        xrayFlow: '',
        xraySSMethod: 'aes-256-gcm',
        realityPbk: '',
        xraySettings: '',
        xraySettingsEntity: {
          fallbacks: [
            {
              name: '',
              alpn: '',
              path: undefined,
              dest: '80',
              xver: 0
            }
          ],
          network: 'tcp',
          accounts: [{ user: '', pass: '' }],
          udp: true
        },
        xrayStreamSettings: '',
        xrayStreamSettingsEntity: {
          network: 'tcp',
          security: 'none',
          tlsSettings: {
            serverName: '',
            alpn: ['h2', 'http/1.1'],
            allowInsecure: false,
            fingerprint: 'chrome'
          },
          realitySettings: {
            dest: '',
            xver: 0,
            serverNames: [],
            fingerprint: 'chrome',
            privateKey: '',
            shortIds: [],
            spiderX: ''
          },
          wsSettings: {
            path: '/trojan-panel-websocket-path',
            headers: {
              Host: ''
            }
          }
        },
        xrayTag: 'user',
        xraySniffing: '',
        xrayAllocate: '',

        trojanGoSni: '',
        trojanGoMuxEnable: 1,
        trojanGoWebsocketEnable: 0,
        trojanGoWebsocketPath: '/trojan-panel-websocket-path',
        trojanGoWebsocketHost: '',
        trojanGoSsEnable: 0,
        trojanGoSsMethod: 'AES-128-GCM',
        trojanGoSsPassword: '',

        hysteriaProtocol: 'udp',
        hysteriaObfs: '',
        hysteriaUpMbps: 100,
        hysteriaDownMbps: 100,
        hysteriaServerName: '',
        hysteriaInsecure: 0,
        hysteriaFastOpen: 0,

        hysteria2ObfsPassword: '',
        hysteria2UpMbps: 100,
        hysteria2DownMbps: 100,
        hysteria2ServerName: '',
        hysteria2Insecure: 0,
        createTime: new Date()
      },
      nodeDetail: {
        id: undefined,
        nodeServerId: undefined,
        nodeSubId: undefined,
        nodeTypeId: 1,
        name: '',
        domain: '',
        port: 443,
        priority: 100,

        password: '',
        uuid: '',
        alterId: 0,

        xrayProtocol: 'vless',
        xrayFlow: '',
        xraySSMethod: 'aes-256-gcm',
        realityPbk: '',
        xraySettings: '',
        xraySettingsEntity: {
          fallbacks: [
            {
              name: '',
              alpn: '',
              path: undefined,
              dest: '80',
              xver: 0
            }
          ],
          network: 'tcp',
          accounts: [{ user: '', pass: '' }],
          udp: true
        },
        xrayStreamSettings: '',
        xrayStreamSettingsEntity: {
          network: 'tcp',
          security: 'none',
          tlsSettings: {
            serverName: '',
            alpn: ['h2', 'http/1.1'],
            allowInsecure: false,
            fingerprint: 'chrome'
          },
          realitySettings: {
            dest: '',
            xver: 0,
            serverNames: [],
            fingerprint: 'chrome',
            privateKey: '',
            shortIds: [],
            spiderX: ''
          },
          wsSettings: {
            path: '/trojan-panel-websocket-path',
            headers: {
              Host: ''
            }
          }
        },
        xrayTag: 'user',
        xraySniffing: '',
        xrayAllocate: '',

        trojanGoSni: '',
        trojanGoMuxEnable: 1,
        trojanGoWebsocketEnable: 0,
        trojanGoWebsocketPath: '/trojan-panel-websocket-path',
        trojanGoWebsocketHost: '',
        trojanGoSsEnable: 0,
        trojanGoSsMethod: 'AES-128-GCM',
        trojanGoSsPassword: '',

        hysteriaProtocol: 'udp',
        hysteriaObfs: '',
        hysteriaUpMbps: 100,
        hysteriaDownMbps: 100,
        hysteriaServerName: '',
        hysteriaInsecure: 0,
        hysteriaFastOpen: 0,

        hysteria2ObfsPassword: '',
        hysteria2UpMbps: 100,
        hysteria2DownMbps: 100,
        hysteria2ServerName: '',
        hysteria2Insecure: 0,

        naiveProxyUsername: '',
        createTime: new Date()
      },
      dialogStatus: '',
      dialogFormVisible: false,
      dialogInfoVisible: false,
      dialogQRCodeVisible: false,
      qrCodeSrc: '',
      nodeTypes: [],
      nodeServers: [],
      textMap: {
        update: this.$t('table.edit'),
        create: this.$t('table.add')
      }
    }
  },
  created() {
    this.setNodeTypes()
    this.setNodeServers()
    this.getList()
  },
  methods: {
    timeStampToDate,
    checkPermission,
    nodeServerFind,
    nodeTypeFind,
    setNodeTypes() {
      selectNodeTypeList().then((response) => {
        const { data } = response
        this.nodeTypes = data
      })
    },
    setNodeServers() {
      selectNodeServerList().then((response) => {
        const { data } = response
        this.nodeServers = data
      })
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        nodeServerId: undefined,
        nodeSubId: undefined,
        nodeTypeId: 1,
        name: '',
        domain: '',
        port: 443,
        priority: 100,

        xrayProtocol: 'vless',
        xrayFlow: '',
        xraySSMethod: 'aes-256-gcm',
        realityPbk: '',
        xraySettings: '',
        xraySettingsEntity: {
          fallbacks: [
            {
              name: '',
              alpn: '',
              path: undefined,
              dest: '80',
              xver: 0
            }
          ],
          network: 'tcp',
          accounts: [{ user: '', pass: '' }],
          udp: true
        },
        xrayStreamSettings: '',
        xrayStreamSettingsEntity: {
          network: 'tcp',
          security: 'none',
          tlsSettings: {
            serverName: '',
            alpn: ['h2', 'http/1.1'],
            allowInsecure: false,
            fingerprint: 'chrome'
          },
          realitySettings: {
            dest: '',
            xver: 0,
            serverNames: [],
            fingerprint: 'chrome',
            privateKey: '',
            shortIds: [],
            spiderX: ''
          },
          wsSettings: {
            path: '/trojan-panel-websocket-path',
            headers: {
              Host: ''
            }
          }
        },
        xrayTag: 'user',
        xraySniffing: '',
        xrayAllocate: '',

        trojanGoSni: '',
        trojanGoMuxEnable: 1,
        trojanGoWebsocketEnable: 0,
        trojanGoWebsocketPath: '/trojan-panel-websocket-path',
        trojanGoWebsocketHost: '',
        trojanGoSsEnable: 0,
        trojanGoSsMethod: 'AES-128-GCM',
        trojanGoSsPassword: '',

        hysteriaProtocol: 'udp',
        hysteriaObfs: '',
        hysteriaUpMbps: 100,
        hysteriaDownMbps: 100,
        hysteriaServerName: '',
        hysteriaInsecure: 0,
        hysteriaFastOpen: 0,

        hysteria2ObfsPassword: '',
        hysteria2UpMbps: 100,
        hysteria2DownMbps: 100,
        hysteria2ServerName: '',
        hysteria2Insecure: 0,
        createTime: new Date()
      }
    },
    getList() {
      this.listLoading = true
      selectNodePage(this.listQuery).then((response) => {
        this.list = response.data.nodes
        this.total = response.data.total

        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
    },
    handleFilter() {
      this.listQuery.pageNum = 1
      this.getList()
    },
    handleCreate() {
      this.resetTemp()
      // 设置节点部分属性的默认值
      nodeDefault().then((response) => {
        this.temp.realityPbk = response.data.publicKey
        this.temp.xrayStreamSettingsEntity.realitySettings.privateKey =
          response.data.privateKey
        this.temp.xrayStreamSettingsEntity.realitySettings.shortIds = [
          response.data.shortId
        ]
        this.temp.xrayStreamSettingsEntity.realitySettings.spiderX =
          response.data.spiderX
      })
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['nodeForm'].$refs['dataForm'].clearValidate()
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign(this.temp, row)
      selectNodeById({ id: row.id }).then((response) => {
        this.temp = handleNodeUpdate(this.temp, response.data)
        this.dialogStatus = 'update'
        this.dialogFormVisible = true
        this.$nextTick(() => {
          this.$refs['nodeForm'].$refs['dataForm'].clearValidate()
        })
      })
    },
    handleDetail(row) {
      this.nodeDetail = Object.assign(this.temp, row)
      selectNodeInfo({ id: row.id }).then((response) => {
        this.nodeDetail = handleNodeDetail(this.nodeDetail, response.data)
        this.dialogInfoVisible = true
      })
    },
    handleDelete(row, index) {
      MessageBox.confirm(
        this.$t('confirm.deleteNode').toString(),
        this.$t('confirm.warn').toString(),
        {
          confirmButtonText: this.$t('confirm.yes').toString(),
          cancelButtonText: this.$t('confirm.cancel').toString(),
          type: 'warning'
        }
      ).then(() => {
        const tempData = Object.assign({}, row)
        deleteNodeById(tempData).then(() => {
          this.list.splice(index, 1)
          this.$notify({
            title: 'Success',
            message: this.$t('confirm.deleteSuccess').toString(),
            type: 'success',
            duration: 2000
          })
        })
      })
    },
    handleQRCode(row) {
      this.qrCodeSrc = ''
      const tempData = Object.assign({}, row)
      nodeQRCode(tempData).then((response) => {
        this.qrCodeSrc = 'data:image/png;base64,' + response.data
        this.dialogQRCodeVisible = true
      })
    },
    handleCopyURL(row) {
      nodeURL(row).then((response) => {
        if (copy(response.data)) {
          Message({
            showClose: true,
            message: this.$t('confirm.urlCopySuccess').toString(),
            type: 'success'
          })
        } else {
          Message({
            showClose: true,
            message: this.$t('confirm.urlCopyFail').toString(),
            type: 'error'
          })
        }
      })
    },
    handleClashSubscribe() {
      clashSubscribe().then((response) => {
        if (
          copy(
            window.location.protocol +
              '//' +
              window.location.host +
              response.data
          )
        ) {
          Message({
            showClose: true,
            message: this.$t('confirm.urlCopySuccess').toString(),
            type: 'success'
          })
        } else {
          Message({
            showClose: true,
            message: this.$t('confirm.urlCopyFail').toString(),
            type: 'error'
          })
        }
      })
    }
  }
}
</script>

<style scoped>
.el-button {
  margin-left: 10px;
}
</style>