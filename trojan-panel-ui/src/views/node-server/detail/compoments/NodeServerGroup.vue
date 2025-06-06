<template>
  <el-row :gutter="40" class="panel-group">
    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col">
      <div class="card-panel">
        <div class="card-panel-icon-wrapper icon-sysinfo">
          <svg-icon icon-class="sysinfo" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">{{ $t('dashboard.cpuUsed') }}</div>
          <h3 :style="nodeServerGroupData.cpuUsed | useWarning">
            {{ nodeServerGroupData.cpuUsed | toPercent }}
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
          <h3 :style="nodeServerGroupData.memUsed | useWarning">
            {{ nodeServerGroupData.memUsed | toPercent }}
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
          <h3 :style="nodeServerGroupData.diskUsed | useWarning">
            {{ nodeServerGroupData.diskUsed | toPercent }}
          </h3>
        </div>
      </div>
    </el-col>
  </el-row>
</template>

<script>
export default {
  name: 'NodeServerGroup',
  props: {
    nodeServerGroupData: {
      type: Object,
      required: true
    }
  },
  filters: {
    useWarning: function (value) {
      return value >= 80 ? 'color: #FF0000;' : ''
    },
    toPercent: function (value) {
      return value + '%'
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
