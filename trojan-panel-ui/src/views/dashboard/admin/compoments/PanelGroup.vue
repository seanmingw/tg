<template>
  <el-row :gutter="40" class="panel-group">
    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col">
      <div class="card-panel">
        <div class="card-panel-icon-wrapper icon-sysinfo">
          <svg-icon icon-class="sysinfo" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">{{ $t('dashboard.cpuUsed') }}</div>
          <h3 :style="groupData.cpuUsed | useWarning">
            {{ groupData.cpuUsed | toPercent }}
          </h3>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col">
      <div class="card-panel">
        <div class="card-panel-icon-wrapper icon-sysinfo">
          <svg-icon icon-class="sysinfo" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">{{ $t('dashboard.memUsed') }}</div>
          <h3 :style="groupData.memUsed | useWarning">
            {{ groupData.memUsed | toPercent }}
          </h3>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col">
      <div class="card-panel">
        <div class="card-panel-icon-wrapper icon-sysinfo">
          <svg-icon icon-class="sysinfo" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">{{ $t('dashboard.diskUsed') }}</div>
          <h3 :style="groupData.diskUsed | useWarning">
            {{ groupData.diskUsed | toPercent }}
          </h3>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col">
      <div class="card-panel">
        <div class="card-panel-icon-wrapper icon-node">
          <svg-icon icon-class="node" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">{{ $t('dashboard.nodeCount') }}</div>
          <h3>{{ groupData.nodeCount }}</h3>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col">
      <div class="card-panel">
        <div class="card-panel-icon-wrapper icon-user">
          <svg-icon icon-class="account" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">{{ $t('dashboard.accountCount') }}</div>
          <h3>{{ groupData.accountCount }}</h3>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col">
      <div class="card-panel">
        <div class="card-panel-icon-wrapper icon-flow">
          <svg-icon icon-class="flow" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">{{ $t('dashboard.quota') }}</div>
          <h3>{{ quota }}</h3>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col">
      <div class="card-panel">
        <div class="card-panel-icon-wrapper icon-flow">
          <svg-icon icon-class="flow" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">{{ $t('dashboard.residualFlow') }}</div>
          <h3 :style="residualFlowStyle">{{ residualFlow }}</h3>
        </div>
      </div>
    </el-col>
    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col">
      <div class="card-panel">
        <div class="card-panel-icon-wrapper icon-time">
          <svg-icon icon-class="time" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">{{ $t('dashboard.expireTime') }}</div>
          <h3 :style="expireTimeStyle">
            {{ timeStampToDate(groupData.expireTime, false) }}
          </h3>
        </div>
      </div>
    </el-col>
  </el-row>
</template>

<script>
import { timeStampToDate } from '@/utils'
import { getFlow } from '@/utils/account'

export default {
  name: 'PanelGroup',
  props: {
    groupData: {
      type: Object,
      required: true
    }
  },
  methods: {
    getFlow,
    timeStampToDate
  },
  filters: {
    useWarning: function (value) {
      return value >= 80 ? 'color: #FF0000;' : ''
    },
    toPercent: function (value) {
      return value + '%'
    }
  },
  computed: {
    quota: function () {
      return this.groupData.quota < 0
        ? this.$t('dashboard.unlimited')
        : getFlow(this.groupData.quota)
    },
    residualFlowStyle: function () {
      return this.groupData.quota < 0
        ? this.$t('dashboard.unlimited')
        : getFlow(this.groupData.residualFlow)
    },
    residualFlow: function () {
      return this.groupData.quota < 0
        ? this.$t('dashboard.unlimited')
        : getFlow(this.groupData.residualFlow)
    },
    expireTimeStyle: function () {
      return this.groupData.expireTime <= new Date().getTime()
        ? 'color: #FF0000;'
        : ''
    }
  }
}
</script>

<style lang="scss" scoped>
.panel-group {
  margin-top: 18px;

  .card-panel-col {
    margin-bottom: 32px;
  }

  .card-panel {
    height: 108px;
    cursor: pointer;
    font-size: 12px;
    position: relative;
    overflow: hidden;
    color: #666;
    background: #fff;
    box-shadow: 4px 4px 40px rgba(0, 0, 0, 0.05);
    border-color: rgba(0, 0, 0, 0.05);

    &:hover {
      .card-panel-icon-wrapper {
        color: #fff;
      }

      .icon-user {
        background: #40c9c6;
      }

      .icon-flow {
        background: #36a3f7;
      }

      .icon-node {
        background: #f4516c;
      }

      .icon-time {
        background: #ffb6c1;
      }

      .icon-sysinfo {
        background: #d4237a;
      }
    }

    .icon-user {
      color: #40c9c6;
    }

    .icon-flow {
      color: #36a3f7;
    }

    .icon-node {
      color: #f4516c;
    }

    .icon-time {
      color: #ffb6c1;
    }

    .icon-sysinfo {
      color: #d4237a;
    }

    .card-panel-icon-wrapper {
      float: left;
      margin: 14px 0 0 14px;
      padding: 16px;
      transition: all 0.38s ease-out;
      border-radius: 6px;
    }

    .card-panel-icon {
      float: left;
      font-size: 48px;
    }

    .card-panel-description {
      float: right;
      font-weight: bold;
      margin: 26px;
      margin-left: 0px;

      .card-panel-text {
        line-height: 18px;
        color: rgba(0, 0, 0, 0.45);
        font-size: 16px;
        margin-bottom: 12px;
      }

      .card-panel-num {
        font-size: 20px;
      }
    }
  }
}

@media (max-width: 550px) {
  .card-panel-description {
    display: none;
  }

  .card-panel-icon-wrapper {
    float: none !important;
    width: 100%;
    height: 100%;
    margin: 0 !important;

    .svg-icon {
      display: block;
      margin: 14px auto !important;
      float: none !important;
    }
  }
}
</style>
