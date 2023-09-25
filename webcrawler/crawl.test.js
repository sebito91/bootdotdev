const {test, expect} = require('@jest/globals');
const {normalizeURL} = require('./crawl.js');
const {getURLsFromHTML} = require('./crawl.js');

// tests for the normalizeURL function
test('normalizeURL: normal URL + normal path', () => {
  expect(normalizeURL('https://this.is.a.test.com/tester')).toBe(
    'this.is.a.test.com/tester',
  );
});

test('normalizeURL: normal URL + normal path + trailing /', () => {
  expect(normalizeURL('https://this.is.a.test.com/tester/')).toBe(
    'this.is.a.test.com/tester',
  );
});

test('normalizeURL: normal URL + normal path + query', () => {
  expect(
    normalizeURL('https://this.is.a.test.com/tester?tester=a&newtester=b'),
  ).toBe('this.is.a.test.com/tester');
});

test('normalizeURL: normal URL + diff port + normal path', () => {
  expect(normalizeURL('https://this.is.a.test.com:9191/tester')).toBe(
    'this.is.a.test.com:9191/tester',
  );
});

test('normalizeURL: normal URL + no path/query', () => {
  expect(normalizeURL('https://this.is.a.test.com')).toBe('this.is.a.test.com');
});

// tests for the getURLsFromHTML function

test('getURLsFromHTML: found a single basic href', () => {
  expect(
    getURLsFromHTML(
      `<a href="https://boot.dev">Learn Backend Development</a>`,
      `https://boot.dev`,
    ),
  ).toStrictEqual(['https://boot.dev/']);
});

test('getURLsFromHTML: found two basic href', () => {
  expect(
    getURLsFromHTML(
      `<a href="https://boot.dev">Learn Backend Development</a>
      <a href="https://boot.org">Learn Backend Development Org</a>`,
      `https://boot.dev`,
    ),
  ).toStrictEqual(['https://boot.dev/', 'https://boot.org/']);
});
