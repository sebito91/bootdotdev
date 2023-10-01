const {test, expect} = require('@jest/globals');
const {normalizeURL} = require('./crawl.js');
const {getURLsFromHTML} = require('./crawl.js');
const {sortPages, printReport} = require('./report.js');

// tests for the normalizeURL function
test('normalizeURL: normal URL + normal path', () => {
  expect(normalizeURL('https://this.is.a.test.com/tester')).toBe(
    'https://this.is.a.test.com/tester',
  );
});

test('normalizeURL: normal URL + normal path + trailing /', () => {
  expect(normalizeURL('https://this.is.a.test.com/tester/')).toBe(
    'https://this.is.a.test.com/tester',
  );
});

test('normalizeURL: normal URL + normal path + query', () => {
  expect(
    normalizeURL('https://this.is.a.test.com/tester?tester=a&newtester=b'),
  ).toBe('https://this.is.a.test.com/tester');
});

test('normalizeURL: normal URL + diff port + normal path', () => {
  expect(normalizeURL('https://this.is.a.test.com:9191/tester')).toBe(
    'https://this.is.a.test.com:9191/tester',
  );
});

test('normalizeURL: normal URL + no path/query', () => {
  expect(normalizeURL('https://this.is.a.test.com')).toBe(
    'https://this.is.a.test.com',
  );
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

// tests for the sortPages function
test('sortPages: test out a random sample of object', () => {
  const testObject = {
    one: 10,
    two: 3,
    three: 5,
    four: 14,
  };

  const expectedObject = {
    two: 3,
    three: 5,
    one: 10,
    four: 14,
  };

  expect(sortPages(testObject)).toStrictEqual(expectedObject);
});

test('sortPages: test out an already sorted object', () => {
  const testObject = {
    two: 3,
    three: 5,
    one: 10,
    four: 14,
  };

  const expectedObject = {
    two: 3,
    three: 5,
    one: 10,
    four: 14,
  };

  expect(sortPages(testObject)).toStrictEqual(expectedObject);
});

test('sortPages: test out a null/empty object', () => {
  const testObject = {};
  const expectedObject = {};

  expect(sortPages(testObject)).toStrictEqual(expectedObject);
});

// test for the printReport function
test('printReport: check that the function works with a non-empty object', () => {
  const testObject = {
    two: 3,
    three: 5,
    one: 10,
    four: 14,
  };

  expect(printReport(testObject)).toBeUndefined();
});
