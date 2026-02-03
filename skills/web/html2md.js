import { decode } from 'https://esm.sh/he';

export function html2md(html) {
    return decode(html)
        .replace(/<br>/gi, '\n')
        .replace(/<head[^>]*>([\s\S]*?)<\/head>/gi, '')
        .replace(/<script[^>]*>([\s\S]*?)<\/script>/gi, '')
        .replace(/<style[^>]*>([\s\S]*?)<\/style>/gi, '')

        // 1. Code blocks (Pre/Code)
        .replace(/<pre[^>]*><code[^>]*>([\s\S]*?)<\/code><\/pre>/gi, '```\n$1\n```\n')
        .replace(/<code[^>]*>(.*?)<\/code>/gi, '`$1`')

        // 2. Headers (H1 - H6)
        .replace(/<h1[^>]*>(.*?)<\/h1>/gi, '\n# $1\n')
        .replace(/<h2[^>]*>(.*?)<\/h2>/gi, '\n## $1\n')
        .replace(/<h3[^>]*>(.*?)<\/h3>/gi, '\n### $1\n')
        .replace(/<h4[^>]*>(.*?)<\/h4>/gi, '\n#### $1\n')
        .replace(/<h5[^>]*>(.*?)<\/h5>/gi, '\n##### $1\n')
        .replace(/<h6[^>]*>(.*?)<\/h6>/gi, '\n###### $1\n')

        // 3. Lists (Simple)
        .replace(/<li[^>]*>(.*?)<\/li>/gi, '* $1\n')

        // 4. Text format
        .replace(/<strong[^>]*>(.*?)<\/strong>|<b[^>]*>(.*?)<\/b>/gi, '**$1$2**')
        .replace(/<em[^>]*>(.*?)<\/em>|<i[^>]*>(.*?)<\/i>/gi, '*$1$2*')

        // 5. Structure and cleaning
        .replace(/<p[^>]*>(.*?)<\/p>/gi, '$1\n\n')
        .replace(/<br\s*\/?>/gi, '\n')
        .replace(/<a[^>]*href=["']([^"']*)["'][^>]*>([\s\S]*?)<\/a>/gi, (match, url, text) => {
            if (url.startsWith('#')) return text;

            text = text.replace(/[\n\r]*/g, '');
            if (!text.trim() || text.includes('<svg')) return `[${url}](${url})`;

            return `[${text.trim()}](${url})`;
        })
        .replace(/<[^>]+>/g, '') // Remove any remaining tags (ul, ol, div, etc.)

        // 6. Normalization of spaces and line breaks
        .replace(/\r/g, '')               // Remove Windows carriage returns
        .replace(/[ \t]+\n/g, '\n')       // Remove spaces/tabs at the end of each line
        .replace(/\n{3,}/g, '\n\n')       // Collapse 3 or more line breaks into 2
        .trim();                          // Clean start and end
}