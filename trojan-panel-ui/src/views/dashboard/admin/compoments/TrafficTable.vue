<template>
  <el-table :data="list" style="width: 100%">
    <el-table-column
      :label="$t('dashboard.ranking')"
      width="100"
      align="center"
      type="index"
    />
    <el-table-column
      :label="$t('dashboard.username')"
      width="200"
      align="center"
    >
      <template slot-scope="scope">
        {{ scope.row.username }}
      </template>
    </el-table-column>
    <el-table-column
      :label="$t('dashboard.trafficUsed')"
      min-width="200"
      align="center"
    >
      <template slot-scope="scope">
        {{ getFlow(scope.row.trafficUsed) }}</template
      >
    </el-table-column>
  </el-table>
</template>

<script>
import { trafficRank } from '@/api/dashboard'
import { getFlow } from '@/utils/account'

export default {
  name: 'trafficTable',
  data() {
    return {
      list: null
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    getFlow,
    fetchData() {
      trafficRank().then((response) => {
        this.list = response.data
      })
    }
  }
}
</script>
