import numpy as np
import matplotlib.pyplot as plt

def parse_file(fname):
    res = []
    with open(fname, 'r') as f:
        l = f.readline()
        while l:
            l = f.readline()
            data = l.split("\t")
            cores = int(data[0].split(":")[1])
            mu = float(data[1].split(":")[1])
            intr_lambda = float(data[0].split(":")[1])

            l = f.readline() # skip the collector name
            l = f.readline() # skip label
            l = f.readline()
            tmp = l.split("\t")

            avg = float(tmp[1])
            p50 = float(tmp[3])
            p90 = float(tmp[4])
            p95 = float(tmp[5])
            p99 = float(tmp[6])
            qps = float(tmp[7])

            # compute rho
            res.append((cores, mu, intr_lambda, qps, avg, p50, p90, p95, p99))

            l = f.readline()
    return res

def plot_data(data, name, p):
    cores, mu, intr_lambda, qps, avg, p50, p90, p95, p99 = zip(*data)

    if p == "avg":
        to_plot = avg
    elif p == 50:
        to_plot = p50
    elif p == 90:
        to_plot = p90
    elif p == 95:
        to_plot = p95
    elif p == 99:
        to_plot = p99

    # Assuming mu and cores are always the same
    cores = cores[0]
    mu = mu[0]

    y = map(lambda a: a*mu, to_plot)
    x = map(lambda a: a/(cores*mu), qps)
    plt.plot(x, y, label="{} {}".format(name, p), marker='+')

def main():

    # Plotting goes here
    data = parse_file("data/single_queue.dat")
    plot_data(data,"M/M/8", 99)

    # plot horizontal line at 1
    plt.axhline(y=1)

    axes = plt.gca()
    axes.set_ylim([0,10])


    plt.xlabel("RPS")
    plt.ylabel("lantency normalized to service_time")
    plt.legend()
    plt.show()

if __name__ == "__main__":
    main()
