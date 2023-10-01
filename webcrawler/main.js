const {argv} = require('node:process');
const {crawlPage} = require('./crawl.js');
const {printReport, tester} = require('./report.js');

async function startCrawl(baseURL, currentURL) {
  const pages = await crawlPage(baseURL, currentURL, {});
  return pages;
}

async function main() {
  if (argv.length != 3) {
    throw new Error(
      `expected a single baseURL, got ${argv.length - 2} args instead`,
    );
  }

  const baseURL = argv[2];
  const currentURL = argv[2];

  console.log(`Calling the crawler for ${baseURL}...`);

  const pages = await crawlPage(baseURL, currentURL, {});
  printReport(pages);
}

main();
