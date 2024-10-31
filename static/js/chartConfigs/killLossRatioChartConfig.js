// static/js/chartConfigs/killLossRatioChartConfig.js

import { truncateLabel, getColor, getCommonOptions, validateChartData } from '../utils.js';

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
                        const dataPoint = context.raw; // Access the data point directly
                        const kills = dataPoint.kills || 0;
                        const losses = dataPoint.losses || 0;
                        const ratio = dataPoint.y !== undefined ? dataPoint.y.toFixed(2) : '0.00';
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
                title: {
                    display: true,
                    text: 'Characters',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
            y: {
                beginAtZero: true,
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                title: {
                    display: true,
                    text: 'Ratio',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
        },
    }),
    processData: function (data) {
        const chartName = 'Kill-to-Loss Ratio Chart';
        if (!validateChartData(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        const labels = data.map(item => item.CharacterName || 'Unknown');
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const ratios = data.map(item => item.Ratio || 0);

        const datasets = [{
            label: 'Kill-to-Loss Ratio',
            data: data.map(item => ({
                x: item.CharacterName || 'Unknown',
                y: item.Ratio || 0,
                kills: item.Kills || 0,
                losses: item.Losses || 0
            })),
            backgroundColor: 'rgba(75, 192, 192, 0.7)',
            borderColor: 'rgba(75, 192, 192, 1)',
            borderWidth: 1,
        }];

        return { labels: truncatedLabels, datasets, fullLabels: labels };
    },
};

export default killLossRatioChartConfig;
