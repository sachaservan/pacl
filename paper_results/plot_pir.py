from plot_common import *


def plot(baseline, pacl, db_size, item_bytes, cluster_size):
    ######################## PLOT CODE ########################
    ax = plt.figure().gca()
    ax.yaxis.grid(color=gridcolor, linestyle=linestyle, linewidth=0.5)
    # ax.xaxis.grid(color=gridcolor, linestyle=linestyle, linewidth=0.5)

    fig = matplotlib.pyplot.gcf()
    fig.set_size_inches(width, height/1.5)

    X = np.arange(cluster_size)
    bar_width = 0.12

    cluster = 0
    next_pos = -3*bar_width + bar_width/4
    for i in range(0, len(baseline), cluster_size):
        sort = np.argsort(db_size[i:i + cluster_size])
        ybaseline = baseline[i:i + cluster_size][sort]
        ypacl = pacl[i:i + cluster_size][sort]

        # eliminate odd clusters to create space
        # for j in range(1, cluster_size, 2):
        #     ybaseline[j] *= 0
        #     ypacl[j] *= 0

        ax.bar(
            X + next_pos,
            ybaseline[:, 0],
            hatch='',
            edgecolor='white',
            linewidth=0.25,
            width=bar_width,
            color=bar_colors[cluster],
            zorder=3,  # https://stackoverflow.com/questions/23357798/how-to-draw-grid-lines-behind-matplotlib-bar-graph
            yerr=ybaseline[:, 1],
        )

        ax.bar(
            X + next_pos + bar_width,
            ypacl[:, 0],
            hatch='xxxx',
            edgecolor='white',
            linewidth=0.25,
            color=bar_colors[cluster],
            width=bar_width,
            zorder=3,
            yerr=ypacl[:, 1],
        )

        # fake bars for legend only
        ax.bar(
            X + next_pos + bar_width,
            np.zeros(cluster_size),
            edgecolor='#333',
            linewidth=0.25,
            color=bar_colors[cluster],
            label=str(item_bytes[i]) + " B",
        )

        # x = np.arange(len(db_size))
        ax.set_xticks([x for x in range(0, len(X))])
        labels_raw = np.sort(np.unique(db_size)).tolist()
        labels = [r'$2^{{{}}}$'.format(
            int(math.log2(labels_raw[i]))) for i in range(0, len(labels_raw))]
        ax.set_xticklabels(labels)

        next_pos += bar_width * 2 + bar_width/8

        cluster += 1

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

    db_size = []
    item_bytes = []
    server_baseline_ms = []
    server_pacl_ms = []

    num_results = 0
    num_trials = len(results[0]["server_pir_processing_ms"])

    # first we extract the relevant bits
    for i in range(len(results)):

        db_size.append(results[i]["db_size"])
        item_bytes.append(results[i]["item_size"])

        avg = np.mean(results[i]["server_pir_processing_ms"])
        std = np.std(results[i]["server_pir_processing_ms"])
        server_baseline_ms.append([avg, confidence95(std, num_trials)])

        avg = np.mean(results[i]["server_pir_pacl_processing_ms"])
        std = np.std(results[i]["server_pir_pacl_processing_ms"])
        server_pacl_ms.append([avg, confidence95(std, num_trials)])

        num_results += 1

    db_size = np.array(db_size)
    item_bytes = np.array(item_bytes)
    server_pacl_ms = np.array(server_pacl_ms)
    server_baseline_ms = np.array(server_baseline_ms)

    item_bytes = item_bytes.astype(int)

    sort = np.argsort(item_bytes)  # cluster by the username bytes
    item_bytes = item_bytes[sort]
    db_size = db_size[sort]
    server_pacl_ms = server_pacl_ms[sort]
    server_baseline_ms = server_baseline_ms[sort]

    mask = np.ones(len(item_bytes), dtype=bool)
    # mask[db_size == 2**21] = False
    # mask[db_size == 2**22] = False
    # mask[db_size == 2**18] = False
    mask[db_size == 2**18] = False
    mask[db_size == 2**17] = False
    mask[db_size == 2**16] = False
    mask[db_size == 2**15] = False
    mask[db_size == 2**14] = False

    item_bytes = item_bytes[mask]
    db_size = db_size[mask]
    server_baseline_ms = server_baseline_ms[mask]
    server_pacl_ms = server_pacl_ms[mask]

    # PLOT SERVER CPU TIME

    cluster_size = len(np.unique(db_size))
    print("Cluster size: " + str(cluster_size))

    print(item_bytes)

    ax = plot(server_baseline_ms, server_pacl_ms,
              db_size, item_bytes, cluster_size)
    ax.set_xlabel('Database size (number of items)')
    ax.set_ylabel('Server CPU time\n(milliseconds)')
    # ax.set_xscale('log', base=2)
    # ax.set_yscale('log', base=10)
    ax.figure.tight_layout()
    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)

    # ghost axis for the pacl/vanilla legend
    styles = ['-', '-.']
    style_labels = ['PIR', 'PIR w/ VDPF-PACL']
    ax2 = ax.twinx()
    ax2.bar(np.NaN, np.NaN,
            label=style_labels[0], edgecolor='#333', color='#666', linewidth=0.5)
    ax2.bar(np.NaN, np.NaN, label=style_labels[1], edgecolor='#333',
            linewidth=0.5, hatch='xxxx', color='white')
    ax2.get_yaxis().set_visible(False)
    # ax2.set_yscale('log', base=10)
    ax2.spines['top'].set_visible(False)
    ax2.spines['right'].set_visible(False)

    ax.legend(title="Item size", loc='upper left', framealpha=1, bbox_to_anchor=(
        0, 0.95), ncol=3, fancybox=False, edgecolor=edgecolor, facecolor='white', handlelength=handlelength)
    ax2.legend(ncol=2, loc='upper left', bbox_to_anchor=(
        0, 1.125), fancybox=False,  handlelength=handlelength)

    ax.figure.savefig('plot_pir_server_processing.pdf', bbox_inches='tight')
