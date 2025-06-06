<template>
  <div class="navbar">
    <hamburger
      :is-active="sidebar.opened"
      class="hamburger-container"
      @toggleClick="toggleSideBar"
    />

    <breadcrumb class="breadcrumb-container" />

    <div class="right-menu">
      <template v-if="device !== 'mobile'">
        <screenfull id="screenfull" class="right-menu-item hover-effect" />

        <el-tooltip
          :content="$t('navbar.size')"
          effect="dark"
          placement="bottom"
        >
          <SizeSelect id="size-select" class="right-menu-item hover-effect" />
        </el-tooltip>

        <lang-select class="right-menu-item hover-effect" />
      </template>
      <el-dropdown class="avatar-container" trigger="click">
        <div class="avatar-wrapper">
          <img src="/api/image/logo" class="user-avatar" />
          <i class="el-icon-caret-bottom" />
        </div>
        <el-dropdown-menu slot="dropdown" class="user-dropdown">
          <router-link to="/modify" custom v-slot="{ navigate }">
            <el-dropdown-item>
              <span @click="navigate" role="link">
                {{ $t('navbar.profile') }}
              </span>
            </el-dropdown-item>
          </router-link>
          <a
            target="_blank"
            href="https://github.com/trojanpanel"
            v-if="checkPermission(['sysadmin', 'admin'])"
          >
            <el-dropdown-item>{{ $t('navbar.github') }}</el-dropdown-item>
          </a>
          <a
            target="_blank"
            href="https://trojanpanel.github.io"
            v-if="checkPermission(['sysadmin', 'admin'])"
          >
            <el-dropdown-item>{{ $t('navbar.doc') }}</el-dropdown-item>
          </a>
          <el-dropdown-item divided @click.native="logout">
            <span style="display: block">{{ $t('navbar.logout') }}</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'
import Screenfull from '@/components/Screenfull'
import SizeSelect from '@/components/SizeSelect'
import LangSelect from '@/components/LangSelect'
import checkPermission from '@/utils/permission' // 权限判断指令

export default {
  components: {
    Breadcrumb,
    Hamburger,
    Screenfull,
    SizeSelect,
    LangSelect
  },
  computed: {
    ...mapGetters(['sidebar', 'avatar', 'device'])
  },
  methods: {
    checkPermission,
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    async logout() {
      await this.$store.dispatch('account/logout')
      await this.$router.push(`/login`)
    }
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background 0.3s;
    -webkit-tap-highlight-color: transparent;

    &:hover {
      background: rgba(0, 0, 0, 0.025);
    }
  }

  .breadcrumb-container {
    float: left;
  }

  .errLog-container {
    display: inline-block;
    vertical-align: top;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background 0.3s;

        &:hover {
          background: rgba(0, 0, 0, 0.025);
        }
      }
    }

    .avatar-container {
      margin-right: 30px;

      .avatar-wrapper {
        margin-top: 5px;
        position: relative;

        .user-avatar {
          cursor: pointer;
          width: 40px;
          height: 40px;
          border-radius: 10px;
        }

        .el-icon-caret-bottom {
          cursor: pointer;
          position: absolute;
          right: -20px;
          top: 25px;
          font-size: 12px;
        }
      }
    }
  }
}
</style>
