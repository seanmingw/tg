<template>
  <div class="app-container">
    <el-form
      ref="dataForm"
      :rules="updateRules"
      :model="systemConfig"
      label-position="left"
    >
      <el-form-item :label="$t('config.systemLogo')" prop="systemName">
        <upload-logo />
      </el-form-item>
      <el-form-item :label="$t('config.systemName')" prop="systemName">
        <el-input v-model="systemConfig.systemName" clearable />
      </el-form-item>
      <el-form-item :label="$t('config.clashRule')" prop="clashRule">
        <el-input
          type="textarea"
          :autosize="{ minRows: 2, maxRows: 8 }"
          v-model="systemConfig.clashRule"
          clearable
        />
      </el-form-item>
      <el-form-item :label="$t('config.xrayTemplate')" prop="xrayTemplate">
        <JsonEditorVue
          v-model="systemConfig.xrayTemplateEntity"
          v-bind="systemConfig.xrayTemplateEntity"
          mode="text"
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
import JsonEditorVue from 'json-editor-vue'
import UploadLogo from '@/components/UploadLogo'

export default {
  name: 'templateConfig',
  components: { JsonEditorVue, UploadLogo },
  props: {
    systemConfig: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      updateRules: {
        systemName: [
          {
            required: true,
            message: this.$t('valid.systemName'),
            trigger: ['change', 'blur']
          },
          {
            min: 2,
            max: 32,
            message: this.$t('valid.systemNameRange'),
            trigger: ['change', 'blur']
          }
        ],
        clashRule: [
          {
            min: 0,
            max: 102400,
            message: this.$t('valid.clashRuleRange'),
            trigger: ['change', 'blur']
          }
        ],
        xrayTemplate: [
          {
            min: 0,
            max: 10240,
            message: this.$t('valid.xrayTemplateRange'),
            trigger: ['change', 'blur']
          }
        ]
      }
    }
  },
  methods: {
    updateData() {
      if (typeof this.systemConfig.xrayTemplateEntity !== 'object')
        this.systemConfig.xrayTemplateEntity = JSON.parse(
          this.systemConfig.xrayTemplateEntity
        )
      this.systemConfig.xrayTemplate = JSON.stringify(
        this.systemConfig.xrayTemplateEntity
      )
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
