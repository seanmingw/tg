<template>
  <div class="app-container">
    <el-upload
      class="upload-demo"
      ref="upload"
      action=""
      :file-list="fileList"
      :http-request="uploadFile"
      :auto-upload="false"
      accept=".zip"
      :on-change="handleChange"
      :before-upload="beforeUpload"
      :limit="1"
    >
      <el-button slot="trigger" size="small" type="primary"
        >{{ $t('config.webFileSelect') }}
      </el-button>
      <el-button
        size="small"
        type="success"
        @click="submitUpload"
        :disabled="fileList.length === 0"
      >
        {{ $t('config.webFileBtn') }}
      </el-button>
      <div slot="tip" class="el-upload__tip">
        {{ $t('config.webFileTip') }}
      </div>
    </el-upload>
  </div>
</template>

<script>
import { uploadWebFile } from '@/api/system'
import { Message } from 'element-ui'

export default {
  name: 'web-file',
  data() {
    return {
      fileList: []
    }
  },
  methods: {
    submitUpload() {
      this.$refs.upload.submit()
    },
    uploadFile(params) {
      let formData = new FormData()
      formData.append('file', params.file)
      uploadWebFile(formData).then(() => {
        this.$notify({
          title: 'Success',
          message: this.$t('confirm.uploadWebFileSuccess'),
          type: 'success',
          duration: 2000
        })
      })
      this.fileList = []
    },
    handleChange(file, fileList) {
      // 上传自动覆盖
      this.fileList = [file]
    },
    beforeUpload(file) {
      if (!file.name.endsWith('.zip')) {
        Message({
          message: this.$t('confirm.uploadWebFileFormat'),
          type: 'error',
          duration: 5 * 1000
        })
        return false
      }
      if (file.size / 1024 / 1024 > 10) {
        Message({
          message: this.$t('confirm.uploadWebFileSize'),
          type: 'error',
          duration: 5 * 1000
        })
        return false
      }
    }
  }
}
</script>

<style scoped>
.el-button {
  margin-left: 10px;
}
</style>
