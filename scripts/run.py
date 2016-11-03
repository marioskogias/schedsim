#!/usr/bin/python

import numpy as np
from subprocess import call

lambdas = np.arange(0.001,0.02, 0.0015)

def run_rtc():
    with open("mm1_rtc.dat", 'w') as f:
        [call(["schedsim","--lambda={}".format(x)], stdout=f) for x in lambdas]

def run_ps():
    quantums = range(10, 200, 20)
    for q in quantums:
        with open("mm1_ts_{}.dat".format(q), 'w') as f:
            [call(["schedsim","--lambda={}".format(x), "--system=ts",
                "--quantum={}".format(q)], stdout=f) for x in lambdas]

def main():
    run_ps()

if __name__ == "__main__":
    main()
