<template>
  <div class="dashboard-editor-container">
    <panel-group :group-data="panelGroupData" />
    <el-row :gutter="8">
      <el-col
        :xs="{ span: 24 }"
        :sm="{ span: 24 }"
        :md="{ span: 24 }"
        :lg="{ span: 12 }"
        :xl="{ span: 12 }"
      >
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>{{ $t('dashboard.trafficRank') }}</span>
          </div>
          <div class="component-item">
            <traffic-table />
          </div>
        </el-card>
      </el-col>
      <el-col
        :xs="{ span: 24 }"
        :sm="{ span: 12 }"
        :md="{ span: 12 }"
        :lg="{ span: 6 }"
        :xl="{ span: 6 }"
      >
      </el-col>
      <el-col
        :xs="{ span: 24 }"
        :sm="{ span: 12 }"
        :md="{ span: 12 }"
        :lg="{ span: 6 }"
        :xl="{ span: 6 }"
      >
      </el-col>
    </el-row>
  </div>
</template>

<script>
import PanelGroup from '@/views/dashboard/admin/compoments/PanelGroup.vue'
import TrafficTable from '@/views/dashboard/admin/compoments/TrafficTable'
import { panelGroup } from '@/api/dashboard'

export default {
  name: 'Admin',
  components: {
    PanelGroup,
    TrafficTable
  },
  data() {
    return {
      panelGroupData: {
        totalFlow: 0,
        residualFlow: 0,
        nodeNum: 0,
        expireTime: new Date(),
        accountCount: 0,
        cpuUsed: 0,
        memUsed: 0,
        diskUsed: 0
      }
    }
  },
  created() {
    panelGroup().then((response) => {
      const { data } = response
      this.panelGroupData = data
    })
  }
}
</script>

<style lang="scss" scoped>
.dashboard-editor-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  position: relative;

  .github-corner {
    position: absolute;
    top: 0px;
    border: 0;
    right: 0;
  }

  .chart-wrapper {
    background: #fff;
    padding: 16px 16px 0;
    margin-bottom: 32px;
  }
}

@media (max-width: 1024px) {
  .chart-wrapper {
    padding: 8px;
  }
}
</style>
