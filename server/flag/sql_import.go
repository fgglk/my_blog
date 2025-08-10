package flag

import (
	"fmt"
	"os"
	"os/exec"
	"server/global"

	"go.uber.org/zap"
)

// importMySQL 从SQL文件导入数据到MySQL
func importMySQL(filePath string) error {
	mysqlConf := global.Config.Mysql

	// 检查导入文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("导入文件不存在: %s", filePath)
	}

	// 构建mysql命令参数
	args := []string{
		fmt.Sprintf("-u%s", mysqlConf.Username),
		fmt.Sprintf("-p%s", mysqlConf.Password),
		fmt.Sprintf("-h%s", mysqlConf.Host),
		fmt.Sprintf("-P%d", mysqlConf.Port),
		"--default-character-set=utf8mb4",
		"--ssl-mode=DISABLED",
		mysqlConf.DBName,
	}

	// 打开SQL文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开导入文件: %v", err)
	}
	defer file.Close()

	// 执行mysql命令导入数据
	cmd := exec.Command("mysql", args...)
	cmd.Stdin = file
	output, err := cmd.CombinedOutput()
	if err != nil {
		errorMsg := fmt.Sprintf("数据导入失败: %v\n错误详情: %s", err, string(output))
		global.ZapLog.Error(errorMsg)
		return fmt.Errorf(errorMsg)
	}

	global.ZapLog.Info("MySQL数据导入成功", zap.String("文件路径", filePath))
	return nil
}
