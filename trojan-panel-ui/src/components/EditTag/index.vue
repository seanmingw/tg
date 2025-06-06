<template>
  <div>
    <el-tag
      :key="tag"
      v-for="tag in dynamicTagsProps"
      effect="dark"
      type="success"
      size="medium"
      closable
      :disable-transitions="false"
      @close="handleClose(tag)"
    >
      {{ tag }}
    </el-tag>
    <el-input
      class="input-new-tag"
      v-if="inputVisible"
      v-model="inputValue"
      ref="saveTagInput"
      size="small"
      @keyup.enter.native="handleInputConfirm"
      @blur="handleInputConfirm"
    >
    </el-input>
    <el-button
      icon="el-icon-plus"
      v-else
      class="button-new-tag"
      size="mini"
      type="primary"
      @click="showInput"
    >
    </el-button>
  </div>
</template>

<script>
import { Message } from 'element-ui'

export default {
  data() {
    return {
      inputVisible: false,
      inputValue: ''
    }
  },
  props: {
    dynamicTagsProps: {
      type: Array,
      require: true
    },
    valueCanEmpty: {
      type: Boolean,
      require: false,
      default: false
    }
  },
  methods: {
    handleClose(tag) {
      this.dynamicTagsProps.splice(this.dynamicTagsProps.indexOf(tag), 1)
    },

    showInput() {
      this.inputVisible = true
      this.$nextTick((_) => {
        this.$refs.saveTagInput.$refs.input.focus()
      })
    },

    handleInputConfirm() {
      let inputValue = this.inputValue
      if ((this.valueCanEmpty && inputValue === '') || inputValue) {
        if (this.dynamicTagsProps.indexOf(inputValue) === -1) {
          this.dynamicTagsProps.push(inputValue)
        } else {
          Message({
            message: this.$t('confirm.alreadyExists').toString(),
            type: 'error',
            duration: 5 * 1000
          })
        }
      }
      this.inputVisible = false
      this.inputValue = ''
    }
  }
}
</script>

<style scoped>
.el-tag + .el-tag {
  margin-left: 10px;
}

.button-new-tag {
  margin-left: 10px;
  height: 32px;
  line-height: 30px;
  padding-top: 0;
  padding-bottom: 0;
}

.input-new-tag {
  width: 120px;
  margin-left: 10px;
  vertical-align: bottom;
}
</style>
