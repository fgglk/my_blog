@echo off
echo 恢复文件...

:: 检查是否有备份目录
if not exist "backup_*" (
    echo 没有找到备份目录！
    pause
    exit /b 1
)

:: 显示可用的备份
echo 可用的备份目录：
dir /b backup_* | findstr /r "backup_.*"

:: 询问用户选择备份目录
set /p backup_dir="请输入要恢复的备份目录名称: "

if not exist "%backup_dir%" (
    echo 备份目录不存在！
    pause
    exit /b 1
)

echo 正在从 %backup_dir% 恢复文件...

:: 恢复前端文件
echo 恢复前端文件...
copy "%backup_dir%\ArticleDetail.vue" "web\src\views\ArticleDetail.vue" /Y
copy "%backup_dir%\Write.vue" "web\src\views\Write.vue" /Y
copy "%backup_dir%\article.ts" "web\src\types\article.ts" /Y

:: 恢复后端文件
echo 恢复后端文件...
copy "%backup_dir%\article_service.go" "server\service\article.go" /Y
copy "%backup_dir%\article_api.go" "server\api\article.go" /Y
copy "%backup_dir%\article_request.go" "server\model\request\article.go" /Y

echo 恢复完成！
pause
