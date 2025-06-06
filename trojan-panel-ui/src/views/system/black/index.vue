<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.ip"
        :placeholder="$t('table.blackListIp')"
        style="width: 200px"
        class="filter-item"
        clearable
        @keyup.enter.native="handleFilter"
        @clear="handleFilter"
        maxlength="64"
      />
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
        :label="$t('table.id')"
        sortable="custom"
        align="center"
        width="80"
        type="index"
      />
      <el-table-column
        :label="$t('table.blackListIp')"
        width="180"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ row.ip }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.createTime')"
        width="150"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ timeStampToDate(row.createTime, false) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.actions')"
        align="center"
        class-name="small-padding fixed-width"
      >
        <template slot-scope="{ row, $index }">
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

    <el-dialog :title="$t('table.add')" :visible.sync="dialogFormVisible">
      <el-form
        ref="dataForm"
        :rules="createRules"
        :model="temp"
        label-position="left"
      >
        <el-form-item :label="$t('table.blackListIp')" prop="ip">
          <el-input v-model="temp.ip" clearable />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false"
          >{{ $t('table.cancel') }}
        </el-button>
        <el-button type="primary" @click="createData()">
          {{ $t('table.confirm') }}
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { timeStampToDate } from '@/utils'
import Pagination from '@/components/Pagination'
import { MessageBox } from 'element-ui'
import {
  createBlackList,
  deleteBlackListByIp,
  selectBlackListPage
} from '@/api/black-list'
import checkPermission from '@/utils/permission' // 权限判断指令

export default {
  name: 'Black',
  components: { Pagination },
  data() {
    const validateIp = (rule, value, callback) => {
      if (this.temp.ip === '127.0.0.1') {
        callback(new Error('IP cannot be 127.0.0.1'))
      } else {
        callback()
      }
    }
    return {
      tableKey: 0,
      listLoading: true,
      list: null,
      total: 0,
      listQuery: {
        pageNum: 1,
        pageSize: 20,
        ip: undefined
      },
      temp: {
        id: undefined,
        ip: '',
        createTime: new Date()
      },
      dialogFormVisible: false,
      createRules: {
        ip: [
          {
            required: true,
            message: this.$t('valid.ip'),
            trigger: ['change', 'blur']
          },
          {
            min: 4,
            max: 64,
            message: this.$t('valid.ipRange'),
            trigger: ['change', 'blur']
          },
          {
            validator: validateIp,
            trigger: ['change', 'blur']
          }
        ]
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    checkPermission,
    timeStampToDate,
    getList() {
      this.listLoading = true
      selectBlackListPage(this.listQuery).then((response) => {
        this.list = response.data.blackLists
        this.total = response.data.total

        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        ip: '',
        createTime: new Date()
      }
    },
    handleFilter() {
      this.listQuery.pageNum = 1
      this.getList()
    },
    handleCreate() {
      this.resetTemp()
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleDelete(row, index) {
      MessageBox.confirm(
        this.$t('confirm.deleteBlack').toString(),
        this.$t('confirm.warn').toString(),
        {
          confirmButtonText: this.$t('confirm.yes'),
          cancelButtonText: this.$t('confirm.cancel'),
          type: 'warning'
        }
      ).then(() => {
        const tempData = Object.assign({}, row)
        deleteBlackListByIp(tempData).then(() => {
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
          createBlackList(this.temp).then(() => {
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
    }
  }
}
</script>

<style scoped>
.el-button {
  margin-left: 10px;
}
</style>
