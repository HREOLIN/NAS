import { readFile } from "node:fs/promises";
import path from "node:path";

export function json(response, statusCode, payload) {
  response.writeHead(statusCode, { "Content-Type": "application/json; charset=utf-8" });
  response.end(JSON.stringify(payload, null, 2));
}

export async function readJsonBody(request) {
  const chunks = [];

  for await (const chunk of request) {
    chunks.push(chunk);
  }

  if (chunks.length === 0) {
    return {};
  }

  const raw = Buffer.concat(chunks).toString("utf8");
  return JSON.parse(raw);
}

export async function serveStatic(response, publicDir, pathname) {
  const filePath = path.join(publicDir, pathname === "/" ? "index.html" : pathname.slice(1));
  const ext = path.extname(filePath).toLowerCase();
  const contentType = {
    ".html": "text/html; charset=utf-8",
    ".css": "text/css; charset=utf-8",
    ".js": "application/javascript; charset=utf-8"
  }[ext];

  if (!contentType) {
    return false;
  }

  try {
    const content = await readFile(filePath);
    response.writeHead(200, { "Content-Type": contentType });
    response.end(content);
    return true;
  } catch {
    return false;
  }
}

