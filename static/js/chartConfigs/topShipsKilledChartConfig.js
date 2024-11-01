// static/js/chartConfigs/topShipsKilledWordCloudConfig.js

import { getCommonOptions, validateChartData } from '../utils.js';

/**
 * Configuration for the Top Ships Killed Word Cloud
 */
const topShipsKilledWordCloudConfig = {
    id: 'topShipsKilledWordCloud',
    instance: null,
    dataKeys: {
        mtd: 'mtdTopShipsKilledData',
        ytd: 'ytdTopShipsKilledData',
        lastMonth: 'lastMTopShipsKilledData',
    },
    type: 'wordCloud', // Using wordCloud chart type from chartjs-chart-wordcloud
    options: getCommonOptions('Top Ships Killed Word Cloud', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const shipName = context.label || 'Unknown';
                        const killCount = context.parsed || 0;
                        return `${shipName}: ${killCount} Killmails`;
                    },
                },
            },
        },
        scales: {
            x: {
                display: false,
            },
            y: {
                display: false,
            },
        },
        responsive: true,
        maintainAspectRatio: false,
    }),
    processData: function (data) {
        const chartName = 'Top Ships Killed Word Cloud';
        if (!validateChartData(data, chartName)) {
            // Return empty data to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Convert data to {text, weight} format required by word cloud
        const wordCloudData = data.map(item => ({
            text: item.Name || 'Unknown',
            weight: item.KillCount || 1,
        }));

        const dataset = {
            label: 'Killmails',
            data: wordCloudData,
            // Optional: Customize font size based on weight
            fontSize: function(context) {
                const weight = context.parsed.weight;
                return Math.max(10, weight); // Adjust as needed
            },
            rotation: function() {
                const rotations = [-90, 0, 90];
                return rotations[Math.floor(Math.random() * rotations.length)];
            },
            // Optional: Customize colors
            color: function(context) {
                const index = context.dataIndex;
                const colors = [
                    '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0',
                    '#9966FF', '#FF9F40', '#C9CBCF', '#8B0000',
                    '#008000', '#00008B', // Add more colors as needed
                ];
                return colors[index % colors.length];
            },
        };

        const labels = []; // Not used in word cloud
        const datasets = [dataset];

        return { labels, datasets };
    },
};

export default topShipsKilledWordCloudConfig;
