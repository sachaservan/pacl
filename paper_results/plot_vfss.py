from email.mime import base
from plot_common import *
from mpl_toolkits.axes_grid1.inset_locator import zoomed_inset_axes
from mpl_toolkits.axes_grid1.inset_locator import mark_inset


def plot(ax, axins, num_keys, num_subkeys, group_exp, baseline, pacl_pk, amortize):
    ######################## PLOT CODE ########################
    ax.yaxis.grid(color=gridcolor, linestyle=linestyle, linewidth=0.5)
    ax.xaxis.grid(color=gridcolor, linestyle=linestyle, linewidth=0.5)
    ax.set_xticks(xticks)
    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)

    groupNumber = 0
    groupSize = int(len(num_keys) / len(np.unique(num_subkeys)))
    for subkeys in range(0, len(num_keys), groupSize):
        start = groupSize*groupNumber
        stop = start + groupSize

        sort = np.argsort(num_keys[start:stop])

        num_keys[start:stop] = num_keys[start:stop][sort]
        baseline[start:stop] = baseline[start:stop][sort]
        pacl_pk[start:stop] = pacl_pk[start:stop][sort]

        div = np.ones(len(num_keys[start:stop]))
        if amortize:
            div = num_keys[start:stop]

        print(groupNumber)
        print(num_keys[start:stop])
        print(pacl_pk[start:stop][:, 0]/div)

        # only plot one baseline
        if groupNumber == 0:
            ax.plot(
                num_keys[start:stop],
                baseline[start:stop][:, 0]/div,
                linestyle=linestyles[0],
                color=colors[0],
                marker=markers[0],
                lw=linewidth,
                label=str("VFSS (baseline)")
            )

            # plot exponentiation time
            group_exp[start:stop] = group_exp[start:stop][sort]
            ax.plot(
                num_keys[start:stop],
                group_exp[start:stop]/div,
                linestyle='--',
                color='red',
                lw=1,
                zorder=10,
            )

            axins.plot(
                num_keys[start:stop],
                baseline[start:stop][:, 0]/div,
                linestyle=linestyles[0],
                color=colors[0],
                marker=markers[0],
                lw=linewidth)

            axins.plot(
                num_keys[start:stop],
                group_exp[start:stop]/div,
                linestyle='--',
                color='red',
                lw=1,
                zorder=10)

        ax.plot(
            num_keys[start:stop],
            pacl_pk[start:stop][:, 0]/div,
            linestyle=linestyles[groupNumber+1],
            color=colors[groupNumber+1],
            marker=markers[groupNumber+1],
            lw=linewidth,
            label=str("w/ PACL ") + r"($\ell =$" +
            str(num_subkeys[subkeys]) + ")"
        )

        axins.plot(
            num_keys[start:stop],
            pacl_pk[start:stop][:, 0]/div,
            linestyle=linestyles[groupNumber+1],
            color=colors[groupNumber+1],
            marker=markers[groupNumber+1],
            lw=linewidth
        )

        ax.fill_between(
            num_keys[start:stop],
            (pacl_pk[start:stop][:, 0] -
             pacl_pk[start:stop][:, 1])/div,
            (pacl_pk[start:stop][:, 0] +
             pacl_pk[start:stop][:, 1])/div,
            color=colors[groupNumber+1],
            alpha=error_opacity
        )

        groupNumber += 1

    return ax


