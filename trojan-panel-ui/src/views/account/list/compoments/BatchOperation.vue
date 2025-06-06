<template>
  <el-dialog
    :title="$t('table.createBatch').toString()"
    :visible="dialogFormVisibleProps"
    @close="$emit('update:dialogFormVisibleProps', false)"
  >
    <el-form
      ref="dataForm"
      :rules="createRules"
      :model="temp"
      label-position="left"
    >
      <el-form-item :label="$t('table.num').toString()" prop="num">
        <el-input-number
          v-model.number="temp.num"
          controls-position="right"
          type="number"
        />
      </el-form-item>
      <el-form-item :label="$t('table.presetQuota').toString()" prop="presetQuota">
        <el-input-number
          v-model.number="temp.presetQuota"
          controls-position="right"
          type="number"
        />
      </el-form-item>
      <el-form-item
        :label="$t('table.presetExpire').toString()"
        prop="presetExpire"
      >
        <el-input-number
          v-model.number="temp.presetExpire"
          controls-position="right"
          type="number"
        />
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="$emit('update:dialogFormVisibleProps', false)"
        >{{ $t('table.cancel') }}
      </el-button>
      <el-button type="primary" @click="createAccountBatch">
        {{ $t('table.confirm') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script>
import { createAccountBatch } from '@/api/account'

export default {
  name: 'BatchOperation',
  props: {
    dialogFormVisibleProps: {
      type: Boolean,
      required: true
    },
    getListProps: {
      type: Function,
      required: true
    }
  },
  data() {
    return {
      temp: {
        num: 5,
        presetQuota: 1024,
        presetExpire: 7
      },
      createRules: {
        num: [
          {
            required: true,
            message: this.$t('valid.createBatchNum'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 5,
            max: 500,
            message: this.$t('valid.createBatchNumRange'),
            trigger: ['change', 'blur']
          }
        ],
        presetQuota: [
          {
            required: true,
            message: this.$t('valid.createBatchPresetQuota'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: -1,
            max: 1024000,
            message: this.$t('valid.createBatchPresetQuotaRange'),
            trigger: ['change', 'blur']
          }
        ],
        presetExpire: [
          {
            required: true,
            message: this.$t('valid.createBatchPresetExpire'),
            trigger: ['change', 'blur']
          },
          {
            type: 'number',
            min: 1,
            max: 365,
            message: this.$t('valid.createBatchPresetExpireRange'),
            trigger: ['change', 'blur']
          }
        ]
      }
    }
  },
  methods: {
    clearValidate() {
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createAccountBatch() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          createAccountBatch(this.temp).then(() => {
            this.getListProps()
            this.$emit('update:dialogFormVisibleProps', false)
            this.$notify({
              title: 'Success',
              message: this.$t('confirm.createSuccess').toString(),
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
