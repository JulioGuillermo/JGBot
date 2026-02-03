import showdown from 'https://esm.sh/showdown';

const args = GetArgs();

async function fn() {
  const response = HttpRequest()
    .SetURL("https://github.com/JulioGuillermo/JGBot")
    .Get();
  print(response.StatusCode == 200 ? 'Is ok' : 'not ok');

  const converter = new showdown.Converter();
  const html = response.BodyString();
  const markdown = converter.makeMarkdown(html);
  return markdown;
}

export default fn();
