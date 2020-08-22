from collections import deque

class TreeNode:
    def __init__(self, x):
        self.val = x
        self.left = None
        self.right = None

class Solution:
    def rightSideView(self, root: TreeNode) -> list:
       layer_dict = dict()
       max_depth = -1

       queue = deque([(root, 0)])

       while queue:
            node, depth = queue.popleft()
            if node is not None:
                max_depth = max(max_depth, depth)
                layer_dict[max_depth] = node.val

                queue.append((node.left, max_depth + 1))
                queue.append((node.right, max_depth + 1))

       return [layer_dict[depth] for depth in range(max_depth + 1)]




