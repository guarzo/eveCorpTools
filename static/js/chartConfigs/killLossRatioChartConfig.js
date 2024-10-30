// static/js/chartConfigs/killLossRatioChartConfig.js

import { truncateLabel, getCommonOptions } from '../utils.js';

/**
 * Configuration for the Kill-to-Loss Ratio Chart
 */
const killLossRatioChartConfig = {
    id: 'killLossRatioChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdKillLossRatioData',
        ytd: 'ytdKillLossRatioData',
        lastMonth: 'lastMKillLossRatioData',
    },
    type: 'bar',
    options: getCommonOptions('Kill-to-Loss Ratio', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const data = context.chart.data.originalData[context.dataIndex];
                        const kills = data.Kills || 0;
                        const losses = data.Losses || 0;
                        const ratio = context.parsed.y.toFixed(2);
                        return `Kills: ${kills}, Losses: ${losses}, Ratio: ${ratio}`;
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'category',
                ticks: {
                    color: '#ffffff',
                    maxRotation: 45,
                    minRotation: 45,
                    autoSkip: false,
                },
                grid: { display: false },
            },
            y: {
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                beginAtZero: true,
            },
        },
    }),
    processData: function (data) {
        const labels = data.map(item => item.CharacterName || 'Unknown');
        const ratios = data.map(item => item.Ratio || 0);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const datasets = [{
            label: 'Kill-to-Loss Ratio',
            data: ratios,
            backgroundColor: 'rgba(75, 192, 192, 0.7)',
        }];

        return { labels: truncatedLabels, datasets, fullLabels, originalData: data };
    },
};

export default killLossRatioChartConfig;
