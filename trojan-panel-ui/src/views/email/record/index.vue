<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.toEmail"
        :placeholder="$t('table.toEmail')"
        style="width: 200px"
        class="filter-item"
        clearable
        @keyup.enter.native="handleFilter"
        @clear="handleFilter"
        maxlength="64"
      />
      <el-select
        v-model="listQuery.state"
        :placeholder="$t('table.status')"
        style="width: 200px"
        class="filter-item"
        clearable
        @clear="
          () => {
            listQuery.state = undefined
          }
        "
      >
        <el-option
          :label="item.label"
          :value="item.value"
          :key="item.value"
          v-for="item in states"
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
      <el-table-column :label="$t('table.toEmail')" width="180" align="center">
        <template slot-scope="{ row }">
          <span>{{ row.toEmail }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.subject')" width="250" align="center">
        <template slot-scope="{ row }">
          <span>{{ row.subject }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.content')" width="400" align="center">
        <template slot-scope="{ row }">
          <span>{{ row.content }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.status')" width="180" align="center">
        <template slot-scope="{ row }">
          <el-tag :type="row.state | stateFilter">
            {{ stateDescFilter(row.state) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('table.createTime')"
        align="center"
        class-name="small-padding fixed-width"
      >
        <template slot-scope="{ row }">
          <span>{{ timeStampToDate(row.createTime, false) }}</span>
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
import { timeStampToDate } from '@/utils'
import { selectEmailRecordPage } from '@/api/email-record'

export default {
  name: 'Record',
  components: { Pagination },

  filters: {
    stateFilter(state) {
      let stateMap = new Map()
      stateMap.set(-1, 'danger')
      stateMap.set(0, 'warn')
      stateMap.set(1, 'success')
      return stateMap.get(state)
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
        toEmail: undefined,
        state: undefined
      },
      temp: {
        id: undefined,
        toEmail: '',
        subject: '',
        content: '',
        state: 0,
        createTime: new Date()
      },
      states: [
        { value: -1, label: this.$t('table.sendFail') },
        { value: 0, label: this.$t('table.sendWait') },
        { value: 1, label: this.$t('table.sendSuccess') }
      ],
      stateDesc: [
        { value: -1, desc: this.$t('table.sendFail') },
        { value: 0, desc: this.$t('table.sendWait') },
        { value: 1, desc: this.$t('table.sendSuccess') }
      ]
    }
  },
  created() {
    this.getList()
  },
  methods: {
    timeStampToDate,
    getList() {
      this.listLoading = true
      selectEmailRecordPage(this.listQuery).then((response) => {
        this.list = response.data.emailRecords
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
    stateDescFilter(state) {
      let stateDesc = this.stateDesc.find((item) => item.value === state)
      if (stateDesc && stateDesc.desc) {
        return stateDesc.desc
      } else {
        return ''
      }
    }
  }
}
</script>

<style scoped>
.el-button {
  margin-left: 10px;
}
</style>
