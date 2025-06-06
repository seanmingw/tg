<template>
  <div class="login-container">
    <el-form
      ref="registerForm"
      :model="registerForm"
      :rules="registerRules"
      class="login-form"
      auto-complete="on"
      label-position="left"
    >
      <div class="title-container">
        <h3 class="title">{{ $t('register.title') }}</h3>
      </div>

      <el-form-item prop="username">
        <span class="svg-container">
          <svg-icon icon-class="username" />
        </span>
        <el-input
          ref="username"
          v-model="registerForm.username"
          :placeholder="$t('register.username')"
          name="username"
          type="text"
          tabindex="1"
          auto-complete="on"
          clearable
        />
      </el-form-item>

      <el-form-item prop="pass">
        <span class="svg-container">
          <svg-icon icon-class="pass" />
        </span>
        <el-input
          ref="pass"
          v-model="registerForm.passOne"
          :type="passwordType"
          :placeholder="$t('register.passOne')"
          name="password"
          tabindex="2"
          auto-complete="on"
          clearable
        />
        <span class="show-pwd" @click="showPwd">
          <svg-icon
            :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'"
          />
        </span>
      </el-form-item>

      <el-form-item prop="pass">
        <span class="svg-container">
          <svg-icon icon-class="pass" />
        </span>
        <el-input
          ref="pass"
          v-model="registerForm.pass"
          :type="passwordType"
          :placeholder="$t('register.pass')"
          name="password"
          tabindex="3"
          auto-complete="on"
          clearable
          @keyup.enter.native="handleRegister"
        />
        <span class="show-pwd" @click="showPwd">
          <svg-icon
            :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'"
          />
        </span>
      </el-form-item>

      <!-- 验证码 -->
      <el-form-item prop="captchaCode" v-if="captchaEnable">
        <span class="svg-container">
          <svg-icon icon-class="valid-code" />
        </span>
        <el-input
          v-model="registerForm.captchaCode"
          auto-complete="off"
          :placeholder="$t('login.code')"
          style="width: 65%"
          @keyup.enter="handleRegister"
        />

        <div class="captcha">
          <img :src="captchaImg" @click="handleCaptchaGenerate" height="38px" />
        </div>
      </el-form-item>

      <el-button
        :loading="loading"
        type="primary"
        style="width: 100%; margin-bottom: 30px"
        @click.native.prevent="handleRegister"
        >{{ $t('register.register') }}
      </el-button>
      <el-button type="primary" style="width: 100%; margin: 0">
        <router-link to="/login" custom v-slot="{ navigate }">
          <span @click="navigate" role="link">{{ $t('register.login') }}</span>
        </router-link>
      </el-button>
    </el-form>
  </div>
</template>

<script>
import { generateCaptcha } from '@/api/account'
import { setting } from '@/api/system'

export default {
  name: 'index',
  data() {
    const validatePass = (rule, value, callback) => {
      if (this.registerForm.passOne !== this.registerForm.pass) {
        callback(new Error(this.$t('valid.passNotSame')))
      } else {
        callback()
      }
    }
    const validateUsername = (rule, value, callback) => {
      if (this.registerForm.username.trim().indexOf('admin') >= 0) {
        callback(new Error(this.$t('valid.usernameNotExistAdmin')))
      } else {
        callback()
      }
    }
    return {
      registerForm: {
        username: '',
        passOne: '',
        pass: '',
        captchaId: '',
        captchaCode: ''
      },
      captchaImg: '',
      captchaEnable: 0,
      registerRules: {
        username: [
          {
            required: true,
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
            validator: validateUsername,
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
            validator: validatePass,
            trigger: ['change', 'blur']
          }
        ],
        captchaCode: [
          {
            required: true,
            message: this.$t('valid.code'),
            trigger: ['change', 'blur']
          }
        ]
      },
      loading: false,
      passwordType: 'password'
    }
  },
  created() {
    this.handleCaptchaGenerate()
    this.setting()
  },
  methods: {
    handleCaptchaGenerate() {
      generateCaptcha().then((response) => {
        this.registerForm.captchaId = response.data.captchaId
        this.captchaImg = response.data.captchaImg
      })
    },
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
      this.$nextTick(() => {
        this.$refs.pass.focus()
      })
    },
    handleRegister() {
      this.$refs.registerForm.validate((valid) => {
        if (valid) {
          this.loading = true
          this.$store
            .dispatch('account/register', this.registerForm)
            .then(() => {
              this.$router.push({ path: '/login' })
              this.loading = false
            })
            .catch(() => {
              this.loading = false
              this.handleCaptchaGenerate()
            })
        } else {
          // console.log('error submit!!')
          return false
        }
      })
    },
    setting() {
      setting().then((response) => {
        const { data } = response
        this.captchaEnable = data.captchaEnable
      })
    }
  }
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg: #283443;
$light_gray: #fff;
$cursor: #fff;

@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
  .login-container .el-input input {
    color: $cursor;
  }
}

/* reset element-ui css */
.login-container {
  .el-input {
    display: inline-block;
    height: 47px;
    width: 85%;

    input {
      background: transparent;
      border: 0px;
      -webkit-appearance: none;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: $light_gray;
      height: 47px;
      caret-color: $cursor;

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px $bg inset !important;
        -webkit-text-fill-color: $cursor !important;
      }
    }
  }

  .el-form-item {
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.1);
    border-radius: 5px;
    color: #454545;
  }
}
</style>

<style lang="scss" scoped>
$bg: #2d3a4b;
$dark_gray: #889aa4;
$light_gray: #eee;

.login-container {
  min-height: 100%;
  width: 100%;
  background-color: $bg;
  overflow: hidden;

  .login-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 160px 35px 0;
    margin: 0 auto;
    overflow: hidden;
  }

  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: $light_gray;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }

  .captcha {
    position: absolute;
    right: 0;
    top: 0;

    img {
      height: 48px;
      cursor: pointer;
      vertical-align: middle;
    }
  }
}
</style>
