function sortPages(pages) {
  let sortedPages = {};

  let sortable = [];
  for (var page in pages) {
    sortable.push([page, pages[page]]);
  }

  sortable.sort(function (a, b) {
    return a[1] - b[1];
  });

  sortable.forEach((item) => {
    sortedPages[item[0]] = item[1];
  });

  return sortedPages;
}

function printReport(pages) {
  console.log('Report starting...');
  sortedPages = sortPages(pages);

  const pageset = Object.keys(sortedPages);
  pageset.forEach((page) => {
    console.log(`Found ${sortedPages[page]} internal links to ${page}`);
  });
}

module.exports = {
  printReport,
  sortPages,
};
