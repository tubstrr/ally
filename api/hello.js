import { readFileSync } from "fs";
import path from "path";

export default function handler(req, res) {
  const file = path.join(process.cwd(), "ally/admin", "index.html");
  const stringified = readFileSync(file, "utf8");

  res.setHeader("Content-Type", "text/html; charset=utf-8");
  return res.end(stringified);
}
