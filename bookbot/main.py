#!/usr/bin/env python3
"""Module to render a book bot."""


from pathlib import Path
from collections import defaultdict


def read_book(book: Path = None):
    """Read the given book book.

    :param book: Path - path to the book we want to read
    """
    with book.open() as f:
        file_contents = f.read()

    return file_contents


def count_words(content: str) -> int:
    """Count the number of words in the given content.

    :param content: str - input string to count content
    :return: int - count of words from content
    """
    return len(content.split())


def parse_letters(content: str) -> dict:
    """Parse out the individual lowercase letters from the content.

    :param content: str - input string to parse content
    :return: dict - dict of lowercase counts of letters from content
    """
    letters = defaultdict(int)
    for word in content.split():
        for char in word:
            letters[char.lower()] += 1

    return dict(sorted(letters.items(), key=lambda x: x[1], reverse=True))


def generate_report(book: str, word_count: int, letter_count: dict[str, int]) -> str:
    """Generate a report from the counts of the book.

    :param book: str - path to book content
    :param word_count: int - number of words from the book
    :param letter_count: dict[str, int] - dict of characters and their count from the book's contents
    :return: str - generated report of the given document
    """
    output = [f"--- Begin report of {book} ---"]
    output.append(f"{word_count} words found in the document")
    output.append("")

    output += [f"The '{char}' was found {count} times" for char, count in letter_count.items() if char.isalpha()]

    output.append("--- End report ---")
    return "\n".join(output)


if __name__ == "__main__":
    book = "books/frankenstein.txt"
    title = Path("/home/sborza/src/github.com/sebito91/bookbot/" + book)

    file_contents = read_book(title)
    word_count = count_words(file_contents)
    letter_count = parse_letters(file_contents)

    print(generate_report(book, word_count, letter_count))
