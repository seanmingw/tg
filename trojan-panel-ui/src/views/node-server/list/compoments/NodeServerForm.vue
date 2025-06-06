<template>
  <el-dialog
    :title="textMap[dialogStatus]"
    :visible="dialogVisible"
    @close="$emit('update:dialogVisible', false)"
    width="150"
  >
    <el-form
      ref="dataForm"
      :rules="dialogStatus === 'create' ? createRules : updateRules"
      :model="nodeServer"
      label-position="left"
    >
      <el-form-item :label="$t('table.nodeServerName')" prop="name">
        <el-input v-model="nodeServer.name" clearable />
      </el-form-item>
      <el-form-item :label="$t('table.nodeServerIp')" prop="ip">
        <el-input v-model="nodeServer.ip" clearable />
      </el-form-item>
      <el-form-item :label="$t('table.nodeServerGrpcPort')" prop="grpcPort">
        <el-input-number
          v-model.number="nodeServer.grpcPort"
          controls-position="right"
          type="number"
        />
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="$emit('update:dialogVisible', false)"
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
</template>

<script>
import { createNodeServer, updateNodeServerById } from '@/api/node-server'

export default {
  name: 'NodeServerForm',
  props: {
    nodeServer: {
      type: Object,
      required: true
    },
    dialogStatus: {
      type: String,
      required: true
    },
    dialogVisible: {
      type: Boolean,
      required: true
    },
    getList: {
      type: Function,
      required: true
    }
  },
  data() {
    return {
      textMap: {
        update: this.$t('table.edit'),
        create: this.$t('table.add')
      },
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
          }
        ],
        name: [
          {
            required: true,
            message: this.$t('valid.nodeServerName'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 20,
            message: this.$t('valid.nodeServerNameRange'),
            trigger: ['change', 'blur']
          }
        ],
        grpcPort: [
          {
            required: true,
            message: this.$t('valid.nodePort'),
            trigger: ['change', 'blur']
          },
          {
            pattern:
              /^([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/,
            message: this.$t('valid.nodePortRange'),
            trigger: ['change', 'blur']
          }
        ]
      },
      updateRules: {
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
          }
        ],
        name: [
          {
            required: true,
            message: this.$t('valid.nodeServerName'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 20,
            message: this.$t('valid.nodeServerNameRange'),
            trigger: ['change', 'blur']
          }
        ],
        grpcPort: [
          {
            required: true,
            message: this.$t('valid.nodePort'),
            trigger: ['change', 'blur']
          },
          {
            pattern:
              /^([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/,
            message: this.$t('valid.nodePortRange'),
            trigger: ['change', 'blur']
          }
        ]
      }
    }
  },
  methods: {
    clearValidate() {
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          createNodeServer(this.nodeServer).then(() => {
            this.getList()
            this.$emit('update:dialogVisible', false)
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
          const tempData = Object.assign({}, this.nodeServer)
          updateNodeServerById(tempData).then(() => {
            this.getList()
            this.$emit('update:dialogVisible', false)
            this.$notify({
              title: 'Success',
              message: this.$t('confirm.modifySuccess'),
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

<style scoped></style>
