import { getDate } from "/date.js";

async function fn() {
  console.log("Running");
  const d = await getDate();
  onResult(d);
  return d;
}

export default fn();
