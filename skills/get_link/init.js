import { TurndownService } from "/turndown.js";

const args = GetArgs();

const exec = () => {
  const response = HttpRequest().SetURL(args.url).Get();

  if (response.StatusCode != 200 && response.StatusCode != 201) {
    return `Request return status code ${response.StatusCode} with status message ${response.Status}`;
  }

  const content = response.BodyString();

  var turndownService = new TurndownService();
  var markdown = turndownService.turndown(content);

  return markdown;
};

export default exec();
