

def f(x1, x2, x3, x4):
    return 0.6224 * x1 * x3 * x4 + 1.7781 * x2 * x3**2 + 3.1661 * x4 * x1**2 + 19.84 * x3 * x1**2


print(f(0.8155, 0.4478, 41.1892, 176.6589))


def f2(x1, x2, x3, x4, x5):
    return 0.6224 * (x1 + x2 + x3 + x4 + x5)


print('MFO', f2(5.9830, 5.3167, 4.4973, 3.5136, 2.1616))
print('SMA', f2(6.017757, 5.310892, 4.493758, 3.501106, 2.150159))
print('GCA', f2(6.0100,5.3000,4.4900,3.4900,2.1500))
print('SSA', f2(6.015134526133134, 5.309304676055819, 4.495006716308508, 3.501426286300545, 2.152787908005768))
print('PRO', f2(6.0151, 5.3025, 4.5392, 3.4914, 2.1690))


def f3(x1, x2, x3, x4):
    a = 5000
    b = 1/12*x3*(x1-2*x4)**3+1/6*x2*x4**3+2*x2*x4*((x1-x4)/2)**2

    return a / b


print('SMA', f3(79.9948,50.4200,0.9000,2.3205))
