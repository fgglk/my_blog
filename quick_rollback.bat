@echo off
echo 快速回退到上一个版本...

:: 检查Git状态
git status

echo.
echo 当前提交历史：
git log --oneline -5

echo.
set /p confirm="确认要回退到上一个版本吗？(y/n): "

if /i "%confirm%"=="y" (
    echo 正在回退...
    git reset --hard HEAD~1
    echo 回退完成！
) else (
    echo 取消回退操作。
)

pause
