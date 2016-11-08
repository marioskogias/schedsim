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

            '''
            # get 50th percentile
            l = f.readline()
            p = float(l.split("\t")[0].split(" ")[1])
            '''

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

    # plot mm1 rtc
    rtc_data = parse_file("mm1_rtc.dat")
    x, y1, y2 = zip(*rtc_data)
    y1 = map(lambda a: a/float(service_time), y1)
    y2 = map(lambda a: a/float(service_time), y2)
    #plt.plot(x, y1, label="MM RTC average")
    plt.plot(x, y2, label="MM RTC 99th")
    #plt.plot(x, y2, label="MM RTC 50th")

    # plot md1 rtc
    rtc_data = parse_file("md1_rtc.dat")
    x, y1, y2 = zip(*rtc_data)
    y1 = map(lambda a: a/float(service_time), y1)
    y2 = map(lambda a: a/float(service_time), y2)
    #plt.plot(x, y1, label="MD RTC average")
    plt.plot(x, y2, label="MD RTC 99th")
    #plt.plot(x, y2, label="MD RTC 50th")

    # plot mm1 ps
    rtc_data = parse_file("mm1_ps.dat")
    x, y1, y2 = zip(*rtc_data)
    y1 = map(lambda a: a/float(service_time), y1)
    y2 = map(lambda a: a/float(service_time), y2)
    #plt.plot(x, y1, label="MM PS average")
    plt.plot(x, y2, label="MM PS 99th")
    #plt.plot(x, y2, label="MM PS 50th")

    # plot md1 ps
    rtc_data = parse_file("md1_ps.dat")
    x, y1, y2 = zip(*rtc_data)
    y1 = map(lambda a: a/float(service_time), y1)
    y2 = map(lambda a: a/float(service_time), y2)
    #plt.plot(x, y1, label="MD PS average")
    plt.plot(x, y2, label="MD PS 99th")
    #plt.plot(x, y2, label="MD PS 50th")


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

    # plot horizontal line at 1
    plt.axhline(y=1)
    axes = plt.gca()
    axes.set_ylim([0,100])
    plt.xlabel("rho")
    plt.ylabel("lantency normalized to service_time")
    plt.legend()
    plt.show()


if __name__ == "__main__":
    main()
