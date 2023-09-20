const {test, expect} = require('@jest/globals');
const {normalizeURL} = require('./crawl.js');

test('normal URL + normal path', () => {
  expect(normalizeURL('https://this.is.a.test.com/tester')).toBe(
    'this.is.a.test.com/tester',
  );
});

test('normal URL + normal path + trailing /', () => {
  expect(normalizeURL('https://this.is.a.test.com/tester/')).toBe(
    'this.is.a.test.com/tester',
  );
});

test('normal URL + normal path + query', () => {
  expect(
    normalizeURL('https://this.is.a.test.com/tester?tester=a&newtester=b'),
  ).toBe('this.is.a.test.com/tester');
});

test('normal URL + diff port + normal path', () => {
  expect(normalizeURL('https://this.is.a.test.com:9191/tester')).toBe(
    'this.is.a.test.com:9191/tester',
  );
});

test('normal URL + no path/query', () => {
  expect(normalizeURL('https://this.is.a.test.com')).toBe('this.is.a.test.com');
});
