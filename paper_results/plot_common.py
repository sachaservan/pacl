import sys
import time
import math
import argparse
import numpy as np
import matplotlib
import matplotlib.pyplot as plt
from matplotlib.ticker import StrMethodFormatter, NullFormatter, FuncFormatter
import matplotlib.ticker as ticker
import json
import matplotlib.colors as mcolors


###########################
# PLOT PARAMETERS
###########################
TINY_SIZE = 7
SMALL_SIZE = 7
MEDIUM_SIZE = 8
BIGGER_SIZE = 10
matplotlib.rcParams['font.family'] = 'serif'
plt.rc('font', size=SMALL_SIZE)          # controls default text sizes
plt.rc('axes', titlesize=MEDIUM_SIZE)    # fontsize of the axes title
plt.rc('axes', labelsize=MEDIUM_SIZE)    # fontsize of the x and y labels
plt.rc('xtick', labelsize=TINY_SIZE)     # fontsize of the tick labels
plt.rc('ytick', labelsize=SMALL_SIZE)    # fontsize of the tick labels
plt.rc('legend', fontsize=SMALL_SIZE)    # legend fontsize
plt.rc('figure', titlesize=BIGGER_SIZE)  # fontsize of the figure title
plt.rc('hatch', linewidth=0.5)

colors = list(mcolors.TABLEAU_COLORS.values())
# colors = ['#1b9e77','#d95f02','#7570b3', '#e6ab02', '#e7298a','#66a61e']
# colors = ['#0099CC', '#008000', '#955196', '#003f5c', '#ff6e54', '#dd5182']
colors = ['#08519c', '#ff7f00', '#16a085', '#8e44ad', '#c0392b', '#333']
# colors = ['#af43be', '#fa8090', '#65dc98', '#ffa600', '#defe47']
bar_colors = ['#c6dbef', '#6baed6', '#08519c']
markers = ['x',  '.', '2', '+', '1']
# linestyles = ['-',  '--', '-.', ':']
linestyles = ['-',  '-', '-', '-', '-',  '-', '-', '-']
barhatches = ['////', '\\\\\\\\', 'xxxx', '//////////', '']

annotate_arrow = dict(arrowstyle="<->", color='black', linewidth=0.5,)
annoate_font_size = SMALL_SIZE
annoate_text_color = 'black'

handlelength = 3  # make legend ticks wider

linestyle = 'dashed'
gridcolor = '#ccc'
error_opacity = 0.15
linewidth = 1.15
edgecolor = 'white'
error_bars = dict(ecolor='black', lw=linewidth/2,
                  capsize=3, capthick=linewidth/2)

# PLOT DIM
width = 4
height = 2.75

MILLI_TO_SECONDS = 1.0 / 1000.0
SECONDS_TO_MILI = 1000.0
MILLI_TO_HOURS = 1.0 / 3600000
BYTES_TO_MB = 1.0 / (10**6)
BYTES_TO_KB = 1.0 / (10**3)
BYTES_TO_GB = 1.0 / (10**9)


def confidence95(x, n):
    return 1.96 * x / math.sqrt(n)


def autolabel(ax, rects, group_size):
    """Attach a text label above each bar in *rects*, displaying its height."""
    for r in range(len(rects), group_size):
        rect = rects[r]
        height = rect.get_height() + rect.get_xy()[1]
        if height < 60:
            continue

        ax.annotate('{}'.format(time.strftime('%M:%S', time.gmtime(height))),
                    xy=(rect.get_x() - rect.get_width(), height),
                    xytext=(0, 3),  # 3 points vertical offset
                    textcoords="offset points",
                    color="#333333",
                    ha='center', va='bottom')