if __name__ == '__main__':
    argparser = argparse.ArgumentParser(sys.argv[0])
    argparser.add_argument("--file", type=str, default='')

    args = argparser.parse_args()

    # read experiment file (expected json)
    with open(args.file, 'r') as myfile:
        data = myfile.read()

    # parse the experiment file as json
    results = json.loads(data)

    num_keys = []
    num_subkeys = []

    group_exp_processing = []

    eq_baseline_processing = []
    eq_vdpf_pacl_processing = []
    eq_vdpf_sk_pacl_processing = []

    range_baseline_processing = []
    range_vdpf_pacl_processing = []
    range_vdpf_sk_pacl_processing = []

    num_results = 0
    num_trials = len(results[0]["equality_baseline_processing_us"])

    # first we extract the relevant bits
    for i in range(len(results)):

        num_keys.append(results[i]["num_keys"])
        num_subkeys.append(results[i]["num_subkeys"])
        group_exp_processing.append(results[i]["group_exp_us"])

        avg = np.mean(results[i]["equality_baseline_ver_processing_us"])
        std = np.std(results[i]["equality_baseline_ver_processing_us"])
        eq_baseline_processing.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["equality_vdpf_pacl_processing_us"])
        std = np.std(results[i]["equality_vdpf_pacl_processing_us"])
        eq_vdpf_pacl_processing.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["equality_vdpf_sk_pacl_processing_us"])
        std = np.std(results[i]["equality_vdpf_sk_pacl_processing_us"])
        eq_vdpf_sk_pacl_processing.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["range_baseline_ver_processing_us"])
        std = np.std(results[i]["range_baseline_ver_processing_us"])
        range_baseline_processing.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["range_vdpf_pacl_processing_us"])
        std = np.std(results[i]["range_vdpf_pacl_processing_us"])
        range_vdpf_pacl_processing.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["range_vdpf_sk_pacl_processing_us"])
        std = np.std(results[i]["range_vdpf_sk_pacl_processing_us"])
        range_vdpf_sk_pacl_processing.append(
            [avg, confidence95(std, num_trials)])

    num_results += 1

    num_keys = np.array(num_keys)
    num_subkeys = np.array(num_subkeys)
    group_exp_processing = np.array(group_exp_processing)
    eq_baseline_processing = np.array(eq_baseline_processing)
    eq_vdpf_pacl_processing = np.array(eq_vdpf_pacl_processing)
    eq_vdpf_sk_pacl_processing = np.array(eq_vdpf_sk_pacl_processing)
    range_baseline_processing = np.array(range_baseline_processing)
    range_vdpf_pacl_processing = np.array(range_vdpf_pacl_processing)
    range_vdpf_sk_pacl_processing = np.array(range_vdpf_sk_pacl_processing)

    print("Num datapoints: " + str(len(num_keys)))
    sort = np.argsort(num_subkeys)
    num_keys = num_keys[sort]
    num_subkeys = num_subkeys[sort]
    group_exp_processing = group_exp_processing[sort]
    eq_baseline_processing = eq_baseline_processing[sort]
    eq_vdpf_pacl_processing = eq_vdpf_pacl_processing[sort]
    eq_vdpf_sk_pacl_processing = eq_vdpf_sk_pacl_processing[sort]
    range_baseline_processing = range_baseline_processing[sort]
    range_vdpf_pacl_processing = range_vdpf_pacl_processing[sort]
    range_vdpf_sk_pacl_processing = range_vdpf_sk_pacl_processing[sort]

    xticks = [1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024]
    indices = np.nonzero(np.in1d(num_keys, xticks,))

    # PLOT SERVER CPU TIME

    fig, (ax1, ax2) = plt.subplots(1, 2)
    fig.set_size_inches(width, 1.75)

    # make a little zoomed-in window for the tail end (zoom level = 35)
    axins1 = zoomed_inset_axes(ax1, 35, loc='center right')
    axins2 = zoomed_inset_axes(ax2, 35, loc='center right')

    amortize = False
    plot(ax1, axins1,
         num_keys[indices],
         num_subkeys[indices],
         group_exp_processing[indices],
         eq_baseline_processing[indices],
         eq_vdpf_pacl_processing[indices], amortize)
    ax1.set_xticks(xticks)
    ax1.set_xscale('log', base=2)
    # ax1.set_yscale('log', base=10)
    ax1.set_title("VDPF-PACL")
    ax1.set_xlabel('Number of evaluations')
    ax1.set_ylabel("Amortized CPU time ($\mu$s)")

    plot(ax2, axins2,
         num_keys[indices],
         num_subkeys[indices],
         group_exp_processing[indices],
         range_baseline_processing[indices],
         range_vdpf_pacl_processing[indices], amortize)
    ax2.set_xticks(xticks)
    ax2.set_xscale('log', base=2)
    ax2.set_title("VDMPF-PACL")
    ax2.set_xlabel('Number of evaluations')

    # make zoomed-in window

    # x-range for zoomed region
    x1 = 490
    x2 = 540

    # y-range for zoomed region
    y1 = -10
    y2 = 150

    # make the zoom-in plot:
    # https://stackoverflow.com/questions/8938449/how-to-extract-data-from-matplotlib-plot
    axins1.set_xlim(x1, x2)
    axins1.set_ylim(y1, y2)
    axins2.set_xlim(x1, x2)
    axins2.set_ylim(y1, y2)
    axins1.set_xticks([])
    axins1.set_yticks([])
    axins2.set_xticks([])
    axins2.set_yticks([])
    mark_inset(ax1, axins1, loc1=3, loc2=4,
               linewidth=0.75, fc="gray", ec="gray")
    mark_inset(ax2, axins2, loc1=3, loc2=4,
               linewidth=0.75, fc="gray", ec="gray")

    # make ghost plot for extra legend

    axghost = ax1.twinx()
    axghost.plot(np.NaN, np.NaN,
                 label="Exponentiation", color='red', linestyle="--", linewidth=1)
    axghost.get_yaxis().set_visible(False)
    axghost.spines['top'].set_visible(False)
    axghost.spines['right'].set_visible(False)
    axghost.legend(loc='upper right', edgecolor='white',
                   framealpha=1, fancybox=False)

    axghost = ax2.twinx()
    axghost.plot(np.NaN, np.NaN,
                 label="Exponentiation", color='red', linestyle="--", linewidth=1)
    axghost.get_yaxis().set_visible(False)
    axghost.spines['top'].set_visible(False)
    axghost.spines['right'].set_visible(False)
    axghost.legend(loc='upper right', edgecolor='white',
                   framealpha=1, fancybox=False)

    fig.tight_layout()

    ax1.legend(ncol=2, loc='center', framealpha=1, bbox_to_anchor=(
        1.1, -0.65), facecolor='white', fancybox=False,  handlelength=handlelength)
    fig.savefig('plot_vfss.pdf', bbox_inches='tight')
