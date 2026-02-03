import { getDate } from "/date.js";
require("https://cdn.jsdelivr.net/npm/node-html-markdown@2.0.0/+esm");

const args = GetArgs();

async function fn() {
  print(args);
  console.Log("Running #### 123");

  VF.WriteStrFile("/test.md", "hello world");
  print(VF.ReadDir("./"));
  try {
    print(VF.ReadStrFile("/test.md"));
  } catch (e) {
    console.log(e);
  }

  const formData = HttpFormData()
    .AddField("name", "test")
    .AddFile("file", "test.txt", new Uint8Array([1, 2, 3, 4, 5]));

  const response = HttpRequest()
    .SetURL("https://httpbin.org/get")
    .SetBodyFormData(formData)
    .Get();
  print(response.StatusCode == 200 ? "Is ok" : "not ok");
  console.log(response.BodyString());
  const d = await getDate();
  return d;
}

export default fn();
