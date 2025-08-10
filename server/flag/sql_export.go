package flag

import (
	"fmt"
	"os"
	"os/exec"
	"server/global"
	"strings"

	"go.uber.org/zap"
)

// exportMySQL 导出MySQL数据到SQL文件
func exportMySQL(filePath, tables string) error {
	mysqlConf := global.Config.Mysql

	args := []string{
		fmt.Sprintf("-u%s", mysqlConf.Username),
		fmt.Sprintf("-p%s", mysqlConf.Password),
		fmt.Sprintf("-h%s", mysqlConf.Host),
		fmt.Sprintf("-P%d", mysqlConf.Port),
		"--default-character-set=utf8mb4",
		"--ssl-mode=DISABLED",
	}

	if tables == "" {
		args = append(args, mysqlConf.DBName)
	} else {
		args = append(args, mysqlConf.DBName, "--tables", strings.ReplaceAll(tables, ",", " "))
	}

	args = append(args, "--result-file", filePath)

	// global.ZapLog.Info("即将执行mysqldump命令", zap.Strings("参数", args), zap.String("数据库名称", mysqlConf.DBName))

	// global.ZapLog.Info("MySQL配置信息", zap.Any("mysqlConf", mysqlConf))

	cmd := exec.Command("mysqldump", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errorMsg := fmt.Sprintf("数据导出失败: %v\n错误详情: %s", err, string(output))
		global.ZapLog.Error(errorMsg)
		return fmt.Errorf(errorMsg)
	}

	if _, err := os.Stat(filePath); err != nil || os.IsNotExist(err) {
		msg := fmt.Sprintf("导出文件不存在: %s", filePath)
		global.ZapLog.Error(msg)
		return fmt.Errorf(msg)
	}

	global.ZapLog.Info("MySQL数据导出成功", zap.String("文件路径", filePath))
	return nil
}
