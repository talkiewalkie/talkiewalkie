#!/usr/bin/env node

const puppeteer = require("puppeteer");

let args = process.argv.slice(2);
if (args.length !== 1) {
  throw Error("expected usage: `node file.js customToken`");
}

let customToken = args[0];

(async () => {

  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto(
    `file://${process.cwd()}/index.html?customToken=${customToken}`
  );
  const el = await page.waitForSelector("#idToken", {
    visible: true,
    timeout: 5 * 1_000,
  });
  const text = await page.evaluate((el) => el.textContent, el);
  console.log(text);

  await browser.close();
})();
