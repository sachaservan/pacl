from email.mime import base
from plot_common import *


def plot(ax, num_keys, num_subkeys, baseline,  pacl_pk, amortize):
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

        div = np.ones(len(num_keys[start:stop]))
        if amortize:
            div = num_keys[start:stop]

        sort = np.argsort(num_keys[start:stop])

        num_keys[start:stop] = num_keys[start:stop][sort]
        baseline[start:stop] = baseline[start:stop][sort]
        pacl_pk[start:stop] = pacl_pk[start:stop][sort]
        print(num_keys[start:stop])

        # only plot one baseline
        if groupNumber == 0:
            ax.plot(
                num_keys[start:stop],
                baseline[start:stop][:, 0]/div,
                linestyle=linestyles[0],
                color=colors[0],
                marker=markers[0],
                lw=linewidth,
                label=str("FSS (baseline)")
            )

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

    eq_baseline_processing = []
    eq_dpf_pacl_processing = []
    eq_dpf_sk_pacl_processing = []

    range_baseline_processing = []
    range_dpf_pacl_processing = []
    range_dpf_sk_pacl_processing = []

    num_results = 0
    num_trials = len(results[0]["equality_baseline_processing_us"])

    # first we extract the relevant bits
    for i in range(len(results)):

        num_keys.append(results[i]["num_keys"])
        num_subkeys.append(results[i]["num_subkeys"])

        avg = np.mean(results[i]["equality_baseline_processing_us"])
        std = np.std(results[i]["equality_baseline_processing_us"])
        eq_baseline_processing.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["equality_dpf_pacl_processing_us"])
        std = np.std(results[i]["equality_dpf_pacl_processing_us"])
        eq_dpf_pacl_processing.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["equality_dpf_sk_pacl_processing_us"])
        std = np.std(results[i]["equality_dpf_sk_pacl_processing_us"])
        eq_dpf_sk_pacl_processing.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["range_baseline_processing_us"])
        std = np.std(results[i]["range_baseline_processing_us"])
        range_baseline_processing.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["range_dpf_pacl_processing_us"])
        std = np.std(results[i]["range_dpf_pacl_processing_us"])
        range_dpf_pacl_processing.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["range_dpf_sk_pacl_processing_us"])
        std = np.std(results[i]["range_dpf_sk_pacl_processing_us"])
        range_dpf_sk_pacl_processing.append(
            [avg, confidence95(std, num_trials)])

    num_results += 1

    num_keys = np.array(num_keys)
    num_subkeys = np.array(num_subkeys)
    eq_baseline_processing = np.array(eq_baseline_processing)
    eq_dpf_pacl_processing = np.array(eq_dpf_pacl_processing)
    eq_dpf_sk_pacl_processing = np.array(eq_dpf_sk_pacl_processing)
    range_baseline_processing = np.array(range_baseline_processing)
    range_dpf_pacl_processing = np.array(range_dpf_pacl_processing)
    range_dpf_sk_pacl_processing = np.array(range_dpf_sk_pacl_processing)

    print("Num datapoints: " + str(len(num_keys)))
    sort = np.argsort(num_subkeys)
    num_keys = num_keys[sort]
    num_subkeys = num_subkeys[sort]
    eq_baseline_processing = eq_baseline_processing[sort]
    eq_dpf_pacl_processing = eq_dpf_pacl_processing[sort]
    eq_dpf_sk_pacl_processing = eq_dpf_sk_pacl_processing[sort]
    range_baseline_processing = range_baseline_processing[sort]
    range_dpf_pacl_processing = range_dpf_pacl_processing[sort]
    range_dpf_sk_pacl_processing = range_dpf_sk_pacl_processing[sort]

    xticks = [1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024]
    indices = np.nonzero(np.in1d(num_keys, xticks,))

    # PLOT SERVER CPU TIME

    fig, (ax1, ax3) = plt.subplots(1, 2)
    fig.set_size_inches(width, 1.75)

    amortize = False
    plot(ax1,
         num_keys[indices],
         num_subkeys[indices],
         eq_baseline_processing[indices],
         eq_dpf_pacl_processing[indices], amortize)
    ax1.set_xticks(xticks)
    ax1.set_xscale('log', base=2)
    ax1.set_title("DPF-PACL")
    ax1.set_xlabel('Number of evaluations')
    ax1.set_ylabel("Amortized CPU time ($\mu$s)")

    plot(ax3,
         num_keys[indices],
         num_subkeys[indices],
         range_baseline_processing[indices],
         range_dpf_pacl_processing[indices], amortize)
    ax3.set_xticks(xticks)
    # ax3.set_ylim(bottom=0)
    ax3.set_xscale('log', base=2)
    ax3.set_title("DMPF-PACL")
    ax3.set_xlabel('Number of evaluations')
    # ax3.set_ylabel("Amortized CPU time ($\mu$s)")

fig.tight_layout()

ax1.legend(ncol=2, loc='center', framealpha=1, bbox_to_anchor=(
    1.1, -0.65), facecolor='white', fancybox=False,  handlelength=handlelength)

# ax1.legend(loc='best', framealpha=1, facecolor='white', fancybox=False)
# fig.text(-0.05, 0.52, 'Server CPU time\n     (seconds)', va='center', rotation='vertical')

fig.savefig('plot_fss.pdf', bbox_inches='tight')
