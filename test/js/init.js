import { getDate } from "date.js";
import html2md from 'https://cdn.jsdelivr.net/npm/html2md@0.1.1/+esm'

const args = GetArgs();

async function fn() {
  console.log(typeof html2md);

  const response = HttpRequest()
    .SetURL("https://github.com/JulioGuillermo/JGBot")
    .Get();
  print(response.StatusCode == 200 ? 'Is ok' : 'not ok');

  const html = response.BodyString();
  const md = html2md(html);
  return md;
}

export default fn();
