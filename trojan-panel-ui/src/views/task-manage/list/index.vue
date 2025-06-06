<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.accountUsername"
        :placeholder="$t('table.fileTaskAccountUsername')"
        style="width: 200px"
        class="filter-item"
        clearable
        @keyup.enter.native="handleFilter"
        @clear="handleFilter"
        maxlength="20"
      />
      <el-button
        class="filter-item"
        type="primary"
        icon="el-icon-search"
        @click="handleFilter"
      >
        {{ $t('table.search') }}
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
        :label="$t('table.fileTaskName')"
        width="180"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.fileTaskType')"
        width="100"
        align="center"
      >
        <template slot-scope="{ row }">
          <el-tag type="success">
            <span>{{ typeComputed(row.type) }}</span>
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.fileTaskStatus')"
        width="100"
        align="center"
      >
        <template slot-scope="{ row }">
          <el-tag :type="row.status | statusTypeFilter">
            <span>{{ statusComputed(row.status) }}</span>
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.fileTaskErrMsg')"
        width="180"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ row.status === -1 ? row.errMsg : '' }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.fileTaskAccountUsername')"
        width="180"
        align="center"
      >
        <template slot-scope="{ row }">
          <span>{{ row.accountUsername }}</span>
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
            :disabled="row.status !== 2 || row.type === 3 || row.type === 4"
            type="primary"
            size="mini"
            @click="handleDownload(row, $index)"
            v-if="checkPermission(['sysadmin', 'admin'])"
          >
            {{ $t('table.fileTaskDownload') }}
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
  </div>
</template>

<script>
import Pagination from '@/components/Pagination'
import checkPermission from '@/utils/permission'
import { timeStampToDate } from '@/utils'
import { downloadFileTask, selectFileTaskPage } from '@/api/file-task'

export default {
  name: 'TaskList',
  components: { Pagination },
  data() {
    return {
      tableKey: 0,
      listLoading: true,
      list: null,
      total: 0,
      listQuery: {
        pageNum: 1,
        pageSize: 20,
        accountUsername: undefined
      },
      temp: {
        id: undefined,
        name: '',
        path: '',
        type: 1,
        status: 0,
        errMsg: '',
        accountUsername: '',
        createTime: new Date()
      }
    }
  },
  created() {
    this.getList()
  },
  computed: {
    typeComputed() {
      return function (type) {
        if (type === 1) {
          return this.$t('table.taskTypeAccountExport')
        } else if (type === 2) {
          return this.$t('table.taskTypeNodeServerExport')
        } else if (type === 3) {
          return this.$t('table.taskTypeAccountImport')
        } else if (type === 4) {
          return this.$t('table.taskTypeNodeServerImport')
        }
        return ''
      }
    },
    statusComputed() {
      return function (status) {
        if (status === -1) {
          return this.$t('table.taskFail')
        } else if (status === 0) {
          return this.$t('table.taskWait')
        } else if (status === 1) {
          return this.$t('table.taskDoing')
        } else if (status === 2) {
          return this.$t('table.taskSuccess')
        }
        return ''
      }
    }
  },
  filters: {
    statusTypeFilter(status) {
      if (status === -1) {
        return 'danger'
      } else if (status === 0) {
        return 'warning'
      } else if (status === 1) {
        return 'info'
      } else if (status === 2) {
        return 'success'
      }
      return ''
    }
  },
  methods: {
    checkPermission,
    timeStampToDate,
    getList() {
      this.listLoading = true
      selectFileTaskPage(this.listQuery).then((response) => {
        this.list = response.data.fileTasks
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
    handleDownload(row) {
      const tempData = Object.assign({}, row)
      downloadFileTask(tempData).then((res) => {
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
    }
  }
}
</script>

<style scoped>
.el-button {
  margin-left: 10px;
}
</style>
