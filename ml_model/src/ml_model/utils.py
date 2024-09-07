import os

def remove_graph(filepath):
    if os.path.isfile(filepath):
        os.remove(filepath)
