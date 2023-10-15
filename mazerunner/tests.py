"""Module to run unit tests against the maze builder and its various components."""

import unittest

from window.maze import Maze
from window.window import Window


class Tests(unittest.TestCase):
    def test_maze_create_cells(self):
        """Test out cell creation in the Maze object."""
        num_cols = 12
        num_rows = 10
        m1 = Maze(0, 0, num_rows, num_cols, 10, 10)

        # make sure the number of columns is what we expect
        self.assertEqual(len(m1._cells), num_cols)

        # make sure the number of rows in each column is what we expect
        self.assertEqual(len(m1._cells[0]), num_rows)


if __name__ == "__main__":
    unittest.main()
