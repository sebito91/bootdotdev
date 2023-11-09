"""Module to render individual cells."""

from window.window import Point, Line, Window


class Cell:
    """Define an individual cell on the canvas."""

    def __init__(self, p1, p2, window=None):
        """Initialize a given cell with the Point and Window objects."""
        self.has_left_wall = True
        self.has_right_wall = True
        self.has_top_wall = True
        self.has_bottom_wall = True
        self.visited = False

        self._left_wall_color = "black"
        self._right_wall_color = "black"
        self._top_wall_color = "black"
        self._bottom_wall_color = "black"

        self._x1 = p1.x
        self._x2 = p2.x
        self._y1 = p1.y
        self._y2 = p2.y

        self._p1 = p1
        self._p2 = p2

        self._win = window

    @property
    def x1(self):
        return self._x1

    @property
    def y1(self):
        return self._y1

    @property
    def x2(self):
        return self._x2

    @property
    def y2(self):
        return self._y2

    @property
    def p1(self):
        return self._p1

    @property
    def p2(self):
        return self._p2

    def draw(self, top_left, bottom_right):
        """Draw out the cell using the top-left Point and bottom-right Point given."""
        if not self._win:
            return

        left_wall = Line(top_left, Point(top_left.x, bottom_right.y))
        right_wall = Line(Point(bottom_right.x, top_left.y), bottom_right)
        top_wall = Line(top_left, Point(bottom_right.x, top_left.y))
        bottom_wall = Line(Point(top_left.x, bottom_right.y), bottom_right)

        if not self.has_left_wall:
            self._left_wall_color = "white"

        if not self.has_right_wall:
            self._right_wall_color = "white"

        if not self.has_top_wall:
            self._top_wall_color = "white"

        if not self.has_bottom_wall:
            self._bottom_wall_color = "white"

        self._win.draw_line(left_wall, self._left_wall_color)
        self._win.draw_line(right_wall, self._right_wall_color)
        self._win.draw_line(top_wall, self._top_wall_color)
        self._win.draw_line(bottom_wall, self._bottom_wall_color)

    def draw_move(self, to_cell, undo=False):
        """Draw out a move from one cell to another."""
        color = "red"
        if undo:
            color = "gray"

        line_start = Point((self.x1 + (abs(self.x1 - self.x2) // 2)), (self.y1 + (abs(self.y1 - self.y2) // 2)))
        line_end = Point((to_cell.x1 + (abs(to_cell.x1 - to_cell.x2) // 2)), (to_cell.y1 + (abs(to_cell.y1 - to_cell.y2) // 2)))

        self._win.draw_line(Line(line_start, line_end), color)

    def __str__(self):
        """Print out the cell object signature."""
        return f"cell is at top-left ({self.x1}, {self.y1}) and bottom-right ({self.x2}, {self.y2})"
