#!/usr/bin/python

import numpy as np
from subprocess import call

service_time = 50
lambdas = np.arange(0.001,0.02, 0.0015)
costs = [0, 0.001, 0.01, 0.1, 0.25, 0.5]

# Run to completion
def run_mm_rtc():
    with open("data/mm1_rtc.dat", 'w') as f:
        [call(["schedsim","--lambda={}".format(x)], stdout=f) for x in lambdas]

def run_md_rtc():
    with open("data/md1_rtc.dat", 'w') as f:
        [call(["schedsim","--lambda={}".format(x), "--service=d"],
            stdout=f) for x in lambdas]

def run_mlg_rtc():
    with open("data/mlg1_rtc.dat", 'w') as f:
        [call(["schedsim","--lambda={}".format(x), "--service=lg"],
            stdout=f) for x in lambdas]

def run_mb_rtc():
    with open("data/mb1_rtc.dat", 'w') as f:
        [call(["schedsim","--lambda={}".format(x), "--service=b"],
            stdout=f) for x in lambdas]

def run_mm_rtc_costs():
    for c in costs:
        with open("data/mm1_rtc_c_{}.dat".format(c), 'w') as f:
            [call(["schedsim","--lambda={}".format(x),
            "--ctxCost={}".format(c*service_time)], stdout=f) for x in lambdas]

# processorType sharing
def run_mm_ps():
    with open("data/mm1_ps.dat", 'w') as f:
        [call(["schedsim","--lambda={}".format(x), "--processorType=ps"],
            stdout=f) for x in lambdas]

def run_md_ps():
    with open("data/md1_ps.dat", 'w') as f:
        [call(["schedsim","--lambda={}".format(x), "--processorType=ps", "--service=d"],
            stdout=f) for x in lambdas]

def run_mlg_ps():
    with open("data/mlg1_ps.dat", 'w') as f:
        [call(["schedsim","--lambda={}".format(x), "--processorType=ps", "--service=lg"],
            stdout=f) for x in lambdas]

def run_mb_ps():
    with open("data/mb1_ps.dat", 'w') as f:
        [call(["schedsim","--lambda={}".format(x), "--processorType=ps", "--service=b"],
            stdout=f) for x in lambdas]

# Time sharing
def run_md_ts():
    quantums = range(10, 200, 20)
    for q in quantums:
        with open("data/md1_ts_{}.dat".format(q), 'w') as f:
            [call(["schedsim","--lambda={}".format(x), "--processorType=ts",
                "--quantum={}".format(q), "--service=d"], stdout=f) for x in lambdas]

def run_mm_ts():
    quantums = range(10, 200, 20)
    for q in quantums:
        with open("data/mm1_ts_{}.dat".format(q), 'w') as f:
            [call(["schedsim","--lambda={}".format(x), "--processorType=ts",
                "--quantum={}".format(q)], stdout=f) for x in lambdas]

def run_mlg_ts():
    quantums = range(10, 200, 20)
    for q in quantums:
        with open("data/mlg1_ts_{}.dat".format(q), 'w') as f:
            [call(["schedsim","--lambda={}".format(x), "--processorType=ts",
                "--quantum={}".format(q), "--service=lg"], stdout=f) for x in lambdas]

def run_mb_ts():
    quantums = range(10, 200, 20)
    for q in quantums:
        with open("data/mb1_ts_{}.dat".format(q), 'w') as f:
            [call(["schedsim","--lambda={}".format(x), "--processorType=ts",
                "--quantum={}".format(q), "--service=b"], stdout=f) for x in lambdas]

def main():

    # Run to completion
    print "Run md1 rtc"
    run_md_rtc()
    print "Run mm1 rtc"
    run_mm_rtc()
    print "Run mlg1 rtc"
    run_mlg_rtc()
    print "Run mb1 rtc"
    run_mb_rtc()
    print "Rum mm rtc variable ctx costs"
    run_mm_rtc_costs()

    # Processor sharing
    print "Run md1 ps"
    run_md_ps()
    print "Run mm1 ps"
    run_mm_ps()
    print "Run mlg1 ps"
    run_mlg_ps()
    print "Run mb1 ps"
    run_mb_ps()

    # Time sharing
    print "Run md1 ts"
    run_md_ts()
    print "Run mm1 ts"
    run_mm_ts()
    print "Run mlg1 ts"
    run_mlg_ts()
    print "Run mb1 ts"
    run_mb_ts()

if __name__ == "__main__":
    main()
