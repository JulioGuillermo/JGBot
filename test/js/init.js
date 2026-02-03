import { webSearch } from "./search.js";
import { html2md } from "./html2md.js";

const args = GetArgs();

const exec = () => {
  const html = webSearch("noticias de hoy");
  return html2md(html);
};

export default exec();