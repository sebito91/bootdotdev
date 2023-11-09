"""Module to store off the maze itself."""

import random
import time

from window.cell import Cell
from window.window import Point, Window


class Maze:
    """Object to store the maze window and attributes."""

    def __init__(self, x1, y1, num_rows, num_cols, cell_size_x, cell_size_y, win=None, seed=None):
        """Initialize the maze itself."""
        self._x1 = x1
        self._y1 = y1
        self._num_rows = num_rows
        self._num_cols = num_cols
        self._cell_size_x = cell_size_x
        self._cell_size_y = cell_size_y
        self._win = win

        if seed:
            random.seed(seed)
        else:
            random.seed(0)

        self._create_cells()

    def _create_cells(self):
        """Create the cells for the maze."""
        self._cells = [[self._draw_cell(col, row) for row in range(self._num_rows)] for col in range(self._num_cols)]
        self._break_entrance_and_exit()
        self._break_walls_r(0, 0)
        self._reset_cells_visited()

    def _draw_cell(self, col, row):
        """Draw out the individual cell given the column and row it represents."""
        top_left = Point(self._x1 + (col * self._cell_size_x), self._y1 + (row * self._cell_size_y))
        bottom_right = Point(self._x1 + ((col + 1) * self._cell_size_x), self._y1 + ((row + 1) * self._cell_size_y))

        cell = Cell(top_left, bottom_right, self._win)

        cell.draw(top_left, bottom_right)
        self._animate()

        return cell

    def _animate(self):
        """Animate the cell and pause for a period."""
        if self._win:
            self._win.redraw()
        time.sleep(0.05)

    def _break_entrance_and_exit(self):
        """Break up the entrance and exit Cell objects to allow for the maze ingress and egress."""
        self._cells[0][0].has_top_wall = False
        self._cells[self._num_cols - 1][self._num_rows - 1].has_bottom_wall = False

        entrance_top_left = Point(self._x1, self._y1)
        entrance_bottom_right = Point(self._x1 + self._cell_size_x, self._y1 + self._cell_size_y)

        exit_top_left = Point(self._x1 + ((self._num_cols - 1) * self._cell_size_x), self._y1 + ((self._num_rows - 1) * self._cell_size_y))
        exit_bottom_right = Point(self._x1 + (self._num_cols * self._cell_size_x), self._y1 + (self._num_rows * self._cell_size_y))

        self._cells[0][0].draw(entrance_top_left, entrance_bottom_right)
        self._cells[self._num_cols - 1][self._num_rows - 1].draw(exit_top_left, exit_bottom_right)

    def _break_walls_r(self, col, row):
        """Break down walls throughout the maze using a breadth-first recursive search."""
        if not 0 <= col < self._num_cols or not 0 <= row < self._num_rows:
            return

        self._cells[col][row].visited = True
        while True:
            to_visit = []

            if col - 1 >= 0 and not self._cells[col - 1][row].visited:
                to_visit.append((col - 1, row))

            if col + 1 < self._num_cols and not self._cells[col + 1][row].visited:
                to_visit.append((col + 1, row))

            if row - 1 >= 0 and not self._cells[col][row - 1].visited:
                to_visit.append((col, row - 1))

            if row + 1 < self._num_rows and not self._cells[col][row + 1].visited:
                to_visit.append((col, row + 1))

            if len(to_visit) == 0:
                self._cells[col][row].draw(self._cells[col][row].p1, self._cells[col][row].p2)
                return

            col_to_visit, row_to_visit = to_visit.pop(random.randrange(len(to_visit)))

            if col_to_visit == col and row_to_visit < row:
                self._cells[col][row].has_top_wall = False
                self._cells[col][row_to_visit].has_bottom_wall = False
            elif col_to_visit == col and row_to_visit > row:
                self._cells[col][row].has_bottom_wall = False
                self._cells[col][row_to_visit].has_top_wall = False
            elif col_to_visit < col and row_to_visit == row:
                self._cells[col][row].has_left_wall = False
                self._cells[col_to_visit][row].has_right_wall = False
            elif col_to_visit > col and row_to_visit == row:
                self._cells[col][row].has_right_wall = False
                self._cells[col_to_visit][row].has_left_wall = False

            self._break_walls_r(col_to_visit, row_to_visit)

    def _reset_cells_visited(self):
        """Reset all cells to not visited."""
        for col in range(self._num_cols):
            for row in range(self._num_rows):
                self._cells[col][row].visited = False

    def solve(self):
        """Solve the maze."""
        return self._solve_r(0, 0)

    def _solve_r(self, col, row):
        """Solve the maze starting from the given column and row entry point."""
        self._animate()

        self._cells[col][row].visited = True

        # check if we're at the end, and if so let's return success!
        if (col == self._num_cols - 1) and (row == self._num_rows - 1):
            return True

        for next_col, next_row, direction in [(col, row - 1, "up"), (col, row + 1, "down"), (col - 1, row, "left"), (col + 1, row, "right")]:
            if not (0 <= next_col < self._num_cols and 0 <= next_row < self._num_rows):
                continue

            if self._cells[next_col][next_row].visited:
                continue

            if direction == "up" and self._cells[next_col][next_row].has_bottom_wall:
                continue

            if direction == "down" and self._cells[next_col][next_row].has_top_wall:
                continue

            if direction == "left" and self._cells[next_col][next_row].has_right_wall:
                continue

            if direction == "right" and self._cells[next_col][next_row].has_left_wall:
                continue

            self._cells[col][row].draw_move(self._cells[next_col][next_row])
            if self._solve_r(next_col, next_row):
                return True

            self._cells[col][row].draw_move(self._cells[next_col][next_row], undo=True)

        return False
