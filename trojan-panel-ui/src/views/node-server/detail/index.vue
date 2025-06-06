<template>
  <div class="dashboard-editor-container">
    <NodeServerGroup :node-server-group-data="nodeServerGroupData" />
  </div>
</template>

<script>
import NodeServerGroup from '@/views/node-server/detail/compoments/NodeServerGroup'
import { nodeServerState } from '@/api/node-server'
import Cookies from 'js-cookie'

export default {
  name: 'index',
  components: { NodeServerGroup },
  data() {
    return {
      nodeServerGroupData: {
        cpuUsed: 0,
        memUsed: 0,
        diskUsed: 0
      }
    }
  },
  created() {
    nodeServerState({ id: Cookies.get('nodeServerId') }).then((response) => {
      const { data } = response
      this.nodeServerGroupData = data
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
