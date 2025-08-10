import MarkdownIt from 'markdown-it'

// 创建Markdown渲染器实例
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true // 支持换行
})

/**
 * 渲染Markdown内容为HTML
 * @param content Markdown内容
 * @returns 渲染后的HTML字符串
 */
export const renderMarkdown = (content: string): string => {
  if (!content) return ''
  return md.render(content)
}

/**
 * 获取纯文本摘要（移除Markdown语法）
 * @param content Markdown内容
 * @param maxLength 最大长度
 * @returns 纯文本摘要
 */
export const getPlainTextSummary = (content: string, maxLength: number = 140): string => {
  if (!content) return ''
  
  // 移除Markdown语法
  let plainText = content
    .replace(/!\[([^\]]*)\]\([^)]*\)/g, '$1') // 移除图片语法，保留alt文本
    .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1') // 移除链接语法，保留链接文本
    .replace(/\*\*([^*]+)\*\*/g, '$1') // 移除粗体
    .replace(/\*([^*]+)\*/g, '$1') // 移除斜体
    .replace(/`([^`]+)`/g, '$1') // 移除行内代码
    .replace(/^#{1,6}\s+/gm, '') // 移除标题语法
    .replace(/^[-*+]\s+/gm, '') // 移除列表语法
    .replace(/^\d+\.\s+/gm, '') // 移除有序列表语法
    .replace(/^>\s+/gm, '') // 移除引用语法
    .replace(/```[\s\S]*?```/g, '') // 移除代码块
    .replace(/~~([^~]+)~~/g, '$1') // 移除删除线
    .replace(/\n+/g, ' ') // 将多个换行替换为单个空格
    .trim()
  
  // 截取指定长度
  if (plainText.length > maxLength) {
    plainText = plainText.substring(0, maxLength) + '...'
  }
  
  return plainText
}

/**
 * 获取HTML格式的摘要（保留基本格式）
 * @param content Markdown内容
 * @param maxLength 最大长度
 * @returns HTML格式的摘要
 */
export const getHtmlSummary = (content: string, maxLength: number = 140): string => {
  if (!content) return ''
  
  // 先获取纯文本摘要
  const plainSummary = getPlainTextSummary(content, maxLength)
  
  // 如果内容较短，直接渲染
  if (content.length <= maxLength) {
    return renderMarkdown(content)
  }
  
  // 否则返回纯文本摘要
  return plainSummary
}
