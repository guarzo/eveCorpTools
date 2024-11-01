// static/js/chartConfigs/topShipsKilledWordCloudConfig.js

import { getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Top Ships Killed Word Cloud
 */
const topShipsKilledWordCloudConfig = {
    id: 'topShipsKilledWordCloud',
    instance: {}, // Initialize as an object to handle multiple timeframes
    dataKeys: {
        mtd: { dataVar: 'mtdTopShipsKilledData', canvasId: 'topShipsKilledWordCloud_mtd' },
        ytd: { dataVar: 'ytdTopShipsKilledData', canvasId: 'topShipsKilledWordCloud_ytd' },
        lastMonth: { dataVar: 'lastMTopShipsKilledData', canvasId: 'topShipsKilledWordCloud_lastM' },
    },
    type: 'wordCloud', // Using wordCloud chart type from chartjs-chart-wordcloud
    dataType: 'array', // Specify that this chart expects array data
    options: getCommonOptions('Top Ships Killed Word Cloud', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const shipName = context.raw.text || 'Unknown';
                        const killCount = context.raw.weight || 0;
                        return `${shipName}: ${killCount} Killmails`;
                    },
                },
                // Optional: Adjust tooltip display settings if needed
                mode: 'nearest',
                intersect: true,
            },
            datalabels: {
                color: '#ffffff',
                anchor: 'center',
                align: 'center',
                formatter: (value, context) => {
                    const shipName = context.raw.text || 'Unknown';
                    const killCount = context.raw.weight || 0;
                    // Calculate percentage relative to the total kill count in the word cloud
                    const totalKillCount = context.chart.data.datasets[0].data.reduce((sum, item) => sum + (item.weight || 0), 0);
                    const percentage = totalKillCount > 0 ? ((killCount / totalKillCount) * 100).toFixed(1) : '0.0';
                    return `${killCount}\n(${percentage}%)`;
                },
                font: {
                    size: function(context) {
                        const weight = context.raw.weight || 1;
                        // Adjust font size based on kill count
                        return Math.min(40, 10 + weight * 2); // Example scaling
                    },
                    weight: 'bold',
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
        if (!validateChartDataArray(data, chartName)) {
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
