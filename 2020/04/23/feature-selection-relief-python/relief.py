from numpy import random, zeros, matmul, multiply
from math import floor
from sklearn.datasets import load_iris, load_breast_cancer


class Relief:
    def __init__(self, repeat_time, c, X, Y):
        self.repeat_time = repeat_time
        self.c = c
        self.X = X
        self.Y = Y
        # m 样本的个数
        # n 特征的维数
        self.m, self.n = self.X.shape
        self.map = dict()
        # 用于保存计算出权重
        self.weights = zeros(self.n)

        for i, y in enumerate(self.Y):
            if y not in self.map:
                self.map.setdefault(y, [])
            self.map[y].append(i)

    def get_important_weights(self) -> list:
        # print(self.weights, sum(self.weights))
        sum_weight = sum(self.weights)
        return [(i, w / sum_weight) for i, w in enumerate(self.weights) if w / sum_weight > self.c]

    def get_nearest_distance(self, x, indices) -> float:
        nearest_distance = float('inf')
        nearest_i = 0

        for i in indices:
            distance = matmul(x - self.X[i], (x - self.X[i]))
            if distance < nearest_distance:
                nearest_distance = distance
                nearest_i = i

        return multiply(x - self.X[nearest_i], x - self.X[nearest_i])

    def run(self):
        for i in range(self.repeat_time):
            select = floor(random.random() * self.m)
            x = self.X[select]
            y = self.Y[select]
            # 在相同类型中找距离最近的，直接返回欧式距离的平方
            distance_same = self.get_nearest_distance(x, [j for j in self.map[y] if j != select])

            # 找到不同的标签，这里只两类
            y_diff = [y_diff for y_diff in self.map.keys() if y_diff != y][0]
            distance_diff = self.get_nearest_distance(x, [j for j in self.map[y_diff]])

            self.weights = self.weights - distance_same / self.repeat_time + distance_diff / self.repeat_time


X, Y = load_iris(return_X_y=True)
r = Relief(20, .0, X, Y)
r.run()
weights = r.get_important_weights()

print('iris')
for i, w in weights:
    print(i, w)

X, Y = load_breast_cancer(return_X_y=True)
r = Relief(20, .0, X, Y)
r.run()
weights = r.get_important_weights()

print('breast_cancer')
for i, w in weights:
    print(i, w)
