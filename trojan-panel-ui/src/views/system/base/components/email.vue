<template>
  <div class="app-container">
    <el-form
      ref="dataForm"
      :rules="updateRules"
      :model="systemConfig"
      label-position="left"
    >
      <el-form-item
        :label="$t('config.emailEnable')"
        prop="emailEnable"
        clearable
      >
        <el-switch
          v-model="systemConfig.emailEnable"
          :active-value="1"
          :inactive-value="0"
          class="ml-2"
          active-color="#13ce66"
          inactive-color="#ff4949"
        />
      </el-form-item>

      <el-form-item :label="$t('config.emailHost')" prop="emailHost">
        <el-input
          :disabled="emailDisable"
          v-model="systemConfig.emailHost"
          clearable
        />
      </el-form-item>

      <el-form-item :label="$t('config.emailPort')" prop="emailPort">
        <el-input-number
          v-model.number="systemConfig.emailPort"
          controls-position="right"
          type="number"
          :disabled="emailDisable"
        />
      </el-form-item>

      <el-form-item
        :label="$t('config.emailUsername')"
        prop="emailUsername"
        clearable
      >
        <el-input
          :disabled="emailDisable"
          v-model="systemConfig.emailUsername"
        />
      </el-form-item>

      <el-form-item
        :label="$t('config.emailPassword')"
        prop="emailPassword"
        clearable
      >
        <el-input
          :disabled="emailDisable"
          type="password"
          v-model="systemConfig.emailPassword"
        />
      </el-form-item>

      <el-form-item
        :label="$t('config.expireWarnEnable')"
        prop="expireWarnEnable"
      >
        <el-switch
          v-model="systemConfig.expireWarnEnable"
          :active-value="1"
          :inactive-value="0"
          class="ml-2"
          active-color="#13ce66"
          inactive-color="#ff4949"
          :disabled="emailDisable"
        />
      </el-form-item>

      <el-form-item :label="$t('config.expireWarnDay')" prop="expireWarnDay">
        <el-input-number
          v-model.number="systemConfig.expireWarnDay"
          controls-position="right"
          type="number"
          :disabled="emailDisable || expireWarnDisable"
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
  name: 'email',
  props: {
    systemConfig: {
      type: Object,
      required: true
    }
  },
  computed: {
    emailDisable: function () {
      return this.systemConfig.emailEnable === 0
    },
    expireWarnDisable: function () {
      return this.systemConfig.expireWarnEnable === 0
    }
  },
  data() {
    return {
      updateRules: {
        emailEnable: [
          {
            required: true,
            message: this.$t('valid.emailEnable'),
            trigger: ['change', 'blur']
          }
        ],
        emailHost: [
          {
            min: 4,
            max: 64,
            message: this.$t('valid.emailHostRange'),
            trigger: ['change', 'blur']
          }
        ],
        emailPort: [
          {
            required: true,
            message: this.$t('valid.emailPort'),
            trigger: ['change', 'blur']
          },
          {
            pattern:
              /^([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/,
            message: this.$t('valid.emailPortRange'),
            trigger: ['change', 'blur']
          }
        ],
        emailUsername: [
          {
            min: 3,
            max: 32,
            message: this.$t('valid.emailUsernameRange'),
            trigger: ['change', 'blur']
          }
        ],
        emailPassword: [
          {
            min: 3,
            max: 32,
            message: this.$t('valid.emailPasswordRange'),
            trigger: ['change', 'blur']
          }
        ],
        expireWarnEnable: [
          {
            required: true,
            message: this.$t('valid.expireWarnEnable'),
            trigger: ['change', 'blur']
          }
        ],
        expireWarnDay: [
          {
            required: true,
            message: this.$t('valid.expireWarnDay'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 0,
            max: 365,
            message: this.$t('valid.expireWarnDayRange'),
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
