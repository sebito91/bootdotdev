"""Module to run the Maze Runner program."""

import time

from window.maze import Maze
from window.window import Window


if __name__ == "__main__":
    new_maze = Window(800, 800)

    time_seed = int(time.time() * 1_000_000)
    new_new_maze = Maze(20, 20, 30, 30, 25, 25, new_maze, time_seed)

    if not new_new_maze.solve():
        print("could not solve the maze :(")
    else:
        print("solved the maze! :D")

    new_maze.wait_for_close()
