<template>
  <div class="app-container">
    <aside v-if="emailEnable === 1">
      {{ this.$t('remark.emailEnable') }}
    </aside>
    <el-form
      ref="dataForm"
      :rules="updateRules"
      :model="temp"
      label-position="left"
    >
      <el-form-item :label="$t('table.email')" prop="email">
        <el-input
          v-model="temp.email"
          :placeholder="$t('table.email')"
          clearable
        />
      </el-form-item>
      <el-form-item :label="$t('table.username')" prop="username">
        <el-input
          v-model="temp.username"
          :placeholder="$t('table.username')"
          clearable
        />
      </el-form-item>
      <el-form-item :label="$t('table.pass')" prop="pass">
        <el-input
          v-model="temp.pass"
          type="password"
          :placeholder="$t('table.pass')"
          clearable
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
import { updateAccountProperty } from '@/api/account'
import { setting } from '@/api/system'

export default {
  name: 'ModifyProperty',
  data() {
    const validateEmailAndUsername = (rule, value, callback) => {
      if (!this.temp.email && !this.temp.username) {
        callback(new Error(this.$t('valid.emailAOrUsernameExist')))
      } else {
        callback()
      }
    }
    return {
      temp: {
        email: undefined,
        username: undefined,
        pass: ''
      },
      emailEnable: 0,
      updateRules: {
        email: [
          {
            min: 4,
            max: 64,
            message: this.$t('valid.emailRange'),
            trigger: ['change', 'blur']
          },
          {
            pattern:
              /^([A-Za-z0-9_\-\.])+\@(163.com|126.com|qq.com|gmail.com)$/,
            message: this.$t('valid.emailElement'),
            trigger: ['change', 'blur']
          }
        ],
        username: [
          {
            message: this.$t('valid.username'),
            trigger: ['change', 'blur']
          },
          {
            min: 6,
            max: 20,
            message: this.$t('valid.usernameRange'),
            trigger: ['change', 'blur']
          },
          {
            pattern: /^[A-Za-z0-9]+$/,
            message: this.$t('valid.usernameElement'),
            trigger: ['change', 'blur']
          },
          {
            validator: validateEmailAndUsername,
            trigger: ['change', 'blur']
          }
        ],
        pass: [
          {
            required: true,
            message: this.$t('valid.pass'),
            trigger: ['change', 'blur']
          },
          {
            min: 6,
            max: 20,
            message: this.$t('valid.passRange'),
            trigger: ['change', 'blur']
          },
          {
            pattern: /^[A-Za-z0-9]+$/,
            message: this.$t('valid.passElement'),
            trigger: ['change', 'blur']
          },
          {
            validator: validateEmailAndUsername,
            trigger: ['change', 'blur']
          }
        ]
      }
    }
  },
  created() {
    this.setting()
  },
  methods: {
    resetTemp() {
      this.temp = {
        email: undefined,
        username: undefined,
        pass: ''
      }
    },
    setting() {
      setting().then((response) => {
        const { data } = response
        this.emailEnable = data.emailEnable
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)
          updateAccountProperty(tempData).then(() => {
            this.resetTemp()
            this.$nextTick(() => {
              this.$refs['dataForm'].clearValidate()
            })
            this.$notify({
              title: 'Success',
              message: this.$t('confirm.updateAccountProperty'),
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
