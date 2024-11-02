// static/js/chartConfigs/averageFleetSizeChartConfig.js

import { truncateLabel, getColor, getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Average Fleet Size Over Time Chart
 */
const averageFleetSizeChartConfig = {
    id: 'averageFleetSizeChart',
    instance: {}, // Initialize as an object to store chart instances per timeframe
    dataKeys: {
        mtd: { dataVar: 'mtdAverageFleetSizeData', canvasId: 'averageFleetSizeChart_mtd' },
        ytd: { dataVar: 'ytdAverageFleetSizeData', canvasId: 'averageFleetSizeChart_ytd' },
        lastMonth: { dataVar: 'lastMAverageFleetSizeData', canvasId: 'averageFleetSizeChart_lastM' },
    },
    type: 'line', // Using line chart
    dataType: 'array', // Expecting an array of FleetSizeData
    options: getCommonOptions('Average Fleet Size Over Time', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const avgFleetSize = context.parsed.y !== undefined ? context.parsed.y : 0;
                        return `Average Fleet Size: ${avgFleetSize}`;
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'time',
                time: {
                    unit: 'day', // Change to 'week' if interval is weekly
                    tooltipFormat: 'MMM dd, yyyy',
                },
                title: {
                    display: true,
                    text: 'Time',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
                ticks: {
                    color: '#ffffff',
                },
                grid: { display: false },
            },
            y: {
                beginAtZero: true,
                title: {
                    display: true,
                    text: 'Average Fleet Size',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
                ticks: {
                    color: '#ffffff',
                },
                grid: { color: '#444' },
            },
        },
    }),
    processData: function (data) {
        const chartName = 'Average Fleet Size Over Time Chart';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Prepare labels and data
        const labels = data.map(item => item.Time);
        const fleetSizes = data.map(item => item.FleetSize);

        // Check if there are at least 7 days of data
        if (fleetSizes.length < 3) {
            console.warn(`Not enough data points (${fleetSizes.length}) for ${chartName}.`);
            return { labels: [], datasets: [], noDataMessage: 'Not enough data to display the chart.' };
        }

        const datasets = [{
            label: 'Average Fleet Size',
            data: fleetSizes,
            backgroundColor: 'rgba(75, 192, 192, 0.7)',
            borderColor: 'rgba(75, 192, 192, 1)',
            borderWidth: 2,
            fill: true,
            tension: 0.1,
        }];

        return { labels, datasets };
    },
};

export default averageFleetSizeChartConfig;
