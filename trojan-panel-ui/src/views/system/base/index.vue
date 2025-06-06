<template>
  <div class="app-container">
    <el-tabs v-model="activeName">
      <el-tab-pane :label="$t('config.accountConfig')" name="account-config">
        <account :system-config="systemConfig" />
      </el-tab-pane>
      <el-tab-pane :label="$t('config.configEmail')" name="config-email">
        <email :system-config="systemConfig" />
      </el-tab-pane>
      <el-tab-pane :label="$t('config.configWebFile')" name="config-web-file">
        <web-file />
      </el-tab-pane>
      <el-tab-pane
        :label="$t('config.templateConfig')"
        name="config-template-config"
      >
        <template-config :system-config="systemConfig" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import Account from './components/account'
import WebFile from './components/web-file'
import Email from './components/email'
import TemplateConfig from './components/template-config'
import { selectSystemByName } from '@/api/system'

export default {
  name: 'index',
  components: { Account, WebFile, Email, TemplateConfig },
  data() {
    return {
      activeName: 'account-config',
      systemConfig: {
        emailEnable: 0,
        emailHost: undefined,
        emailPassword: undefined,
        emailPort: 0,
        emailUsername: undefined,
        expireWarnDay: 0,
        expireWarnEnable: 0,
        id: 1,
        registerEnable: 1,
        registerExpireDays: 0,
        registerQuota: 0,
        resetDownloadAndUploadMonth: 0,
        trafficRankEnable: 1,
        captchaEnable: 0,
        systemName: '',
        clashRule: '',
        xrayTemplate: '',
        xrayTemplateEntity: {}
      }
    }
  },
  created() {
    this.selectDate()
  },
  methods: {
    selectDate() {
      selectSystemByName().then((response) => {
        this.systemConfig = response.data
        this.systemConfig.xrayTemplateEntity = JSON.parse(
          this.systemConfig.xrayTemplate
        )
      })
    }
  }
}
</script>

<style scoped></style>
