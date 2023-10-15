"""Module to represent a given window for the maze."""

from tkinter import Tk, Canvas, CENTER


class Window:
    """Widget to render a window for the maze object."""

    def __init__(self, width, height):
        """Initialize the maze window object."""
        self.__root = Tk()
        self.__window_running = False

        self.__root.title("Maze Runner in all its glory!")
        self.__root.protocol("WM_DELETE_WINDOW", self.close)

        self.__canvas = Canvas(self.__root, {"height": height, "width": width, "bg": "white"})
        self.__canvas.pack()

    def redraw(self):
        """Redraw the canvas object."""
        self.__root.update_idletasks()
        self.__root.update()

    def draw_line(self, line, fill_color):
        """Draw out the line."""
        line.draw(self.__canvas, fill_color)

    def wait_for_close(self):
        """Keep the window alive and redraw until close received."""
        self.__window_running = True
        while self.__window_running:
            self.redraw()

    def close(self):
        """Close out the running window."""
        self.__window_running = False


class Point:
    """Store the x- and y-coordinate for a given point."""

    def __init__(self, x, y):
        """Initialize the Point object to store the x- and y-coordinate."""
        self.__x = x
        self.__y = y

    @property
    def x(self):
        """Return the x-coordinate of the Point."""
        return self.__x

    @property
    def y(self):
        """Return the y-coordinate of the Point."""
        return self.__y


class Line:
    """Render a line within our Window object."""

    def __init__(self, p1, p2):
        """Initialize the Line object to render a line between two points."""
        self.__p1 = p1
        self.__p2 = p2

    def draw(self, canvas, fill_color):
        """Draw out the line on the given Canvas."""
        canvas.create_line(self.__p1.x, self.__p1.y, self.__p2.x, self.__p2.y, fill=fill_color, width=2)
        canvas.pack()
