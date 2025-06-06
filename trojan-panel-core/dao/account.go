package dao

import (
	"errors"
	"fmt"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan-panel-core/core"
	"trojan-panel-core/model"
	"trojan-panel-core/model/bo"
	"trojan-panel-core/model/constant"
	"trojan-panel-core/model/vo"
)

func UpdateAccountFlowByPassOrHash(pass *string, hash *string, download int, upload int) error {
	if download == 0 && upload == 0 {
		return nil
	}

	mySQLConfig := core.Config.MySQLConfig

	var values []interface{}
	downloadUpdateSql := ""
	if download != 0 {
		downloadUpdateSql = "download = download + ?"
		values = append(values, download)
	}
	uploadUpdateSql := ""
	if upload != 0 {
		if downloadUpdateSql == "" {
			uploadUpdateSql = "upload = upload + ?"
		} else {
			uploadUpdateSql = ",upload = upload + ?"
		}
		values = append(values, upload)
	}

	sql := fmt.Sprintf("update %s set %s where", mySQLConfig.AccountTable, downloadUpdateSql+uploadUpdateSql)

	if pass != nil && *pass != "" {
		sql += " pass = ?"
		values = append(values, *pass)
	}
	if hash != nil && *hash != "" {
		sql += " hash = ?"
		values = append(values, *hash)
	}
	_, err := db.Exec(sql, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}

// SelectAccounts query all valid accounts
func SelectAccounts() ([]bo.AccountBo, error) {
	mySQLConfig := core.Config.MySQLConfig
	var accounts []model.Account
	var (
		values []interface{}
		err    error
	)

	sql := fmt.Sprintf("select id,username,pass,hash from %s where quota < 0 or (quota > download + upload)", mySQLConfig.AccountTable)
	rows, err := db.Query(sql, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	if err = scanner.Scan(rows, &accounts); err != nil && err != scanner.ErrEmptyResult {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	accountBos := make([]bo.AccountBo, 0)
	if len(accounts) > 0 {
		for _, item := range accounts {
			accountBo := bo.AccountBo{
				Username: *item.Username,
				Pass:     *item.Pass,
				Hash:     *item.Hash,
			}
			accountBos = append(accountBos, accountBo)
		}
	}
	return accountBos, nil
}

func SelectAccountByPass(pass string) (*vo.AccountHysteriaVo, error) {
	mySQLConfig := core.Config.MySQLConfig
	var account model.Account

	buildSelect, values, err := builder.NamedQuery(fmt.Sprintf("select id,username,ip_limit,download_speed_limit,upload_speed_limit from %s where (quota < 0 or quota > download + upload) and pass = {{pass}}", mySQLConfig.AccountTable),
		map[string]interface{}{
			"pass": pass,
		})
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	err = scanner.Scan(rows, &account)
	if err == scanner.ErrEmptyResult {
		return nil, errors.New(constant.UsernameOrPassError)
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	ipLimit := 0
	if account.IpLimit != nil {
		ipLimit = *account.IpLimit
	}

	downloadSpeedLimit := 0
	if account.DownloadSpeedLimit != nil {
		downloadSpeedLimit = *account.DownloadSpeedLimit
	}

	uploadSpeedLimit := 0
	if account.UploadSpeedLimit != nil {
		uploadSpeedLimit = *account.UploadSpeedLimit
	}

	AccountHysteriaVo := vo.AccountHysteriaVo{
		Id:                 *account.Id,
		Username:           *account.Username,
		IpLimit:            ipLimit,
		DownloadSpeedLimit: downloadSpeedLimit,
		UploadSpeedLimit:   uploadSpeedLimit,
	}
	return &AccountHysteriaVo, nil
}

// SelectAccountByHash 根据hash查询账户信息
func SelectAccountByHash(hash string) (*model.Account, error) {
	mySQLConfig := core.Config.MySQLConfig
	var account model.Account

	buildSelect, values, err := builder.NamedQuery(fmt.Sprintf("select id,username,pass,hash,quota,download,upload,ip_limit,download_speed_limit,upload_speed_limit from %s where hash = {{hash}}", mySQLConfig.AccountTable),
		map[string]interface{}{
			"hash": hash,
		})
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	rows, err := db.Query(buildSelect, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}
	defer rows.Close()

	err = scanner.Scan(rows, &account)
	if err == scanner.ErrEmptyResult {
		return nil, nil
	} else if err != nil {
		logrus.Errorln(err.Error())
		return nil, errors.New(constant.SysError)
	}

	return &account, nil
}

// UpdateAccountFlowByHash 根据hash更新账户流量
func UpdateAccountFlowByHash(hash *string, download int, upload int) error {
	if download == 0 && upload == 0 {
		return nil
	}

	mySQLConfig := core.Config.MySQLConfig

	var values []interface{}
	downloadUpdateSql := ""
	if download != 0 {
		downloadUpdateSql = "download = download + ?"
		values = append(values, download)
	}
	uploadUpdateSql := ""
	if upload != 0 {
		if downloadUpdateSql == "" {
			uploadUpdateSql = "upload = upload + ?"
		} else {
			uploadUpdateSql = ",upload = upload + ?"
		}
		values = append(values, upload)
	}

	sql := fmt.Sprintf("update %s set %s where hash = ?", mySQLConfig.AccountTable, downloadUpdateSql+uploadUpdateSql)
	values = append(values, *hash)
	
	_, err := db.Exec(sql, values...)
	if err != nil {
		logrus.Errorln(err.Error())
		return errors.New(constant.SysError)
	}
	return nil
}
