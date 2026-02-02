import { getDate } from "/date.js";

const args = GetArgs();

async function fn() {
  print(args);
  console.Log("Running #### 123");

  const formData = HttpFormData()
    .AddField("name", "test")
    .AddFile("file", "test.txt", new Uint8Array([1, 2, 3, 4, 5]));

  const response = HttpRequest()
    .SetURL("https://httpbin.org/get")
    .SetBodyFormData(formData)
    .Get();
  console.log(response.BodyString());
  const d = await getDate();
  onResult(d);
  return d;
}

export default await fn();
