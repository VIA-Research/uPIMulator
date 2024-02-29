import graphviz
import os


if __name__ == "__main__":
    script_dirpath = os.path.dirname(__file__)

    bin_dirpath = os.path.join(script_dirpath, "..", "bin")
    host_dirpath = os.path.join(bin_dirpath, "host")

    nodes_filepath = os.path.join(host_dirpath, "ast_nodes.txt")
    edges_filepath = os.path.join(host_dirpath, "ast_edges.txt")

    g = graphviz.Digraph(filename="ast", format="pdf", directory=host_dirpath)

    with open(nodes_filepath, "r") as f:
        for line in f.readlines():
            node = line.strip()

            g.node(name=node)

    with open(edges_filepath, "r") as f:
        for line in f.readlines():
            tail, head = line.split("->")

            tail = tail.strip()
            head = head.strip()

            g.edge(tail_name=tail, head_name=head)

    g.view()
