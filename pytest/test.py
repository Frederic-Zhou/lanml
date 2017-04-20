#encoding: utf-8
import numpy
import matplotlib.pyplot as plt

def load_data():
    data = []
    label = []
    file_object = open('test.txt')
    for line in file_object.readlines():
        arr = line.strip().split()
        data.append([1.0, float(arr[0]), float(arr[1])])
        label.append(int(arr[2]))


    return data, label

def sigmoid(x):
    return 1.0 / (1 + numpy.exp(-x))

def grad_ascent(data, label):
    data_mat = numpy.mat(data)
    label_mat = numpy.mat(label).transpose()

    # 获得矩阵列数
    n, m = data_mat.shape

    alpha = 0.001
    max_step = 500
    weights = numpy.ones((m, 1))

    
    for i in range(max_step):
        h = sigmoid(data_mat * weights)
        err = (label_mat - h)
        weights = weights + alpha * data_mat.transpose() * err
    return weights

#随机梯度上升
def stoc_grad_ascent0(data, label):
    data = numpy.array(data)
    n, m = data.shape
    alpha = 0.01
    weights = numpy.ones(m)
    for i in range(n):
        h = sigmoid(numpy.sum(data[i] * weights))
        err = label[i] - h
        weights = weights + data[i] * alpha * err
    return weights

#优化后的随机梯度上升
def stoc_grad_ascent1(data, label, max_step):
    data = numpy.array(data)
    n, m = data.shape
    weights = numpy.ones(m)
    for i in range(max_step):
        data_index = range(n)
        for j in range(n):
            alpha = 4 / (1.0 + i + j) + 0.01
            rand_index = int(numpy.random.uniform(0, len(data_index)))
            h = sigmoid(numpy.sum(data[rand_index] * weights))
            error = label[rand_index] - h
            weights = weights + alpha * error * data[rand_index]
            del(data_index[rand_index])
    return weights

def plot(weights):
    data, label = load_data()
    n = len(data)
    x1 = []
    y1 = []
    x2 = []
    y2 = []
    for i in range(n):
        if label[i] == 1:
            x1.append(data[i][1])
            y1.append(data[i][2])
        else:
            x2.append(data[i][1])
            y2.append(data[i][2])
    fig = plt.figure()
    ax = fig.add_subplot(111)
    ax.scatter(x1, y1, s = 30, c = 'red', marker = 's')
    ax.scatter(x2, y2, s = 30, c = 'green')
    x = numpy.arange(-3.0, 3.0, 0.1)
    y = (-weights[0] - weights[1] * x) / weights[2]
    ax.plot(x, y)
    plt.xlabel('X1')
    plt.ylabel('X2')
    plt.show()

if __name__=="__main__":
    data, label = load_data()

    weights = grad_ascent(data, label)
    # print weights
    plot(weights.getA())

    # weights = stoc_grad_ascent0(data, label)
    # # print weights
    # plot(weights)

    # weights = stoc_grad_ascent1(data, label, 500)
    # # print weights
    # plot(weights)