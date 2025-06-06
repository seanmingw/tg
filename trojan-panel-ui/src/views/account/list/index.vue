<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.username"
        :placeholder="$t('table.username')"
        style="width: 200px"
        class="filter-item"
        clearable
        @keyup.enter.native="handleFilter"
        @clear="handleFilter"
        maxlength="20"
      />
      <el-select
        v-model="listQuery.deleted"
        :placeholder="$t('table.deleted').toString()"
        style="width: 200px"
        class="filter-item"
        clearable
        @clear="
          () => {
            listQuery.deleted = undefined
          }
        "
      >
        <el-option
          :label="item.label"
          :value="item.value"
          :key="item.value"
          v-for="item in deletedList"
        ></el-option>
      </el-select>
      <el-select
        v-model="listQuery.lastLoginTime"
        :placeholder="$t('table.lastLoginTimeUnused').toString()"
        style="width: 200px"
        class="filter-item"
        clearable
        @clear="
          () => {
            listQuery.lastLoginTime = undefined
          }
        "
      >
        <el-option
          :label="item.label"
          :value="item.value"
          :key="item.value"
          v-for="item in lastLoginTimeList"
        ></el-option>
      </el-select>
      <el-select
        v-model="orderFieldArr"
        :placeholder="$t('table.orderFields').toString()"
        style="width: 200px"
        class="filter-item"
        clearable
        multiple
      >
        <el-option
          :label="item.label"
          :value="item.value"
          :key="item.value"
          v-for="item in orderFieldList"
        ></el-option>
      </el-select>
      <el-select
        v-model="listQuery.orderBy"
        :placeholder="$t('table.orderBy').toString()"
        style="width: 200px"
        class="filter-item"
        clearable
        @clear="
          () => {
            listQuery.orderBy = undefined
          }
        "
      >
        <el-option
          :label="item.label"
          :value="item.value"
          :key="item.value"
          v-for="item in orderByList"
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
        v-if="checkPermission(['sysadmin'])"
      >
        {{ $t('table.add') }}
      </el-button>
      <el-button
        class="filter-item"
        type="success"
        icon="el-icon-upload2"
        @click="handleImport"
        v-if="checkPermission(['sysadmin'])"
      >
        {{ $t('table.import') }}
      </el-button>
      <el-button
        class="filter-item"
        type="success"
        icon="el-icon-download"
        @click="handleExport"
        v-if="checkPermission(['sysadmin'])"
      >
        {{ $t('table.export') }}
      </el-button>
      <el-button
        class="filter-item"
        type="success"
        icon="el-icon-document-copy"
        @click="handleCreateBatch"
        v-if="checkPermission(['sysadmin'])"
      >
        {{ $t('table.createBatch') }}
      </el-button>
      <el-button
        class="filter-item"
        type="success"
        icon="el-icon-download"
        @click="exportAccountUnused"
        v-if="checkPermission(['sysadmin'])"
      >
        {{ $t('table.exportAccountUnused') }}
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
        :label="$t('table.username').toString()"
        width="100"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ row.username }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.role').toString()"
        width="100"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ roleFilter(row.roleId) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.download').toString()"
        width="100"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ getFlow(row.download) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.upload').toString()"
        width="100"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ getFlow(row.upload) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.flow').toString()"
        width="100"
        align="center"
      >
        <template slot-scope="{ row }">
          <span
            :style="
              row.quota >= 0 && row.quota - row.upload - row.download <= 0
                ? 'color: #FF0000;'
                : ''
            "
            >{{
              row.quota < 0
                ? $t('dashboard.unlimited')
                : getFlow(row.quota - row.upload - row.download)
            }}
          </span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.quota').toString()"
        width="100"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{
            row.lastLoginTime === 0
              ? '-'
              : row.quota < 0
              ? $t('dashboard.unlimited')
              : getFlow(row.quota)
          }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.expireTime').toString()"
        width="150"
        align="center"
      >
        <template slot-scope="{ row }">
          <span
            :style="
              row.lastLoginTime !== 0 && row.expireTime <= new Date().getTime()
                ? 'color: #FF0000;'
                : ''
            "
            >{{
              row.lastLoginTime === 0
                ? '-'
                : timeStampToDate(row.expireTime, false)
            }}
          </span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.presetQuota').toString()"
        width="100"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{
            row.presetExpire === 0
              ? '-'
              : row.presetQuota < 0
              ? $t('dashboard.unlimited')
              : getFlow(row.presetQuota)
          }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.presetExpire').toString()"
        width="80"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ row.presetExpire === 0 ? '-' : row.presetExpire }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.status').toString()"
        width="80"
        align="center"
      >
        <template slot-scope="{ row }">
          <el-tag :type="row.deleted | deletedFilter">
            {{ row.deleted === 1 ? $t('table.disable') : $t('table.enable') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.lastLoginTime').toString()"
        width="150"
        align="center"
      >
        <template slot-scope="{ row }">
          <span
            >{{
              row.lastLoginTime === 0
                ? '-'
                : timeStampToDate(row.lastLoginTime, false)
            }}
          </span>
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
            @click="handleReset(row)"
            v-if="checkPermission(['sysadmin'])"
          >
            {{ $t('table.reset') }}
          </el-button>
          <el-button
            type="primary"
            size="mini"
            @click="handleClashSubscribeForSb(row)"
            v-if="
              checkPermission(['sysadmin', 'admin']) && row.lastLoginTime !== 0
            "
          >
            {{ $t('table.clashSubscribeForSb') }}
          </el-button>
          <el-button
            type="primary"
            size="mini"
            @click="handleUpdate(row)"
            v-if="checkPermission(['sysadmin'])"
          >
            {{ $t('table.edit') }}
          </el-button>
          <el-button
            size="mini"
            type="danger"
            @click="handleDelete(row, $index)"
            v-if="checkPermission(['sysadmin'])"
          >
            {{ $t('table.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      v-if="total > 0"
      :total="total"
      :page.sync="listQuery.pageNum"
      :limit.sync="listQuery.pageSize"
      @pagination="getList"
    />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form
        ref="dataForm"
        :rules="dialogStatus === 'create' ? createRules : updateRules"
        :model="temp"
        label-position="left"
      >
        <el-form-item
          v-if="dialogStatus === 'create'"
          :label="$t('table.username')"
          prop="username"
        >
          <el-input
            v-model="temp.username"
            :placeholder="$t('table.username')"
            clearable
          />
        </el-form-item>
        <el-form-item :label="$t('table.pass')" prop="pass">
          <el-input
            v-model="temp.pass"
            type="password"
            :placeholder="$t('table.pass')"
            clearable
          />
        </el-form-item>
        <el-form-item :label="$t('table.editQuota')" prop="quota">
          <el-input-number
            v-model.number="temp.quota"
            controls-position="right"
            type="number"
          />
        </el-form-item>
        <el-form-item :label="$t('table.ipLimit')" prop="ipLimit">
          <el-input-number
            v-model.number="temp.ipLimit"
            controls-position="right"
            type="number"
            :min="0"
            :max="100"
          />
        </el-form-item>
        <el-form-item :label="$t('table.email')" prop="email">
          <el-input
            v-model="temp.email"
            :placeholder="$t('table.email')"
            clearable
          />
        </el-form-item>
        <el-form-item :label="$t('table.status')" prop="deleted">
          <el-switch
            v-model="temp.deleted"
            active-color="#13ce66"
            inactive-color="#ff4949"
            :active-text="$t('table.enable')"
            :inactive-text="$t('table.disable')"
            :active-value="0"
            :inactive-value="1"
          >
          </el-switch>
        </el-form-item>
        <el-form-item :label="$t('table.expireTime')" prop="expireTime">
          <el-date-picker
            v-model="temp.expireTime"
            type="datetime"
            value-format="timestamp"
            :placeholder="$t('table.expireTime')"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false"
          >{{ $t('table.cancel') }}
        </el-button>
        <el-button
          type="primary"
          @click="dialogStatus === 'create' ? createData() : updateData()"
        >
          {{ $t('table.confirm') }}
        </el-button>
      </div>
    </el-dialog>
    <import-tip
      ref="importTip"
      :dialog-form-visible.sync="importVisible"
      :label="$t('table.coverByAccountName')"
      :import-data="importData"
      :download-template="downloadTemplate"
    />
    <BatchOperation
      ref="batchOperationForm"
      :dialog-form-visible-props.sync="batchOperationDialogFormVisible"
      :get-list-props="getList"
    />
  </div>
</template>

<script>
import {
  clashSubscribeForSb,
  createAccount,
  deleteAccountById,
  exportAccount,
  exportAccountUnused,
  importAccount,
  resetAccountDownloadAndUpload,
  selectAccountPage,
  updateAccountById
} from '@/api/account'
import Pagination from '@/components/Pagination'
import ImportTip from '@/components/ImportTip'

import { Message, MessageBox } from 'element-ui'
import { timeStampToDate } from '@/utils'
import { byteToMb, getFlow, mbToByte } from '@/utils/account'
import { selectRoleList } from '@/api/role'
import checkPermission from '@/utils/permission'
import { setting } from '@/api/system'
import { downloadTemplate } from '@/api/file-task'
import BatchOperation from '@/views/account/list/compoments/BatchOperation'
import copy from 'copy-to-clipboard'

export default {
  name: 'List',
  filters: {
    deletedFilter(deleted) {
      const deletedMap = {
        0: 'success',
        1: 'danger'
      }
      return deletedMap[deleted]
    }
  },
  components: { BatchOperation, Pagination, ImportTip },
  data() {
    const validateUsername = (rule, value, callback) => {
      if (this.temp.username.trim().indexOf('admin') >= 0) {
        callback(new Error(this.$t('valid.usernameNotExistAdmin')))
      } else {
        callback()
      }
    }
    return {
      tableKey: 0,
      listLoading: true,
      list: null,
      total: 0,
      orderFieldArr: ['role_id', 'create_time'],
      listQuery: {
        pageNum: 1,
        pageSize: 20,
        username: undefined,
        deleted: undefined,
        orderFields: 'role_id,create_time',
        orderBy: 'desc',
        lastLoginTime: undefined
      },
      temp: {
        id: undefined,
        quota: 0,
        download: 0,
        upload: 0,
        username: undefined,
        pass: undefined,
        email: undefined,
        roleId: 3,
        deleted: 0,
        lastLoginTime: 0,
        expireTime: new Date().getTime(),
        createTime: new Date(),
        ipLimit: 3 // 默认设备数限制为3
      },
      dialogFormVisible: false,
      batchOperationDialogFormVisible: false,
      textMap: {
        update: this.$t('table.edit'),
        create: this.$t('table.add')
      },
      importVisible: false,
      deletedList: [
        { value: 0, label: this.$t('table.enable') },
        { value: 1, label: this.$t('table.disable') }
      ],
      lastLoginTimeList: [
        { value: 0, label: this.$t('table.no') },
        { value: 1, label: this.$t('table.yes') }
      ],
      orderFieldList: [
        { value: 'quota', label: this.$t('table.quota') },
        { value: 'role_id', label: this.$t('table.role') },
        { value: 'last_login_time', label: this.$t('table.lastLoginTime') },
        { value: 'expire_time', label: this.$t('table.expireTime') },
        { value: 'deleted', label: this.$t('table.deleted') },
        { value: 'create_time', label: this.$t('table.createTime') }
      ],
      orderByList: [
        { value: 'desc', label: this.$t('table.desc') },
        { value: 'asc', label: this.$t('table.asc') }
      ],
      createRules: {
        username: [
          {
            required: true,
            message: this.$t('valid.username'),
            trigger: ['change', 'blur']
          },
          {
            min: 6,
            max: 20,
            message: this.$t('valid.usernameRange'),
            trigger: ['change', 'blur']
          },
          {
            pattern: /^[A-Za-z0-9]+$/,
            message: this.$t('valid.usernameElement'),
            trigger: ['change', 'blur']
          },
          {
            validator: validateUsername,
            trigger: ['change', 'blur']
          }
        ],
        pass: [
          {
            required: true,
            message: this.$t('valid.pass'),
            trigger: ['change', 'blur']
          },
          {
            min: 6,
            max: 20,
            message: this.$t('valid.passRange'),
            trigger: ['change', 'blur']
          },
          {
            pattern: /^[A-Za-z0-9]+$/,
            message: this.$t('valid.passElement'),
            trigger: ['change', 'blur']
          }
        ],
        quota: [
          {
            required: true,
            message: this.$t('valid.quota'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: -1,
            max: 1024000,
            message: this.$t('valid.quotaRange'),
            trigger: ['change', 'blur']
          }
        ],
        ipLimit: [
          {
            required: true,
            message: this.$t('valid.ipLimit'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 0,
            max: 100,
            message: this.$t('valid.ipLimitRange'),
            trigger: ['change', 'blur']
          }
        ],
        email: [
          {
            min: 4,
            max: 64,
            message: this.$t('valid.emailRange'),
            trigger: ['change', 'blur']
          },
          {
            pattern:
              /^([A-Za-z0-9_\-\.])+\@(163.com|126.com|qq.com|gmail.com)$/,
            message: this.$t('valid.emailElement'),
            trigger: ['change', 'blur']
          }
        ],
        deleted: [
          {
            required: true,
            message: this.$t('valid.deleted'),
            trigger: ['change', 'blur']
          }
        ],
        expireTime: [
          {
            required: true,
            message: this.$t('valid.expireTime'),
            trigger: ['change', 'blur']
          }
        ]
      },
      updateRules: {
        username: [
          {
            required: true,
            message: this.$t('valid.username'),
            trigger: ['change', 'blur']
          },
          {
            min: 6,
            max: 20,
            message: this.$t('valid.usernameRange'),
            trigger: ['change', 'blur']
          },
          {
            pattern: /^[A-Za-z0-9]+$/,
            message: this.$t('valid.usernameElement'),
            trigger: ['change', 'blur']
          },
          {
            validator: validateUsername,
            trigger: ['change', 'blur']
          }
        ],
        pass: [
          {
            min: 6,
            max: 20,
            message: this.$t('valid.passRange'),
            trigger: ['change', 'blur']
          },
          {
            pattern: /^[A-Za-z0-9]+$/,
            message: this.$t('valid.passElement'),
            trigger: ['change', 'blur']
          }
        ],
        quota: [
          {
            required: true,
            message: this.$t('valid.quota'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: -1,
            max: 1024000,
            message: this.$t('valid.quotaRange'),
            trigger: ['change', 'blur']
          }
        ],
        ipLimit: [
          {
            required: true,
            message: this.$t('valid.ipLimit'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 0,
            max: 100,
            message: this.$t('valid.ipLimitRange'),
            trigger: ['change', 'blur']
          }
        ],
        email: [
          {
            min: 4,
            max: 64,
            message: this.$t('valid.emailRange'),
            trigger: ['change', 'blur']
          },
          {
            pattern:
              /^([A-Za-z0-9_\-\.])+\@(163.com|126.com|qq.com|gmail.com)$/,
            message: this.$t('valid.emailElement'),
            trigger: ['change', 'blur']
          }
        ],
        deleted: [
          {
            required: true,
            message: this.$t('valid.deleted'),
            trigger: ['change', 'blur']
          }
        ],
        expireTime: [
          {
            required: true,
            message: this.$t('valid.expireTime'),
            trigger: ['change', 'blur']
          }
        ]
      },
      dialogStatus: '',
      roleList: []
    }
  },
  created() {
    this.setRoleList()
    this.getList()
  },
  methods: {
    checkPermission,
    getFlow,
    timeStampToDate,
    setRoleList() {
      selectRoleList().then((response) => {
        const { data } = response
        this.roleList = data
      })
    },
    getList() {
      this.listLoading = true
      this.listQuery.orderFields = this.orderFieldArr.join(',')
      selectAccountPage(this.listQuery).then((response) => {
        this.list = response.data.accounts
        this.total = response.data.total

        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        quota: 0,
        download: 0,
        upload: 0,
        username: undefined,
        pass: undefined,
        email: undefined,
        roleId: 3,
        deleted: 0,
        expireTime: new Date().getTime(),
        createTime: new Date(),
        ipLimit: 3 // 默认设备数限制为3
      }
    },
    handleFilter() {
      this.listQuery.pageNum = 1
      this.getList()
    },
    handleCreate() {
      this.resetTemp()
      setting().then((response) => {
        if (response.data.registerEnable === 1) {
          this.temp.quota = response.data.registerQuota
          this.temp.expireTime =
            new Date().getTime() + response.data.registerExpireDays * 86400000
        }
      })
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row)
      this.temp.quota = byteToMb(row.quota)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleDelete(row, index) {
      MessageBox.confirm(
        this.$t('confirm.deleteUser'),
        this.$t('confirm.warn'),
        {
          confirmButtonText: this.$t('confirm.yes'),
          cancelButtonText: this.$t('confirm.cancel'),
          type: 'warning'
        }
      ).then(() => {
        const tempData = Object.assign({}, row)
        deleteAccountById(tempData).then(() => {
          this.list.splice(index, 1)
          this.$notify({
            title: 'Success',
            message: this.$t('confirm.deleteSuccess'),
            type: 'success',
            duration: 2000
          })
        })
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          // Ensure ipLimit field has a value
          if (this.temp.ipLimit === undefined) {
            this.temp.ipLimit = 3; // Default value
          }
          createAccount(this.temp).then(() => {
            this.getList()
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: this.$t('confirm.createSuccess'),
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
          // 确保ipLimit字段有值
          if (this.temp.ipLimit === undefined) {
            this.temp.ipLimit = 3; // 默认值
          }
          const tempData = Object.assign({}, this.temp)
          updateAccountById(tempData).then(() => {
            const index = this.list.findIndex((v) => v.id === this.temp.id)
            this.list.splice(index, 1, this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: this.$t('notify.success'),
              message: this.$t('notify.updateSuccess'),
              type: 'success',
              duration: 2000
            })
            this.getList()
          })
        }
      })
    },
    roleFilter(roleId) {
      let role = this.roleList.find((item) => item.id === roleId)
      if (role && role.desc) {
        return role.desc
      } else {
        return ''
      }
    },
    handleReset(row) {
      MessageBox.confirm(
        this.$t('confirm.handleReset'),
        this.$t('confirm.warn'),
        {
          confirmButtonText: this.$t('confirm.yes'),
          cancelButtonText: this.$t('confirm.cancel'),
          type: 'warning'
        }
      ).then(() => {
        const tempData = Object.assign({}, row)
        resetAccountDownloadAndUpload(tempData).then(() => {
          this.getList()
          this.$notify({
            title: 'Success',
            message: this.$t('confirm.resetSuccess'),
            type: 'success',
            duration: 2000
          })
        })
      })
    },
    importData(params) {
      this.$refs['importTip'].$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.$refs['importTip'].temp)
          let formData = new FormData()
          formData.append('file', params.file)
          formData.append('cover', tempData.cover)
          importAccount(formData).then(() => {
            this.importVisible = false
            this.$notify({
              title: 'Success',
              message: this.$t('confirm.taskSubmitSuccess'),
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleImport() {
      this.importVisible = true
    },
    handleExport() {
      exportAccount().then(() => {
        this.importVisible = false
        this.$notify({
          title: 'Success',
          message: this.$t('confirm.taskSubmitSuccess').toString(),
          type: 'success',
          duration: 2000
        })
      })
    },
    handleCreateBatch() {
      this.batchOperationDialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['batchOperationForm'].$refs['dataForm'].clearValidate()
      })
    },
    exportAccountUnused() {
      exportAccountUnused().then(() => {
        this.$notify({
          title: 'Success',
          message: this.$t('confirm.taskSubmitSuccess').toString(),
          type: 'success',
          duration: 2000
        })
      })
    },
    downloadTemplate() {
      downloadTemplate({ id: 1 }).then((res) => {
        // 将二进制文件转化为可访问的url
        const blob = new Blob([res.data], {
          type: 'application/octet-stream'
        })
        let url = window.URL.createObjectURL(blob)
        let a = document.createElement('a')
        document.body.appendChild(a)
        a.href = url
        let dis = res.headers['content-disposition']
        a.download = dis.split('attachment; filename=')[1]
        // 模拟点击下载
        a.click()
        window.URL.revokeObjectURL(url)
        this.$notify({
          title: 'Success',
          message: this.$t('confirm.taskDownloadSuccess'),
          type: 'success',
          duration: 2000
        })
      })
    },
    handleClashSubscribeForSb(row) {
      clashSubscribeForSb(row).then((response) => {
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
