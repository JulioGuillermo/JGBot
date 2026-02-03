import { html2md } from "./html2md";
import { fetchContent } from "./fetch";

const args = GetArgs();

const searchUrl = (query) => {
  return `https://www.mojeek.com/search?q=${encodeURIComponent(query)}`;
}

const exec = () => {
  let url = ''
  let headers = false

  if (args.url) {
    url = args.url
  } else if (args.query) {
    url = searchUrl(args.query)
    headers = true
  } else {
    return "No url or query provided"
  }

  const html = fetchContent(url, headers);

  return html2md(html);
};

export default exec();
