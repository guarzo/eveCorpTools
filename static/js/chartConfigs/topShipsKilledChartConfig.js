// static/js/chartConfigs/topShipsKilledChartConfig.js

import { truncateLabel, getCommonOptions } from '../utils.js';

/**
 * Configuration for the Top Ships Killed Chart
 */
const topShipsKilledChartConfig = {
    id: 'topShipsKilledChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdTopShipsKilledData',
        ytd: 'ytdTopShipsKilledData',
        lastMonth: 'lastMTopShipsKilledData',
    },
    type: 'bar',
    options: getCommonOptions('Top Ships Killed', {
        indexAxis: 'y',
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    // Custom tooltip to display Ship Name and Kill Count
                    label: function (context) {
                        const shipName = context.dataset.label;
                        const killCount = context.parsed.x;
                        return `${shipName}: ${killCount}`;
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'linear',
                beginAtZero: true,
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
            },
            y: {
                type: 'category',
                labels: [], // Labels are set dynamically in processData
                ticks: {
                    color: '#ffffff',
                    autoSkip: false,
                },
                grid: { display: false },
            },
        },
    }),
    processData: function (data) {
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
