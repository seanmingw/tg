<template>
  <el-upload
    list-type="picture-card"
    ref="upload"
    action=""
    :file-list="fileList"
    :http-request="uploadFile"
    accept=".png"
    :on-change="handleChange"
    :before-upload="beforeUpload"
    :limit="1"
    :disabled="uploadDisable"
  >
    <i slot="default" class="el-icon-plus"></i>
    <div slot="file" slot-scope="{ file }">
      <img class="el-upload-list__item-thumbnail" :src="file.url" alt="" />
      <span class="el-upload-list__item-actions">
        <span class="el-upload-list__item-delete" @click="handleRemove(file)">
          <i class="el-icon-delete"></i>
        </span>
      </span>
    </div>
    <div slot="tip" class="el-upload__tip">
      {{ $t('config.imageFileTip') }}
    </div>
  </el-upload>
</template>

<script>
import { Message } from 'element-ui'
import { uploadLogo } from '@/api/system'

export default {
  name: 'index',
  data() {
    return {
      fileList: [{ name: 'logo.png', url: '/api/image/logo' }],
      uploadDisable: true
    }
  },
  watch: {
    fileList: function (val) {
      this.uploadDisable = val.length >= 1
    }
  },
  methods: {
    handleRemove(file) {
      this.fileList = []
    },
    beforeUpload(file) {
      const isJPG = file.type === 'image/png'
      const isLt2M = file.size / 1024 / 1024 < 3

      if (!isJPG) {
        Message({
          message: this.$t('confirm.uploadWebFileFormat'),
          type: 'error',
          duration: 5 * 1000
        })
      }
      if (!isLt2M) {
        Message({
          message: this.$t('confirm.uploadWebFileSize'),
          type: 'error',
          duration: 5 * 1000
        })
      }
      return isJPG && isLt2M
    },
    handleChange(file, fileList) {
      // 上传自动覆盖
      this.fileList = [file]
    },
    uploadFile(params) {
      if (this.fileList.length > 0) {
        let formData = new FormData()
        formData.append('file', this.fileList[0].raw)
        uploadLogo(formData).then(() => {
          this.$notify({
            title: 'Success',
            message: this.$t('confirm.modifySuccess'),
            type: 'success',
            duration: 2000
          })
        })
      }
    }
  }
}
</script>
<style></style>
