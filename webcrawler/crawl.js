const {JSDOM} = require('jsdom');

function normalizeURL(url) {
  let newURL = new URL(url);
  let returnURL = `${newURL.protocol}//${newURL.hostname}`;

  if (newURL.port) {
    returnURL = returnURL.concat(`:${newURL.port}`);
  }

  if (newURL.pathname) {
    returnURL = returnURL.concat(`${newURL.pathname}`.replace(/\/+$/, ''));
  }

  return returnURL;
}

function getURLsFromHTML(htmlBody, baseURL) {
  const htmlDOM = new JSDOM(htmlBody, {
    url: `${baseURL}`,
    contentType: 'text/html',
    includeNodeLocations: true,
  });

  let found_urls = [];
  urls_found = htmlDOM.window.document.querySelectorAll('a');
  for (let i = 0; i < urls_found.length; i++) {
    found_urls.push(`${urls_found[i]}`);
  }

  return found_urls;
}

async function crawlPage(baseURL, currentURL, pages) {
  await new Promise((resolve) => setTimeout(resolve, 1000));

  if (!currentURL.startsWith(baseURL)) {
    return pages;
  }

  const currentURL_normalized = normalizeURL(currentURL);

  // if we already exist in the pages object, increment and return
  // else if we are the same as the baseURL we set to 0
  // else we set the currentURL to 1
  if (currentURL_normalized in pages) {
    pages[currentURL_normalized]++;
    return pages;
  } else if (baseURL == currentURL) {
    pages[currentURL_normalized] = 0;
  } else {
    pages[currentURL_normalized] = 1;
  }

  const resp = await fetch(currentURL_normalized);
  if (!resp.ok) {
    console.log(
      `received error code ${resp.status} for ${currentURL_normalized}`,
    );
    return pages;
  }

  if (!resp.headers.get('Content-Type').includes('text/html')) {
    console.log(
      `received different Content-Type header from ${currentURL_normalized}, got ${resp.headers.get(
        'Content-Type',
      )}, expected 'text/html'`,
    );
    return pages;
  }

  const htmlText = await resp.text();
  found_urls = getURLsFromHTML(htmlText, baseURL);

  for (let i = 0; i < found_urls.length; i++) {
    try {
      pages = await crawlPage(baseURL, found_urls[i], pages);
    } catch (error) {
      console.log(`${error}`);
    }
  }

  return pages;
}

module.exports = {
  normalizeURL,
  getURLsFromHTML,
  crawlPage,
};
