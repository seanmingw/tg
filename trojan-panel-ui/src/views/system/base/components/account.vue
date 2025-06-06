<template>
  <div class="app-container">
    <el-form
      ref="dataForm"
      :rules="updateRules"
      :model="systemConfig"
      label-position="left"
    >
      <el-form-item :label="$t('config.registerEnable')" prop="registerEnable">
        <el-switch
          v-model="systemConfig.registerEnable"
          :active-value="1"
          :inactive-value="0"
          class="ml-2"
          active-color="#13ce66"
          inactive-color="#ff4949"
        />
      </el-form-item>
      <el-form-item :label="$t('config.registerQuota')" prop="registerQuota">
        <el-input-number
          v-model.number="systemConfig.registerQuota"
          controls-position="right"
          type="number"
          :disabled="registerDisable"
        />
      </el-form-item>
      <el-form-item
        :label="$t('config.registerExpireDays')"
        prop="registerExpireDays"
      >
        <el-input-number
          v-model.number="systemConfig.registerExpireDays"
          controls-position="right"
          type="number"
          :disabled="registerDisable"
        />
      </el-form-item>
      <el-form-item
        :label="$t('config.resetDownloadAndUploadMonth')"
        prop="resetDownloadAndUploadMonth"
      >
        <el-switch
          v-model="systemConfig.resetDownloadAndUploadMonth"
          :active-value="1"
          :inactive-value="0"
          class="ml-2"
          active-color="#13ce66"
          inactive-color="#ff4949"
        />
      </el-form-item>
      <el-form-item
        :label="$t('config.trafficRankEnable')"
        prop="trafficRankEnable"
      >
        <el-switch
          v-model="systemConfig.trafficRankEnable"
          :active-value="1"
          :inactive-value="0"
          class="ml-2"
          active-color="#13ce66"
          inactive-color="#ff4949"
        />
      </el-form-item>
      <el-form-item :label="$t('config.captchaEnable')" prop="captchaEnable">
        <el-switch
          v-model="systemConfig.captchaEnable"
          :active-value="1"
          :inactive-value="0"
          class="ml-2"
          active-color="#13ce66"
          inactive-color="#ff4949"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="updateData()"
          >{{ $t('table.confirm') }}
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { updateSystemById } from '@/api/system'

export default {
  name: 'account',
  props: {
    systemConfig: {
      type: Object,
      required: true
    }
  },
  computed: {
    registerDisable: function () {
      return this.systemConfig.registerEnable === 0
    }
  },
  data() {
    return {
      updateRules: {
        registerEnable: [
          {
            required: true,
            message: this.$t('valid.registerEnable'),
            trigger: ['change', 'blur']
          }
        ],
        registerQuota: [
          {
            required: true,
            message: this.$t('valid.registerQuota'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 0,
            max: 1024000,
            message: this.$t('valid.registerQuotaRange'),
            trigger: ['change', 'blur']
          }
        ],
        registerExpireDays: [
          {
            required: true,
            message: this.$t('valid.registerExpireDays'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 0,
            max: 365,
            message: this.$t('valid.registerExpireDaysRange'),
            trigger: ['change', 'blur']
          }
        ],
        resetDownloadAndUploadMonth: [
          {
            required: true,
            message: this.$t('valid.resetDownloadAndUploadMonth'),
            trigger: ['change', 'blur']
          }
        ],
        trafficRankEnable: [
          {
            required: true,
            message: this.$t('valid.trafficRankEnable'),
            trigger: ['change', 'blur']
          }
        ]
      }
    }
  },
  methods: {
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.systemConfig)
          updateSystemById(tempData).then(() => {
            this.$nextTick(() => {
              this.$refs['dataForm'].clearValidate()
            })
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

data() {
    return {
        systemConfig: {
            // ... existing 字段 ...
            ipLimit: 0 // 新增设备数限制字段
        },
        rules: {
            // ... existing 验证规则 ...
            ipLimit: [
                { type: 'number', min: 0, max: 100, message: '设备数限制需在0-100之间', trigger: ['change', 'blur'] }
            ]
        }
    }
}
