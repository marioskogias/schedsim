import numpy as np
import matplotlib.pyplot as plt

def parse_file(fname):
    res = []
    with open(fname, 'r') as f:
        l = f.readline()
        while l:
            # get rho
            rho = float(l.split(" ")[2])

            # get avg
            l = f.readline()
            avg = l.split(" ")[3]

            # get 99th percentile
            l = f.readline()
            p = float(l.split("\t")[3].split(" ")[1])

            l = f.readline()
            res.append((rho, avg, p))
    return res

def parse_ps():
    quantums = range(10, 200, 20)
    res = {}
    for q in quantums:
        res[q] = parse_file("mm1_ps_{}.dat".format(q))
    return res

def main():

    service_time = 50
    # plot run to completion
    rtc_data = parse_file("mm1_rtc.dat")
    x, y1, y2 = zip(*rtc_data)
    #plt.plot(x, y1, label="RTC average")
    plt.plot(x, y2, label="RTC 99th")

    # plot processor sharing
    data = parse_ps()
    #for q, v in data.iteritems():
    #    x, _, y2 = zip(*v)
    #    plt.plot(x, y2, label="PS 99th q={}*service_time".format(q/float(service_time)))

    # quantum
    q = 10
    v = data[q]
    x, _, y2 = zip(*v)
    plt.plot(x, y2, label="PS 99th q={}*service_time".format(q/float(service_time)))

    q = 190
    v = data[q]
    x, _, y2 = zip(*v)
    plt.plot(x, y2, label="PS 99th q={}*service_time".format(q/float(service_time)))

    plt.xlabel("rho")
    plt.ylabel("latency")
    plt.legend()
    plt.show()


if __name__ == "__main__":
    main()
