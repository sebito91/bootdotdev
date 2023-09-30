const {argv} = require('node:process');
const {crawlPage} = require('./crawl.js');

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
  const pageset = Object.keys(pages);

  pageset.forEach((page) => {
    console.log(`Found ${pages[page]} visits to ${page}`);
  });
}

main();
