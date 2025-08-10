# 代码回退指南

## 概述
为了避免修改代码时出现错误导致无法恢复，我们提供了多种回退方案。

## 方案一：Git 版本控制（推荐）

### 基本操作
1. **查看当前状态**
   ```bash
   git status
   ```

2. **提交当前修改**
   ```bash
   git add .
   git commit -m "描述您的修改"
   ```

3. **查看提交历史**
   ```bash
   git log --oneline
   ```

4. **回退到上一个版本**
   ```bash
   git reset --hard HEAD~1
   ```

5. **回退到特定版本**
   ```bash
   git reset --hard <commit-hash>
   ```

### 使用快速回退脚本
双击运行 `quick_rollback.bat` 文件，按提示操作即可快速回退。

## 方案二：文件备份恢复

### 创建备份
双击运行 `backup_script.bat` 文件，会自动创建带时间戳的备份目录。

### 恢复备份
双击运行 `restore_script.bat` 文件，选择要恢复的备份目录。

## 方案三：手动操作

### 备份重要文件
在进行重要修改前，手动复制以下文件：
- `web/src/views/ArticleDetail.vue`
- `web/src/views/Write.vue`
- `web/src/types/article.ts`
- `server/service/article.go`
- `server/api/article.go`
- `server/model/request/article.go`

### 恢复文件
将备份的文件复制回原位置即可。

## 最佳实践

1. **每次修改前提交代码**
   ```bash
   git add .
   git commit -m "修改前的状态"
   ```

2. **小步提交**
   每完成一个小功能就提交一次，便于回退。

3. **使用描述性的提交信息**
   ```bash
   git commit -m "修复文章保存为草稿的问题"
   ```

4. **定期创建备份**
   在进行重大修改前，使用备份脚本创建备份。

## 常见问题

### Q: 回退后想恢复被删除的修改怎么办？
A: 使用 `git reflog` 查看操作历史，然后使用 `git reset --hard <commit-hash>` 恢复到指定版本。

### Q: 只想回退部分文件怎么办？
A: 使用 `git checkout HEAD~1 -- <文件路径>` 恢复特定文件。

### Q: 备份文件太多怎么办？
A: 定期清理旧的备份目录，只保留最近的几个备份。

## 注意事项

1. **Git 回退会丢失未提交的修改**，请确保重要修改已提交。
2. **备份文件会占用磁盘空间**，定期清理。
3. **在生产环境中谨慎使用 `git reset --hard`**，建议使用 `git revert`。
