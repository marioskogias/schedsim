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
            avg = float(l.split(" ")[3])

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
        res[q] = parse_file("mm1_ts_{}.dat".format(q))
    return res

def main():

    service_time = 50

    # plot run to completion
    rtc_data = parse_file("mm1_rtc.dat")
    x, y1, y2 = zip(*rtc_data)
    y1 = map(lambda a: a/float(service_time), y1)
    y2 = map(lambda a: a/float(service_time), y2)
    #plt.plot(x, y1, label="RTC average")
    plt.plot(x, y2, label="RTC 99th")

    # plot processor sharing
    rtc_data = parse_file("mm1_ps.dat")
    x, y1, y2 = zip(*rtc_data)
    y1 = map(lambda a: a/float(service_time), y1)
    y2 = map(lambda a: a/float(service_time), y2)
    #plt.plot(x, y1, label="PS average")
    plt.plot(x, y2, label="PS 99th")

    '''
    # plot time sharing
    data = parse_ps()
    #for q, v in data.iteritems():
    #    x, _, y2 = zip(*v)
    #    y2 = map(lambda a: a/float(service_time), y2)
    #    plt.plot(x, y2, label="PS 99th q={}*service_time".format(q/float(service_time)))

    # quantum
    q = 10
    v = data[q]
    x, _, y2 = zip(*v)
    y2 = map(lambda a: a/float(service_time), y2)
    plt.plot(x, y2, label="TS 99th q={}*service_time".format(q/float(service_time)))

    q = 190
    v = data[q]
    x, _, y2 = zip(*v)
    y2 = map(lambda a: a/float(service_time), y2)
    plt.plot(x, y2, label="TS 99th q={}*service_time".format(q/float(service_time)))
    '''

    axes = plt.gca()
    axes.set_ylim([0,100])
    plt.xlabel("rho")
    plt.ylabel("lantency normalized to service_time")
    plt.legend()
    plt.show()


if __name__ == "__main__":
    main()
