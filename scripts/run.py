#!/usr/bin/python

import numpy as np
from subprocess import call
from scipy.stats import expon
from math import log, sqrt

mm116_lambdas = np.arange(0.1, 1.7,0.1)

def single_queue():
    with open("data/single_queue.dat", 'w') as f:
        [call(["schedsim","--lambda={}".format(l), "--mu=0.1", "--duration=1000000"], stdout=f) for l in mm116_lambdas]

def main():
    single_queue()

if __name__ == "__main__":
    main()
