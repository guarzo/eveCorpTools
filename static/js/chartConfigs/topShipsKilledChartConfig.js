// chartConfigs/topShipsKilledChartConfig.js
import { truncateLabel, commonOptions } from '../utils.js';

const topShipsKilledChartConfig = {
    id: 'topShipsKilledChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdTopShipsKilledData',
        ytd: 'ytdTopShipsKilledData',
        lastMonth: 'lastMTopShipsKilledData',
    },
    type: 'bar',
    options: {
        ...commonOptions,
        indexAxis: 'y',
        plugins: {
            ...commonOptions.plugins,
            legend: { display: false },
        },
        scales: {
            x: {
                ...commonOptions.scales.x,
                beginAtZero: true,
            },
            y: {
                ...commonOptions.scales.y,
                ticks: {
                    ...commonOptions.scales.y.ticks,
                    autoSkip: false,
                },
                grid: { display: false },
            },
        },
    },
    processData: function (data) {
        // ... processing logic remains the same ...
        const labels = data.map(item => item.ShipName || 'Unknown');
        const killCounts = data.map(item => item.KillCount || 0);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 15));

        const datasets = [{
            label: 'Ships Killed',
            data: killCounts,
            backgroundColor: 'rgba(255, 77, 77, 0.7)',
        }];

        return { labels: truncatedLabels, datasets, fullLabels };
    },
};

export default topShipsKilledChartConfig;
