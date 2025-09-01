import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

// 创建统一的Markdown渲染器
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true, // 支持换行
  highlight: function (str: string, lang: string): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        const highlighted = hljs.highlight(str, { language: lang }).value;
        return `<pre class="hljs"><code class="language-${lang}">${highlighted}</code></pre>`;
      } catch (__) {}
    }
    // 如果没有语言或高亮失败，使用默认的代码块
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`;
  }
})

// 统一的Markdown渲染函数
export function renderMarkdown(content: string): string {
  let html = md.render(content)
  
  // 为标题添加ID属性
  html = html.replace(
    /<h([1-6])>(.*?)<\/h[1-6]>/g,
    (_match: string, level: string, content: string) => {
      // 生成标题ID
      const text = content.replace(/<[^>]*>/g, '') // 移除HTML标签
      const id = text
        .toLowerCase()
        .replace(/[^\w\u4e00-\u9fa5]+/g, '-') // 替换非字母数字字符为连字符
        .replace(/^-+|-+$/g, '') // 移除首尾连字符
        .substring(0, 50) // 限制长度
      
      return `<h${level} id="${id}">${content}</h${level}>`
    }
  )
  
  // 处理代码块，添加语言标识和行号
  html = html.replace(
    /<pre class="hljs"><code class="language-(\w+)">([\s\S]*?)<\/code><\/pre>/g,
    (_match: string, lang: string, code: string) => {
      const lines: string[] = code.split('\n')
      
      const codeWithLineNumbers = lines.map((line: string, index: number) => {
        return `<div class="code-line">
          <span class="line-number">${index + 1}</span>
          <span class="line-content">${line || ' '}</span>
        </div>`
      }).join('\n')
      
      const result = `
        <div class="code-block" data-language="${lang}">
          <div class="code-header">
            <span class="language-label">${lang}</span>
            <button class="copy-btn" onclick="copyCode(this)">复制代码</button>
          </div>
          <div class="code-content">
            ${codeWithLineNumbers}
          </div>
        </div>
      `
      return result
    }
  )
  
  // 处理没有语言标识的代码块
  html = html.replace(
    /<pre class="hljs"><code>([\s\S]*?)<\/code><\/pre>/g,
    (_match: string, code: string) => {
      const lines: string[] = code.split('\n')
      
      const codeWithLineNumbers = lines.map((line: string, index: number) => {
        return `<div class="code-line">
          <span class="line-number">${index + 1}</span>
          <span class="line-content">${line || ' '}</span>
        </div>`
      }).join('\n')
      
      return `
        <div class="code-block" data-language="text">
          <div class="code-header">
            <span class="language-label">text</span>
            <button class="copy-btn" onclick="copyCode(this)">复制代码</button>
          </div>
          <div class="code-content">
            ${codeWithLineNumbers}
          </div>
        </div>
      `
    }
  )
  
  return html
}

// 复制代码功能
export function copyCode(button: HTMLElement): void {
  const codeBlock = button.closest('.code-block')
  if (codeBlock) {
    const codeContent = codeBlock.querySelector('.code-content')
    if (codeContent) {
      const textContent = codeContent.textContent || ''
      navigator.clipboard.writeText(textContent).then(() => {
        // 临时改变按钮文本
        const originalText = button.textContent
        button.textContent = '已复制!'
        button.style.backgroundColor = '#67c23a'
        button.style.color = 'white'
        
        setTimeout(() => {
          button.textContent = originalText
          button.style.backgroundColor = ''
          button.style.color = ''
        }, 2000)
      }).catch(() => {
        console.error('复制失败')
      })
    }
  }
}

// 获取纯文本摘要（移除Markdown语法）
export function getPlainTextSummary(content: string, maxLength: number = 140): string {
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

// 获取HTML格式的摘要（保留基本格式）
export function getHtmlSummary(content: string, maxLength: number = 140): string {
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

// 导出Markdown渲染器实例（如果需要直接使用）
export { md }
