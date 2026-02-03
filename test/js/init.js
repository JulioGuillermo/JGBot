import { decode } from 'https://esm.sh/he';

function html2md(html) {
  return decode(html)
    .replace(/<h1[^>]*>(.*?)<\/h1>/gi, '# $1\n')
    .replace(/<strong[^>]*>(.*?)<\/strong>/gi, '**$1**')
    .replace(/<p[^>]*>(.*?)<\/p>/gi, '$1\n\n')
    .replace(/<[^>]+>/g, ''); // Limpia cualquier otra etiqueta
}

function fn() {
  const response = HttpRequest()
    .SetURL("https://github.com/JulioGuillermo/JGBot")
    .Get();
  print(response.StatusCode == 200 ? 'Is ok' : 'not ok');

  const html = response.BodyString();
  const markdown = html2md(html);
  console.log(markdown);
  return markdown;
}

export default fn();
