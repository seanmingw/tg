<template>
  <el-dialog
    :title="$t('table.import')"
    :visible="dialogFormVisible"
    @close="$emit('update:dialogFormVisible', false)"
  >
    <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left">
      <el-form-item :label="label" prop="cover">
        <el-switch
          v-model="temp.cover"
          active-color="#13ce66"
          inactive-color="#ff4949"
          :active-text="$t('table.yes')"
          :inactive-text="$t('table.no')"
          :active-value="1"
          :inactive-value="0"
        >
        </el-switch>
      </el-form-item>
      <el-form-item :label="$t('config.webFileSelect')" prop="file">
        <el-upload
          class="upload-demo"
          ref="upload"
          action=""
          :file-list="fileList"
          :http-request="importData"
          :auto-upload="false"
          accept=".json"
          :on-change="handleChange"
          :before-upload="beforeUpload"
          :limit="1"
        >
          <el-button slot="trigger" size="small" type="primary"
            >{{ $t('config.webFileSelect') }}
          </el-button>
          <el-button
            class="filter-item"
            type="success"
            icon="el-icon-download"
            @click="downloadTemplate"
          >
            {{ $t('table.downloadTemplate') }}
          </el-button>
          <div slot="tip" class="el-upload__tip">
            {{ $t('config.jsonFileTip') }}
          </div>
        </el-upload>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="$emit('update:dialogFormVisible', false)">
        {{ $t('table.cancel') }}
      </el-button>
      <el-button
        type="primary"
        @click="submitImport"
        :disabled="fileList.length === 0"
      >
        {{ $t('table.confirm') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script>
import { Message } from 'element-ui'

export default {
  name: 'index',
  props: {
    dialogFormVisible: {
      type: Boolean,
      required: true
    },
    label: {
      type: String,
      required: true
    },
    importData: {
      type: Function,
      required: true
    },
    downloadTemplate: {
      type: Function,
      required: true
    }
  },
  data() {
    return {
      temp: {
        cover: 0
      },
      fileList: [],
      rules: {
        cover: [
          {
            required: true,
            message: this.$t('valid.cover'),
            trigger: ['change', 'blur']
          }
        ]
      }
    }
  },
  methods: {
    handleChange(file, fileList) {
      // 上传自动覆盖
      this.fileList = [file]
    },
    beforeUpload(file) {
      if (!file.name.endsWith('.json')) {
        Message({
          message: this.$t('confirm.uploadWebFileFormat'),
          type: 'error',
          duration: 5 * 1000
        })
        return false
      }
      if (file.size / 1024 / 1024 > 10) {
        Message({
          message: this.$t('confirm.uploadWebFileSize').toString(),
          type: 'error',
          duration: 5 * 1000
        })
        return false
      }
    },
    submitImport() {
      this.$refs.upload.submit()
    }
  }
}
</script>

<style scoped>
.el-button {
  margin-left: 10px;
}
</style>
