@echo off
echo 创建备份文件...

:: 创建备份目录
set backup_dir=backup_%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%
set backup_dir=%backup_dir: =0%
mkdir %backup_dir%

:: 备份前端重要文件
echo 备份前端文件...
xcopy "web\src\views\ArticleDetail.vue" "%backup_dir%\ArticleDetail.vue" /Y
xcopy "web\src\views\Write.vue" "%backup_dir%\Write.vue" /Y
xcopy "web\src\types\article.ts" "%backup_dir%\article.ts" /Y

:: 备份后端重要文件
echo 备份后端文件...
xcopy "server\service\article.go" "%backup_dir%\article_service.go" /Y
xcopy "server\api\article.go" "%backup_dir%\article_api.go" /Y
xcopy "server\model\request\article.go" "%backup_dir%\article_request.go" /Y

echo 备份完成！备份目录：%backup_dir%
pause
