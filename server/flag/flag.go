package flag

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

// 命令行标志定义
var (
	sqlFlag = &cli.BoolFlag{
		Name:  "sql",
		Usage: "初始化Mysql数据库表结构",
	}
	exportFlag = &cli.BoolFlag{
		Name:  "export",
		Usage: "导出MySQL数据到SQL文件",
	}
	exportPathFlag = &cli.StringFlag{
		Name:  "export-path",
		Usage: "指定导出文件路径",
		Value: "backup.sql",
	}
	tablesFlag = &cli.StringFlag{
		Name:  "tables",
		Usage: "指定要导出的表，逗号分隔（如：users,articles）",
	}
	createAdminFlag = &cli.BoolFlag{
		Name:  "create-admin",
		Usage: "创建管理员账户",
	}
	importFlag = &cli.BoolFlag{
		Name:  "import",
		Usage: "从SQL文件导入MySQL数据",
	}
	importPathFlag = &cli.StringFlag{
		Name:  "import-path",
		Usage: "指定导入文件路径",
		Value: "backup.sql",
	}
	createEsIndexFlag = &cli.BoolFlag{
		Name:  "create-es-index",
		Usage: "创建Elasticsearch索引结构",
	}
	exportEsFlag = &cli.BoolFlag{
		Name:  "export-es",
		Usage: "导出Elasticsearch数据到文件",
	}
	exportEsPathFlag = &cli.StringFlag{
		Name:  "export-es-path",
		Usage: "指定Elasticsearch导出文件路径",
		Value: "es_backup.json",
	}
	importEsFlag = &cli.BoolFlag{
		Name:  "import-es",
		Usage: "从文件导入数据到Elasticsearch",
	}
	importEsPathFlag = &cli.StringFlag{
		Name:  "import-es-path",
		Usage: "指定Elasticsearch导入文件路径",
		Value: "es_backup.json",
	}
)

// NewApp 创建CLI应用实例
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "blog_cli"
	app.Usage = "博客系统命令行工具"
	app.Version = "1.0.0"
	return app
}

// setupFlags 配置命令行标志
func setupFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		sqlFlag,
		exportFlag,
		exportPathFlag,
		tablesFlag,
		createAdminFlag,
		importFlag,
		importPathFlag,
		createEsIndexFlag,
		exportEsFlag,
		exportEsPathFlag,
		importEsFlag,
		importEsPathFlag,
	}
}

// setupActions 配置命令行动作处理
func setupActions(app *cli.App) {
	app.Action = func(c *cli.Context) error {
		// 处理SQL迁移命令
		if c.Bool("sql") {
			return migrateDatabase()
		}
		if c.Bool("export") {
			filePath := c.String("export-path")
			tables := c.String("tables")
			if err := exportMySQL(filePath, tables); err != nil {
				return fmt.Errorf("数据导出失败: %v", err)
			}
			return nil
		}
		if c.Bool("create-admin") {
			if err := CreateAdministrator(); err != nil {
				return fmt.Errorf("创建管理员失败: %v", err)
			}
			fmt.Println("管理员账户创建成功")
			return nil
		}
		if c.Bool("import") {
			filePath := c.String("import-path")
			if err := importMySQL(filePath); err != nil {
				return fmt.Errorf("数据导入失败: %v", err)
			}
			return nil
		}
		if c.Bool("create-es-index") {
			if err := createEsIndex(); err != nil {
				return fmt.Errorf("创建ES索引失败: %v", err)
			}
			return nil
		}
		if c.Bool("export-es") {
			filePath := c.String("export-es-path")
			if err := exportEsData(filePath); err != nil {
				return fmt.Errorf("ES数据导出失败: %v", err)
			}
			fmt.Printf("ES数据已成功导出至 %s\n", filePath)
			return nil
		}
		if c.Bool("import-es") {
			filePath := c.String("import-es-path")
			if err := importEsData(filePath); err != nil {
				return fmt.Errorf("ES数据导入失败: %v", err)
			}
			fmt.Printf("ES数据已成功从 %s 导入\n", filePath)
			return nil
		}
		return cli.ShowAppHelp(c)

	}
}

// Run 启动CLI应用
func Run() {
	// 创建应用
	app := NewApp()
	// 配置标志
	setupFlags(app)
	// 配置动作
	setupActions(app)
	// 运行应用
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("命令执行失败: %v\n", err)
		os.Exit(1)
	}
}
