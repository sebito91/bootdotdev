const {JSDOM} = require('jsdom');

function normalizeURL(url) {
  let newURL = new URL(url);
  let returnURL = `${newURL.hostname}`;

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
    console.log(`SEBTEST -- we found the following at ${i}: ${urls_found[i]}`);
    found_urls.push(`${urls_found[i]}`);
  }

  return found_urls;
}

module.exports = {
  normalizeURL,
  getURLsFromHTML,
};
