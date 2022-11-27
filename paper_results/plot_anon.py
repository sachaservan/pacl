from plot_common import *


def plot(ax, a, b, num_accounts, colora="#c0392b", colorb="#8e44ad", labela="", labelb=""):
    ######################## PLOT CODE ########################
    ax.plot(
        num_accounts,
        a[:, 0],
        linestyle="-",
        color=colora,
        lw=linewidth,
        marker='.',
        label=labela,
    )
    ax.fill_between(
        num_accounts,
        a[:, 0] - a[:, 1],
        a[:, 0] + a[:, 1],
        color=colora,
        alpha=error_opacity,
    )

    ax.plot(
        num_accounts,
        b[:, 0],
        linestyle="--",
        color=colorb,
        lw=linewidth,
        marker='x',
        label=labelb,
    )
    ax.fill_between(
        num_accounts,
        b[:, 0] - b[:, 1],
        b[:, 0] + b[:, 1],
        color=colorb,
        alpha=error_opacity,
    )

    return ax


if __name__ == '__main__':
    argparser = argparse.ArgumentParser(sys.argv[0])
    argparser.add_argument("--file", type=str, default='')

    height = 1.75  # override height parameter

    args = argparser.parse_args()

    # read experiment file (expected json)
    with open(args.file, 'r') as myfile:
        data = myfile.read()

    # parse the experiment file as json
    results = json.loads(data)

    num_keys = []
    vanilla_express = []
    pacl_express = []
    vanilla_spectrum = []
    pacl_spectrum = []

    num_results = 0
    num_trials = len(results[0]["server_express_ms"])

    # first we extract the relevant bits
    for i in range(len(results)):

        num_keys.append(results[i]["num_keys"])

        avg = np.mean(results[i]["server_express_ms"])
        std = np.std(results[i]["server_express_ms"])
        vanilla_express.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["server_express_pacl_ms"])
        std = np.std(results[i]["server_express_pacl_ms"])
        pacl_express.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["server_spectrum_ms"])
        std = np.std(results[i]["server_spectrum_ms"])
        vanilla_spectrum.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["server_spectrum_pacl_ms"])
        std = np.std(results[i]["server_spectrum_pacl_ms"])
        pacl_spectrum.append([avg, confidence95(std, num_trials)])

        num_results += 1

    num_keys = np.array(num_keys)
    pacl_express = np.array(pacl_express)
    vanilla_express = np.array(vanilla_express)
    pacl_spectrum = np.array(pacl_spectrum)
    vanilla_spectrum = np.array(vanilla_spectrum)

    sort = np.argsort(num_keys)
    num_keys = num_keys[sort]
    pacl_express = pacl_express[sort]
    vanilla_express = vanilla_express[sort]
    pacl_spectrum = pacl_spectrum[sort]
    vanilla_spectrum = vanilla_spectrum[sort]

    vanilla_express *= MILLI_TO_SECONDS
    pacl_express *= MILLI_TO_SECONDS

    # TODO: maybe don't hardcode these?
    xticks = [16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152]
    indices = np.nonzero(np.in1d(num_keys, xticks,))

    fig, (ax1, ax2) = plt.subplots(1, 2, sharex=False, sharey=False)

    ax1.yaxis.grid(color=gridcolor, linestyle=linestyle, linewidth=0.5)
    ax2.yaxis.grid(color=gridcolor, linestyle=linestyle, linewidth=0.5)
    ax1.xaxis.grid(color=gridcolor, linestyle=linestyle, linewidth=0.5)
    ax2.xaxis.grid(color=gridcolor, linestyle=linestyle, linewidth=0.5)
    ax1.spines['top'].set_visible(False)
    ax1.spines['right'].set_visible(False)
    ax2.spines['top'].set_visible(False)
    ax2.spines['right'].set_visible(False)
    fig.set_size_inches(width, height)

    plot(ax1, vanilla_express[indices], pacl_express[indices], num_keys[indices],
         colora=colors[0], colorb=colors[0], labela="Express", labelb="Express w/ PACL")
    ax1.set_xlabel('Number of mailboxes')
    ax1.set_xticks(xticks)
    ax1.legend(loc='upper left', ncol=1, framealpha=1, facecolor='white',
               bbox_to_anchor=(0, 1.15), fancybox=False, handlelength=handlelength)
    ax1.set_xscale('log', base=2)

    ax1.annotate("",
                 xy=(1048576, 0.5), xycoords='data',
                 xytext=(1048576, 5), textcoords='data',
                 arrowprops=annotate_arrow)
    ax1.annotate(r'57$\times$',
                 xy=(1150000, 2.55),
                 xycoords='data',
                 fontsize=annoate_font_size,
                 color=annoate_text_color)

    vanilla_spectrum *= MILLI_TO_SECONDS
    pacl_spectrum *= MILLI_TO_SECONDS

    plot(ax2, vanilla_spectrum[indices], pacl_spectrum[indices],
         num_keys[indices], colora=colors[1], colorb=colors[1],
         labela="Spectrum", labelb="Spectrum w/ PACL")
    ax2.set_xlabel('Number of mailboxes')
    ax2.set_xticks(xticks)
    ax2.set_xscale('log', base=2)
    ax2.legend(loc='upper left', ncol=1, framealpha=1, facecolor='white',
               bbox_to_anchor=(0, 1.15), fancybox=False, handlelength=handlelength)

    ax2.annotate("",
                 xy=(1048576, 2), xycoords='data',
                 xytext=(1048576, 25), textcoords='data',
                 arrowprops=annotate_arrow)
    ax2.annotate(r'70$\times$',
                 xy=(1150000, 12),
                 xycoords='data',
                 fontsize=annoate_font_size,
                 color=annoate_text_color)

    fig.text(-0.02, 0.52, 'Server CPU time (seconds)',
             va='center', rotation='vertical')

    fig.tight_layout()
    fig.savefig('plot_express_spectrum.pdf', bbox_inches='tight')
