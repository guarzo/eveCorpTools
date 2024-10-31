// static/js/chartConfigs/topShipsKilledChartConfig.js

import { truncateLabel, getColor, getCommonOptions, validateChartData } from '../utils.js';

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
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const label = context.label || '';
                        const value = context.parsed.y !== undefined ? context.parsed.y : context.parsed.x;
                        return `Killmails: ${value}`;
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
                    text: 'Ship Types',
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
                    text: 'Kill Count',
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
        const chartName = 'Top Ships Killed Chart';
        if (!validateChartData(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Extract labels and data
        const labels = data.map(item => item.Name || 'Unknown');
        const truncatedLabels = labels.map(label => truncateLabel(label, 15));

        const killCounts = data.map(item => item.KillCount || 0);

        // Define datasets
        const datasets = [{
            label: 'Killmails',
            data: killCounts,
            backgroundColor: 'rgba(75, 192, 192, 0.7)',
            borderColor: 'rgba(75, 192, 192, 1)',
            borderWidth: 1,
        }];

        return { labels: truncatedLabels, datasets, fullLabels: labels };
    },
};

export default topShipsKilledChartConfig;
