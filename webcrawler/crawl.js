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

module.exports = {
  normalizeURL,
};
