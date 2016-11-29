import numpy as np
import matplotlib.pyplot as plt

def parse_file(fname):
    res = []
    with open(fname, 'r') as f:
        l = f.readline()
        while l:
            # get mu
            mu = float(l.split("\t")[2].split("=")[1])

            # get avg
            l = f.readline()
            avg = float(l.split(" ")[3])

            # get percentiles
            l = f.readline()
            tmp = l.split("\t")
            p50 = float(tmp[0].split(" ")[1])
            p90 = float(tmp[1].split(" ")[1])
            p95 = float(tmp[2].split(" ")[1])
            p99 = float(tmp[3].split(" ")[1])

            # get achieved requests
            l = f.readline()
            qps = float(l.split(":")[1])

            # compute rho
            res.append((qps, mu, avg, p50, p90, p95, p99))

            l = f.readline()
    return res

def plot_data(data, name, p):
    qps, mu, avg, p50, p90, p95, p99 = zip(*data)

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
    y = map(lambda a,b: a*b, to_plot, mu)
    x = map(lambda a,b: a/b, qps, mu)
    plt.plot(x, y, label="{} {}".format(name, p))

def main():

    # Plotting goes here
    data = parse_file("data/multicore_s0_q500.dat")
    plot_data(data,"0 slow", 99)

    # plot horizontal line at 1
    plt.axhline(y=1)

    axes = plt.gca()
    # set log y axis
    #axes.set_yscale("log")
    # y axis limit
    #axes.set_ylim([0,1000])

    axes.get_xaxis().set_ticks([])

    plt.xlabel("RPS")
    plt.ylabel("lantency normalized to service_time")
    plt.legend()
    plt.show()

if __name__ == "__main__":
    main()
