const {argv} = require('node:process');
const {crawlPage} = require('./crawl.js');

function main() {
  if (argv.length != 3) {
    throw new Error(
      `expected a single baseURL, got ${argv.length - 2} args instead`,
    );
  }

  console.log(`Calling the crawler for ${argv[2]}...`);
  crawlPage(argv[2]);
}

main();
