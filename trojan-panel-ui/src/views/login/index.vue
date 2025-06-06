<template>
  <div class="login-container">
    <el-form
      ref="loginForm"
      :model="loginForm"
      :rules="loginRules"
      class="login-form"
      auto-complete="on"
      label-position="left"
    >
      <div class="title-container">
        <h3 class="title">
          {{ systemName }}
        </h3>
        <lang-select class="set-language" />
      </div>

      <el-form-item prop="username">
        <span class="svg-container">
          <svg-icon icon-class="username" />
        </span>
        <el-input
          ref="username"
          v-model="loginForm.username"
          :placeholder="$t('login.username')"
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
          :key="passwordType"
          ref="pass"
          v-model="loginForm.pass"
          :type="passwordType"
          :placeholder="$t('login.password')"
          name="password"
          tabindex="2"
          auto-complete="on"
          @keyup.enter.native="handleLogin"
          clearable
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
          v-model="loginForm.captchaCode"
          auto-complete="off"
          :placeholder="$t('login.code')"
          style="width: 65%"
          @keyup.enter.native="handleLogin"
        />

        <div class="captcha">
          <img
            alt="captcha"
            :src="captchaImg"
            @click="handleCaptchaGenerate"
            height="38px"
          />
        </div>
      </el-form-item>

      <el-button
        :loading="loading"
        type="primary"
        style="width: 100%; margin-bottom: 30px"
        @click.native.prevent="handleLogin"
        >{{ $t('login.logIn') }}
      </el-button>

      <el-button
        type="primary"
        style="width: 100%; margin: 0"
        v-if="registerEnable === 1"
      >
        <router-link to="/register" custom v-slot="{ navigate }">
          <span @click="navigate" role="link">{{ $t('login.register') }}</span>
        </router-link>
      </el-button>
    </el-form>
  </div>
</template>

<script>
import LangSelect from '@/components/LangSelect'
import { setting } from '@/api/system'
import { generateCaptcha } from '@/api/account'

export default {
  name: 'index',
  components: { LangSelect },
  data() {
    return {
      loginForm: {
        username: '',
        pass: '',
        captchaId: '',
        captchaCode: ''
      },
      captchaImg: '',
      loginRules: {
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
      passwordType: 'password',
      redirect: undefined,
      registerEnable: 0,
      systemName: '',
      captchaEnable: 0
    }
  },
  watch: {
    $route: {
      handler: function (route) {
        this.redirect = route.query && route.query.redirect
      },
      immediate: true
    }
  },
  created() {
    this.setting()
    this.handleCaptchaGenerate()
  },
  methods: {
    handleCaptchaGenerate() {
      generateCaptcha().then((response) => {
        this.loginForm.captchaId = response.data.captchaId
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
    handleLogin() {
      this.$refs.loginForm.validate((valid) => {
        if (valid) {
          this.loading = true
          this.$store
            .dispatch('account/login', this.loginForm)
            .then(() => {
              this.$router
                .push({ path: this.redirect || '/' })
                .catch(() => true)
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
        this.registerEnable = data.registerEnable
        this.systemName = data.systemName
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

    .set-language {
      color: #fff;
      position: absolute;
      top: 3px;
      font-size: 18px;
      right: 0px;
      cursor: pointer;
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

  .thirdparty-button {
    position: absolute;
    right: 0;
    bottom: 6px;
  }

  @media only screen and (max-width: 470px) {
    .thirdparty-button {
      display: none;
    }
  }
}
</style>
