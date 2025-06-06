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
        <el-card class="box-card" v-if="trafficRankEnable === 1">
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
import PanelGroup from '@/views/dashboard/user/compoments/PanelGroup'
import { panelGroup } from '@/api/dashboard'
import TrafficTable from '@/views/dashboard/admin/compoments/TrafficTable'
import { setting } from '@/api/system'

export default {
  name: 'User',
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
        expireTime: new Date()
      },
      trafficRankEnable: 0
    }
  },
  created() {
    panelGroup().then((response) => {
      const { data } = response
      this.panelGroupData = data
    })
    this.setting()
  },
  methods: {
    setting() {
      setting().then((response) => {
        const { data } = response
        this.trafficRankEnable = data.trafficRankEnable
      })
    }
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
